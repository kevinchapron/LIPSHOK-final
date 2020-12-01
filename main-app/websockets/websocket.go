package websockets

import (
	"github.com/gorilla/mux"
	"net/http"
)

func CreateWebSocket(name string, uri string, r *mux.Router) {

	ws_wrapper := WebSocketWrapper{}
	hub, ok := hubs[name]
	if !ok {
		hubs[name] = &WebSocketHub{
			clients:       make(map[*WebSocketClient]bool),
			connecting:    make(chan *WebSocketClient),
			disconnecting: make(chan *WebSocketClient),

			Sender:   make(chan WebSocketMessage),
			Receiver: make(chan WebSocketMessage),

			IndexKey: name,
		}
		hub = hubs[name]
	}
	ws_wrapper.Hub = hub
	go ws_wrapper.Hub.run()

	r.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		CreateClientConnection(ws_wrapper.Hub, writer, request)
	})

}
