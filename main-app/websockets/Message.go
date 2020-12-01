package websockets

type WebSocketMessage struct {
	Data []byte           `json:"data"`
	From *WebSocketClient `json:"from"`
	To   *WebSocketClient `json:"-"`

	Type int `json:"type"`
}
