package internal_connectors

import (
	"github.com/gorilla/websocket"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/messaging"
	"net/url"
	"strconv"
)

type internalConnector struct {
	wsConnection *websocket.Conn
	terminate    chan struct{}
	IsConnected  chan bool

	messages *chan messaging.Message
}

func CreateInternalConnector() *internalConnector {
	return &internalConnector{terminate: make(chan struct{}), IsConnected: make(chan bool)}
}

func (conn *internalConnector) Connect(list *chan messaging.Message) {

	u := url.URL{Scheme: "ws", Host: constants.ROUTER_ADDRESS + ":" + strconv.Itoa(constants.SENSOR_WEBSOCKET_PORT), Path: constants.SENSOR_WEBSOCKET_PATH}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		conn.IsConnected <- false
		panic(err)
	}
	conn.messages = list
	conn.IsConnected <- true
	conn.wsConnection = c
	go conn.internalConnection()
}

func (conn *internalConnector) internalConnection() {
	// Action with websocket
	defer conn.wsConnection.Close()
	// listen for messages to send

	go conn.continuouslyReadMessages()
	for {
		select {
		case <-conn.terminate:
			return
		case msg := <-*conn.messages:
			err := conn.wsConnection.WriteMessage(websocket.BinaryMessage, msg.ToBytes())
			if err != nil {
				Logging.Error(err)
				continue
			}
			// create other cases
		}
	}
}

func (conn *internalConnector) continuouslyReadMessages() {
	defer close(conn.terminate)
	for {
		_, message, err := conn.wsConnection.ReadMessage()
		if err != nil {
			Logging.Error(err)
			return
		}
		Logging.Debug("Received Message : ", string(message))
	}
}
