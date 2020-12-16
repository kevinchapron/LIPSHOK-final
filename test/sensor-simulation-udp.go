package main

import (
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
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

	// send depending to the conn.
	n, err := conn.Write(data)
	if err != nil {
		Logging.Error(err)
		return
	}
	Logging.Info("Sent data : ", n)

	select {}
}
