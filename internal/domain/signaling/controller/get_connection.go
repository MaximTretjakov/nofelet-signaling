package controller

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"nofelet/config"
	"nofelet/internal/domain/signaling/controller/view"
	"nofelet/pkg/singleton"
)

const (
	caller = "caller"
	callee = "callee"
)

// GetConnection /connect/:uuid установка sdp сессии
func (c *Controller) GetConnection(ctx *gin.Context) {
	conn, sErr := Upgrader(ctx)
	if sErr != nil {
		c.Logger.Error("socket creation", slog.Any("err", sErr))
	}

	uuid := ctx.Param("uuid")

	room := singleton.NewRoom()
	room.Init(uuid)

	go handler(conn, uuid, room, c.Logger)
}

// handler - обрабатывает коннекты участников
func handler(conn *websocket.Conn, uuid string, room *singleton.Rooms, logger *slog.Logger) {
	defer func() {
		room.DeleteClient(uuid)
		_ = conn.Close()
	}()

	var data view.SDPData

	for {
		readErr := conn.ReadJSON(&data)
		if readErr != nil {
			logger.Error("socket read", slog.Any("err", readErr))
			break
		}

		r := room.Rooms[uuid]

		switch data.Type {
		case "join":
			jErr := Join(data, conn, r, room)
			if jErr != nil {
				logger.Error("handler", slog.Any("join error:", jErr))
			}
		case "offer":
			oErr := Offer(data, conn, r, room)
			if oErr != nil {
				logger.Error("handler", slog.Any("offer error:", oErr))
			}
		case "ice-candidate":
			iceErr := IceCandidate(data, conn, r, room)
			if iceErr != nil {
				logger.Error("handler", slog.Any("ice-candidate error:", iceErr))
			}
		case "answer":
			brErr := room.Broadcast(data, conn)
			if brErr != nil {
				logger.Error("handler", slog.Any("answer error:", brErr))
			}
		}

		// Логирование временное
		if config.Current().Debug {
			printSocketData(data, logger, conn)
		}
	}
}

func printSocketData(data view.SDPData, logger *slog.Logger, conn *websocket.Conn) {
	fmt.Println()
	message := fmt.Sprintf("from=%s | data=%+v\n",
		conn.RemoteAddr().String(),
		data,
	)
	logger.Info("wss", slog.String(":", message))
}
