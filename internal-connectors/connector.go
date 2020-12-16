package internal_connectors

import (
	"github.com/gorilla/websocket"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"net/url"
	"strconv"
)

type internalConnector struct {
	wsConnection *websocket.Conn
	terminate    chan struct{}
	IsConnected  chan bool
}

func CreateInternalConnector() *internalConnector {
	return &internalConnector{terminate: make(chan struct{}), IsConnected: make(chan bool)}
}

func (conn *internalConnector) Connect() {
	u := url.URL{Scheme: "ws", Host: constants.ROUTER_ADDRESS + ":" + strconv.Itoa(constants.SENSOR_WEBSOCKET_PORT), Path: constants.SENSOR_WEBSOCKET_PATH}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		conn.IsConnected <- false
		panic(err)
	}
	conn.IsConnected <- true
	conn.wsConnection = c
	go conn.internalConnection()
}

func (conn *internalConnector) internalConnection() {
	// Action with websocket
	defer conn.wsConnection.Close()

	go conn.continuouslyReadMessages()
	for {
		select {
		case <-conn.terminate:
			return
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
