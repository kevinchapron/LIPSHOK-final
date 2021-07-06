package receivers

import (
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	internal_connectors "github.com/kevinchapron/LIPSHOK/internal-connectors"
	"github.com/stampzilla/gozwave"
)

func CreateZWaveReceiver(logPrefix string) {

	// created internal connection with main app.
	internConnector := internal_connectors.CreateInternalConnector("Z-Wave Manager", "ZWave")
	go internConnector.Connect(&ZWaveMessagesInternal)
	if a := <-internConnector.IsConnected; !a {
		Logging.Error(logPrefix, "Problem while trying to connect.")
		return
	}
	Logging.Info(logPrefix, "Successfully connected... Starting Z-Wave socket ...")
	z, err := gozwave.Connect(constants.ZWAVE_INPUT_PORT, "")
	if err != nil {
		Logging.Error(err)
		return
	}
}
