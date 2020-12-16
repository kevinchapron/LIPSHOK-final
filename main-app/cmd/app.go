package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/main-app/database"
	"github.com/kevinchapron/FSHK-final/main-app/websockets"
	"net/http"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)

	db := database.GetDatabase()
	db.Connect()

	Logging.Debug("The database is connected.")

	var done = make(chan bool)
	muxRouter := mux.NewRouter()
	websockets.CreateWebSocket(constants.SENSOR_WEBSOCKET_NAME, "/sensors", muxRouter)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", constants.ROUTER_ADDRESS, constants.SENSOR_WEBSOCKET_PORT), muxRouter)
	Logging.Info(fmt.Sprintf("[WS] > WebSockets application listening on ws://%s:%d/", constants.ROUTER_ADDRESS, constants.SENSOR_WEBSOCKET_PORT))

	routerInterface := mux.NewRouter()
	// Add of the graphic interface-dev
	routerInterface.PathPrefix("/").Handler(http.FileServer(http.Dir("../interface-build/")))
	go http.ListenAndServe(fmt.Sprintf("%s:%d", constants.ROUTER_ADDRESS, constants.INTERFACE_PORT), routerInterface)

	Logging.Info("Program launched.")

	// reading the protocols to know if it musts start its own

	<-done
	Logging.Info("Program terminated.")
}
