package main

import (
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/internal-connectors"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)
	Logging.Info("Program started. Trying to connect internal websocket ...")

	internConnector := internal_connectors.CreateInternalConnector()
	go internConnector.Connect()

	if a := <-internConnector.IsConnected; !a {
		Logging.Error("Problem while trying to connect.")
		return
	}
	Logging.Info("Successfully connected.")
	Logging.Info("Waiting for UDP sensor connections ...")

	select {}
}
