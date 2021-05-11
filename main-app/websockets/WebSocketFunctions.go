package websockets

import (
	"encoding/json"
	"errors"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
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
		client.Name = m["Name"].(string)
		client.Protocol = m["Protocol"].(string)

		return
	}

	BroadcastToOutput(WebSocketMessage{
		Data: msg.Data,
		From: client,
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

		var returnValue = []map[string]interface{}{}

		var clients = ListAllConnectors()
		for _, client := range clients {
			var m = make(map[string]interface{})
			m["addr"] = client.conn.RemoteAddr().String()
			m["lastSeen"] = client.lastMessageTime.String()
			m["name"] = client.Name
			m["protocol"] = client.Protocol
			returnValue = append(returnValue, m)
		}

		returnBytes, _ := json.Marshal(returnValue)

		client.sending <- WebSocketMessage{Data: returnBytes}

		return
	}

}

func BroadcastToOutput(msg WebSocketMessage) {
	for client, b := range hubs[constants.OUTPUT_WEBSOCKET_NAME].clients {
		if !b {
			continue
		}
		msg.Datetime = time.Now()
		client.sending <- msg
	}
}
