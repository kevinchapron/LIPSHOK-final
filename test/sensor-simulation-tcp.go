package main

import (
	"encoding/json"
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"github.com/kevinchapron/LIPSHOK/test/sensor-simulation"
	"net"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", constants.TCP_ADDR, constants.TCP_PORT))
	if err != nil {
		Logging.Error(err)
		return
	}
	defer conn.Close()

	// Default comportement
	device := sensor_simulation.GetDevice()
	device.Protocol = "TCP"
	data, err := sensor_simulation.ObjectToBytesAuth(device)
	if err != nil {
		Logging.Error(err)
		return
	}

	p := make([]byte, constants.MAX_TCP_PACKET_SIZE)
	// send depending to the conn.
	n, err := conn.Write(data)
	if err != nil {
		Logging.Error(err)
		return
	}
	Logging.Info("Sent data : ", n)
	_, err = conn.Read(p)
	if err != nil {
		Logging.Error(err)
		return
	}
	var m messaging.Message
	m.FromBytes(p)

	var answer messaging.AnswerMessage
	json.Unmarshal(m.Data, &answer)

	Logging.Info("Received : ", answer)

	select {}
}
