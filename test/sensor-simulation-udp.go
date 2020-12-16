package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/messaging"
	"github.com/kevinchapron/FSHK-final/test/sensor-simulation"
	"net"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", constants.UDP_ADDR, constants.UDP_PORT))
	if err != nil {
		Logging.Error(err)
		return
	}
	defer conn.Close()

	// Default comportement
	device := sensor_simulation.GetDevice()
	data, err := sensor_simulation.ObjectToBytesAuth(device)
	if err != nil {
		Logging.Error(err)
		return
	}

	p := make([]byte, constants.MAX_UDP_PACKET_SIZE)
	// send depending to the conn.
	n, err := conn.Write(data)
	if err != nil {
		Logging.Error(err)
		return
	}
	Logging.Info("Sent data : ", n)
	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		Logging.Error(err)
		return
	}
	var m messaging.Message
	m.FromBytes(p)

	var answer messaging.AnswerMessage
	json.Unmarshal(m.Data, &answer)

	select {}
}
