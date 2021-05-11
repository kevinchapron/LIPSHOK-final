package websockets

import "time"

type WebSocketMessage struct {
	Data     []byte           `json:"data"`
	From     *WebSocketClient `json:"from"`
	To       *WebSocketClient `json:"-"`
	Datetime time.Time        `json:"datetime"`

	Type int `json:"type"`
}
