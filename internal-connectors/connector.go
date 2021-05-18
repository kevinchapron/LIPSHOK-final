package internal_connectors

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"github.com/kevinchapron/LIPSHOK/security"
	"net/url"
	"strconv"
)

var listSensor = make(map[string]deviceStruct)

type deviceStruct struct {
	Name     string
	Protocol string
}

func (d *deviceStruct) String() string {
	return "\"" + d.Name + " (" + d.Protocol + ")\""
}
func GetSensorDetails(ip string) *deviceStruct {
	if data, exists := listSensor[ip]; exists {
		return &data
	}
	return nil
}

type internalConnector struct {
	wsConnection *websocket.Conn `json:"-"`
	terminate    chan struct{}   `json:"-"`
	IsConnected  chan bool       `json:"-"`

	Name     string
	Protocol string

	messages *chan messaging.Message `json:"-"`
}

func CreateInternalConnector(name string, protocol string) *internalConnector {
	return &internalConnector{terminate: make(chan struct{}), IsConnected: make(chan bool), Name: name, Protocol: protocol}
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
	var authMsg messaging.Message

	authMsg.AesIV = security.RandomKey()
	authMsg.DataType = constants.MESSAGING_DATATYPE_AUTH
	v, _ := json.Marshal(*conn)
	authMsg.Data = v

	conn.wsConnection.WriteMessage(websocket.BinaryMessage, authMsg.ToBytes())

	defer conn.wsConnection.Close()
	// listen for messages to send

	//go conn.continuouslyReadMessages()
	for {
		select {
		case <-conn.terminate:
			return
		case msg := <-*conn.messages:
			// verifier provenance messages
			if msg.DataType == constants.MESSAGING_DATATYPE_AUTH {
				var tmp deviceStruct
				json.Unmarshal(msg.Data, &tmp)

				listSensor[msg.From] = tmp
			}

			sensor, exists := listSensor[msg.From]
			if !exists {
				Logging.Warning("Received:", msg, ". Yet, nothing registered. Ignoring.")
				continue
			}

			msg.From = sensor.Name + "-" + sensor.Protocol

			err := conn.wsConnection.WriteMessage(websocket.BinaryMessage, msg.ToBytes())
			if err != nil {
				Logging.Error(err)
				continue
			}
			// create other cases
		}
	}
}

//
//func (conn *internalConnector) continuouslyReadMessages() {
//	defer close(conn.terminate)
//	for {
//		_, message, err := conn.wsConnection.ReadMessage()
//		if err != nil {
//			Logging.Error(err)
//			return
//		}
//		Logging.Debug("Received Message : ", string(message))
//	}
//}
