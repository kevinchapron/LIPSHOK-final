package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/main-app/database"
	"github.com/kevinchapron/LIPSHOK/main-app/websockets"
	"github.com/kevinchapron/LIPSHOK/receivers"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)

	//db := database.GetDatabase()
	//db.Connect()
	//
	//Logging.Debug("The database is connected.")

	var done = make(chan bool)

	// Creation of the inner sensor receiver websocket
	muxRouter := mux.NewRouter()
	websockets.CreateWebSocket(constants.SENSOR_WEBSOCKET_NAME, constants.SENSOR_WEBSOCKET_PATH, muxRouter)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", constants.ROUTER_ADDRESS, constants.SENSOR_WEBSOCKET_PORT), muxRouter)
	Logging.Info(fmt.Sprintf("[WS] > WebSockets application listening on ws://%s:%d%s", constants.ROUTER_ADDRESS, constants.SENSOR_WEBSOCKET_PORT, constants.SENSOR_WEBSOCKET_PATH))

	// Creation of the output websocket (no encryption)
	muxRouter2 := mux.NewRouter()
	websockets.CreateRawWebSocket(constants.OUTPUT_WEBSOCKET_NAME, constants.OUTPUT_WEBSOCKET_PATH, muxRouter2)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", constants.ROUTER_ADDRESS, constants.OUTPUT_WEBSOCKET_PORT), muxRouter2)
	Logging.Info(fmt.Sprintf("[WS] > WebSockets output listening on ws://%s:%d%s", constants.ROUTER_ADDRESS, constants.OUTPUT_WEBSOCKET_PORT, constants.OUTPUT_WEBSOCKET_PATH))

	// Add of the graphic interface-dev
	routerInterface := mux.NewRouter()
	routerInterface.PathPrefix("/").Handler(http.FileServer(http.Dir("../interface-build/")))
	go http.ListenAndServe(fmt.Sprintf("%s:%d", constants.ROUTER_ADDRESS, constants.INTERFACE_PORT), routerInterface)

	Logging.Info("Program launched.")

	go startProtocols()

	time.Sleep(time.Second * 5)
	websockets.PrintConnectorsList(websockets.ListAllConnectors())
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
