package websockets

type WebSocketWrapper struct {
	Hub *WebSocketHub
}

func (w *WebSocketWrapper) SendAllClients(m WebSocketMessage) {
	for c, b := range w.Hub.clients {
		if b {
			c.sending <- m
		}
	}
}
