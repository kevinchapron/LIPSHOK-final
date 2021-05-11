package main

import (
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/receivers"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)
	Logging.Info("Program started. Trying to connect internal websocket ...")

	receivers.CreateTCPReceiver("[TCP]")
}
