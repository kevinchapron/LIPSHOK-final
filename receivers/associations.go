package receivers

var ProtocolFunctions = map[string]func(string){
	"UDP":   CreateUDPReceiver,
	"TCP":   CreateTCPReceiver,
	"BLE":   CreateBLEReceiver,
	"ZWAVE": CreateBLEReceiver,
}
