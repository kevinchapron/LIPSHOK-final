package websockets

import (
	"encoding/json"
	"errors"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"strings"
	"time"
)

func RegisterNewFunction(name string, f func(message messaging.Message, client *WebSocketClient, hub *WebSocketHub)) error {
	if _, ok := WebSocketFunction[name]; ok {
		return errors.New("Already a function is registered for this name.")
	}
	WebSocketFunction[name] = f
	return nil
}

var WebSocketFunction = map[string]func(messaging.Message, *WebSocketClient, *WebSocketHub){
	constants.SENSOR_WEBSOCKET_NAME: SensorData,
	constants.OUTPUT_WEBSOCKET_NAME: OutputData,
}

func SensorData(msg messaging.Message, client *WebSocketClient, hub *WebSocketHub) {

	if msg.DataType == constants.MESSAGING_DATATYPE_AUTH {
		Logging.Debug("Received Identification for", client.conn.RemoteAddr().String())
		var m = make(map[string]interface{})
		json.Unmarshal(msg.Data, &m)
		if client.Name == "" {
			client.Protocol = m["Protocol"].(string)
			client.Name = m["Name"].(string)
		} else {
			GetListSensor().UpdateSensor(&WebSocketClient{Protocol: m["Protocol"].(string), Name: m["Name"].(string)})
		}
		return
	}

	var fromMessaging = client
	if msg.From != "" {
		var m = WebSocketClient{}
		splits := strings.Split(msg.From, "-")
		m.Name = splits[0]
		m.Protocol = splits[1]
		fromMessaging = &m
	}

	BroadcastToOutput(WebSocketMessage{
		Data: msg.Data,
		From: fromMessaging,
		To:   nil,
		Type: 0,
	})
	//Logging.Debug("Received: ", msg, "; From:", client.conn.RemoteAddr().String())
	//Logging.Debug("--->", string(msg.Data))

}

func OutputData(msg messaging.Message, client *WebSocketClient, hub *WebSocketHub) {
	//
	//Logging.Debug("Received: ", msg, "; From:", client.conn.RemoteAddr().String())
	//Logging.Debug("--->", string(msg.Data))

	if string(msg.Data) == "status" {
		// Client want global status.
		Logging.Debug(" --> Client want global status.")

		var listConnectors = []map[string]interface{}{}
		var clients = ListAllConnectors()
		for _, client := range clients {
			var m = make(map[string]interface{})
			m["addr"] = client.conn.RemoteAddr().String()
			m["lastSeen"] = client.lastMessageTime.String()
			m["name"] = client.Name
			m["protocol"] = client.Protocol
			listConnectors = append(listConnectors, m)
		}

		var returnValue = make(map[string]interface{})
		returnValue["connectors"] = listConnectors
		returnValue["sensors"] = GetListSensor().ListAllSensors(nil)
		returnBytes, _ := json.Marshal(returnValue)

		client.sending <- WebSocketMessage{Data: returnBytes, Type: constants.MESSAGING_DATATYPE_AUTH}

		return
	}

}

func BroadcastToOutput(msg WebSocketMessage) {
	//Logging.Info("[WS] Broadcasting from ",msg.From.Name," (",msg.From.Protocol,") ...")
	GetListSensor().UpdateSensor(msg.From)

	for client, b := range hubs[constants.OUTPUT_WEBSOCKET_NAME].clients {
		if !b {
			continue
		}
		msg.Datetime = time.Now()
		client.sending <- msg
	}
}
