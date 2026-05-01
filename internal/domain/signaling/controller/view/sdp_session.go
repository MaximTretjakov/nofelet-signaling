package view

type SDPData struct {
	Type        string       `json:"type"`                  // Тип сообщения (join, offer, answer, ice-candidates...)
	SDP         string       `json:"sdp,omitempty"`         // SDP сообщение
	Candidate   IceCandidate `json:"candidate,omitempty"`   // Candidate кандидаты разные
	Participant Participant  `json:"participant,omitempty"` // Дополнительная информация о клиенте
}

type IceCandidate struct {
	Candidate        string `json:"candidate"`                  // Протокол (UDP/TCP), IP-адрес, порт и тип кандидата
	SdpMid           string `json:"sdpMid"`                     // Идентификатор медиа-потока («audio», «video» или «data»)
	SdpMLineIndex    int    `json:"sdpMLineIndex"`              // Индекс строки в SDP-описании (начиная с нуля), указывающий на конкретный медиа-блок
	UsernameFragment string `json:"usernameFragment,omitempty"` // Короткий токен (ufrag) гарантирует, что пакеты приходят именно от того участника, с которым вы пытаетесь соединиться
}

type Participant struct {
	Name string `json:"name"` // Ник участника звонка
	Role string `json:"role"` // Роль участника звонка (Инициатор звонка, собеседник)
}
