package singleton

import (
	"sync"

	"github.com/gorilla/websocket"

	"nofelet/internal/domain/signaling/controller/view"
)

var (
	once     sync.Once
	instance *Rooms
)

type Room struct {
	Initiator    *Participant // Участник звонка
	Callee       *Participant
	PendingOffer *view.SDPData // Данные SDP сессии
}

type Participant struct {
	Conn     *websocket.Conn // Объект коннекшена участника звонка
	Nickname string          // Ник участника звонка
	Role     string          // Роль участника звонка
}

// Rooms - хранит и управляет комнатами
type Rooms struct {
	Mu    sync.RWMutex
	Rooms map[string]*Room
}

// NewRoom - создает коммнату
func NewRoom() *Rooms {
	once.Do(func() {
		instance = &Rooms{
			Rooms: make(map[string]*Room),
		}
	})
	return instance
}

// Init - инициализирует комнату с конкретным uuid и инициализирует ее дефолтно
func (cm *Rooms) Init(uuid string) {
	cm.Mu.Lock()
	if _, ok := cm.Rooms[uuid]; !ok {
		cm.Rooms[uuid] = &Room{
			Initiator:    &Participant{},
			Callee:       &Participant{},
			PendingOffer: &view.SDPData{},
		}
	}
	cm.Mu.Unlock()
}

// DeleteClient - удаляем коннекшен клиента
func (cm *Rooms) DeleteClient(uuid string) {
	cm.Mu.Lock()
	delete(cm.Rooms, uuid)
	cm.Mu.Unlock()
}

// Connections - возвращает количество клиентов
func (cm *Rooms) Connections() int {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	return len(cm.Rooms)
}

// Broadcast - рассылает сообщения все клиентам доя установления SDP сессии
func (cm *Rooms) Broadcast(data view.SDPData, sender *websocket.Conn) error {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()

	for roomID, room := range cm.Rooms {
		initiator := room.Initiator.Conn
		callee := room.Callee.Conn

		if initiator == sender {
			err := callee.WriteJSON(data)
			if err != nil {
				_ = initiator.Close()
				delete(cm.Rooms, roomID)
				return err
			}
		}

		if callee == sender {
			err := initiator.WriteJSON(data)
			if err != nil {
				_ = callee.Close()
				delete(cm.Rooms, roomID)
				return err
			}
		}
	}

	return nil
}
