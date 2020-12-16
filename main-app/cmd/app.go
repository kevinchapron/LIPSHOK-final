package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/main-app/database"
	"github.com/kevinchapron/FSHK-final/main-app/websockets"
	"github.com/kevinchapron/FSHK-final/receivers"
	"io/ioutil"
	"net/http"
	"os"
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

	go startProtocols()

	<-done
	Logging.Info("Program terminated.")
}

func startProtocols() {

	// reading the protocols to know if it musts start its own
	var protocols database.Protocols
	jsonFile, err := os.Open(constants.PROTOCOLS_JSON_CONF)
	if err != nil && !os.IsNotExist(err) {
		Logging.Error(err)
		return
	}
	if os.IsNotExist(err) {
		Logging.Warning(fmt.Sprintf("Trying to read \"%s\", but no file found.", constants.PROTOCOLS_JSON_CONF))
		return
	}

	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(bytes, &protocols)
	if err != nil {
		Logging.Error(err)
		return
	}

	for _, protocol := range protocols.List {
		if protocol.Activated {
			f, ok := receivers.ProtocolFunctions[protocol.Name]
			if !ok {
				Logging.Warning(fmt.Sprintf("NO FUNCTION FOR PROTOCOL : ---%s---", protocol.Name))
				continue
			}
			go f(fmt.Sprintf("[%s]", protocol.Name))
		}
	}
}
