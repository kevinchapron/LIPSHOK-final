package websockets

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK"
	"github.com/kevinchapron/FSHK-final/main-app/database"
	"net/http"
	"time"
)

type WebSocketClient struct {
	hub    *WebSocketHub
	conn   *websocket.Conn
	device *database.DatabaseDevice

	sending         chan WebSocketMessage
	lastMessageTime time.Time
	Name            string
	Protocol        string
}

func (c *WebSocketClient) read() {
	defer func() {
		c.hub.disconnecting <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(FSHK.WEBSOCKET_MAX_SIZE)
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			Logging.Info("[WS] > Client ", c.conn.RemoteAddr().String(), "disconnected.")
			break
		}
		c.lastMessageTime = time.Now()

		var m WebSocketMessage
		m.Data = msg
		m.From = c
		c.hub.Receiver <- m

	}
}

func (c *WebSocketClient) write() {
	defer func() {
		c.conn.Close()
	}()
	for {
		msg, ok := <-c.sending

		encoded_msg, err := json.Marshal(msg)
		if err != nil {
			Logging.Error(err)
			continue
		}

		c.conn.SetWriteDeadline(time.Now().Add(WEBSOCKET_MAX_WRITE_TIME))
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(encoded_msg)

		if err := w.Close(); err != nil {
			return
		}

		Logging.Debug(fmt.Sprintf("[WS] > Sent \"%s\" to %s.", encoded_msg, c.conn.RemoteAddr()))
	}
}

func CreateClientConnection(hub *WebSocketHub, w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		Logging.Error(err)
		return
	}
	c := &WebSocketClient{hub: hub, conn: conn, sending: make(chan WebSocketMessage)}
	c.hub.connecting <- c

	go c.read()
	go c.write()

	Logging.Info("[WS] > New Client connecting : ", c.conn.RemoteAddr())
}
