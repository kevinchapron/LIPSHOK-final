package websockets

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  WEBSOCKET_MAX_PACKET_SIZE,
	WriteBufferSize: WEBSOCKET_MAX_PACKET_SIZE,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
