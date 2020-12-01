package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/main-app/constants"
	"github.com/kevinchapron/FSHK-final/main-app/database"
	"github.com/kevinchapron/FSHK-final/main-app/websockets"
	"net/http"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)

	db := database.GetDatabase()
	db.Connect()

	deviceTypes := db.GetDeviceTypes()

	Logging.Debug("The database is connected.")

	var done = make(chan bool)
	muxRouter := mux.NewRouter()

	for _, deviceType := range deviceTypes {
		if deviceType.Name == database.DATABASE_DEVICE_CONNECTION_WIFI && deviceType.Activated {
			// Wi-Fi activated; we link the websockets.
			websockets.CreateWebSocket(constants.WEBSOCKET_NAME, "/sensors", muxRouter)
			continue
		}

		continue
	}

	go http.ListenAndServe(fmt.Sprintf("%s:%d", constants.ROUTER_ADDRESS, constants.ROUTER_PORT), muxRouter)
	Logging.Info(fmt.Sprintf("[WS] > WebSocket application listening on ws://%s:%d/", constants.ROUTER_ADDRESS, constants.ROUTER_PORT))

	Logging.Info("Program launched.")
	<-done
	Logging.Info("Program terminated.")
}
