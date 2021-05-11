package websockets

import (
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
)

var hubs = make(map[string]*WebSocketHub)

type WebSocketHub struct {
	clients       map[*WebSocketClient]bool
	connecting    chan *WebSocketClient
	disconnecting chan *WebSocketClient

	Sender   chan WebSocketMessage
	Receiver chan WebSocketMessage

	IndexKey string
	Raw      bool
}

func (h *WebSocketHub) run() {
	for {
		select {
		case c := <-h.connecting:
			h.clients[c] = true
		case c := <-h.disconnecting:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.sending)
			}
		case msg := <-h.Sender:
			if msg.To != nil {
				msg.To.sending <- msg
			}
		case msg := <-h.Receiver:

			var m messaging.Message
			switch h.Raw {
			case false:
				err := m.FromBytes(msg.Data)
				if err != nil {
					Logging.Error("Message received from", msg.From.conn.RemoteAddr().String(), "; err : ", err)
					continue
				}
			case true:
				m.Data = msg.Data
			}
			WebSocketFunction[h.IndexKey](m, msg.From, h)

		}
	}
}

func ListAllConnectors() []*WebSocketClient {
	var m []*WebSocketClient
	for c, _ := range hubs[constants.SENSOR_WEBSOCKET_NAME].clients {
		m = append(m, c)
	}
	return m
}
