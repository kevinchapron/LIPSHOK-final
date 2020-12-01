package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/messaging"
	"github.com/kevinchapron/FSHK-final/security"
)

type CalibrationStep struct {
	Description string
	Value       float64
	Completed   bool
}

type Device struct {
	Name             string
	ConnectionType   string
	CalibrationSteps []CalibrationStep
}

func main() {
	Logging.SetLevel(Logging.DEBUG)
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:5001/sensors", nil)
	if err != nil {
		Logging.Error(err)
		return
	}
	defer conn.Close()

	go readData(conn)
	var m = messaging.Message{
		DataType: 0x01,
		AesIV:    security.RandomKey(),
	}

	device := Device{
		Name:           "TestingSensor",
		ConnectionType: "Wi-Fi",
		CalibrationSteps: []CalibrationStep{
			{Description: "Itération 1", Value: -1, Completed: false},
			{Description: "Itération 2", Value: -1, Completed: false},
			{Description: "Itération 3", Value: -1, Completed: false},
		},
	}

	m.Data, err = json.Marshal(device)
	if err != nil {
		Logging.Error(err)
		return
	}

	Logging.Info(m.ToBytes())
	err = conn.WriteMessage(websocket.BinaryMessage, m.ToBytes())
	if err != nil {
		Logging.Error(err)
		return
	}

	select {}
}
func readData(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			Logging.Error(err)
			continue
		}
		Logging.Info("Received Data:", msg)
	}
}
