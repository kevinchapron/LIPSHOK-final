package receivers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/internal-connectors"
	"github.com/kevinchapron/LIPSHOK/main-app/websockets"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"net/http"
)

func CreateBLEReceiver(logPrefix string) {
	Logging.Info(logPrefix, "Creating BLE Receiver")

	// created internal connection with main app.
	internConnector := internal_connectors.CreateInternalConnector("BLE Manager", "BLE")
	go internConnector.Connect(&BLEMessagesInternal)
	if a := <-internConnector.IsConnected; !a {
		Logging.Error(logPrefix, "Problem while trying to connect.")
		return
	}
	Logging.Info(logPrefix, "Successfully connected... Starting BLE socket ...")

	// Creation of the inner sensor receiver websocket
	if err := websockets.RegisterNewFunction(constants.WEBSOCKET_INNER_BLE_NAME, receivedBLEMessages); err != nil {
		Logging.Error("Problem with function adding")
		return
	}
	muxRouter := mux.NewRouter()
	websockets.CreateWebSocket(constants.WEBSOCKET_INNER_BLE_NAME, constants.WEBSOCKET_INNER_BLE_PATH, muxRouter)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", constants.WEBSOCKET_INNER_BLE_ADDR, constants.WEBSOCKET_INNER_BLE_PORT), muxRouter)
	Logging.Info(logPrefix, fmt.Sprintf("BLE WebSockets listening on ws://%s:%d%s", constants.ROUTER_ADDRESS, constants.WEBSOCKET_INNER_BLE_PORT, constants.WEBSOCKET_INNER_BLE_PATH))
}

func receivedBLEMessages(message messaging.Message, client *websockets.WebSocketClient, hub *websockets.WebSocketHub) {
	//var data messaging.InertialData
	//json.Unmarshal(message.Data, &data)

	// TODO: add MAC addr to standard messaging.

	websockets.BroadcastToOutput(websockets.WebSocketMessage{
		Data: message.Data,
		From: client,
	})
}
