package websockets

import (
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/messaging"
)

var WebSocketFunction = map[string]func(messaging.Message, *WebSocketClient, *WebSocketHub){
	constants.SENSOR_WEBSOCKET_NAME: SensorSettings,
}

func SensorSettings(msg messaging.Message, client *WebSocketClient, hub *WebSocketHub) {

	Logging.Debug("Received: ", msg, "; From:", client.conn.RemoteAddr().String())
	Logging.Debug("--->", string(msg.Data))

}
