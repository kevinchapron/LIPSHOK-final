package main

import (
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/test/sensor-simulation"
	"math/rand"
	"net"
	"time"
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
	device := sensor_simulation.GetDevice("Sensor on TCP", "TCP")
	data, err := sensor_simulation.ObjectToBytesAuth(device)
	if err != nil {
		Logging.Error(err)
		return
	}

	authAnswer, err := sensor_simulation.AskAndAnswer(conn, data)
	if err != nil {
		Logging.Error(err)
		return
	}
	Logging.Debug("Answer:", authAnswer)

	values := map[string]interface{}{
		"value": 0,
	}

	for {
		time.Sleep(time.Second)

		values["value"] = rand.Intn(1000)

		p, _ := sensor_simulation.ObjectToBytesData(values)

		answer, err := sensor_simulation.AskAndAnswer(conn, p)
		if err != nil {
			Logging.Error(err)
			continue
		}
		Logging.Debug("Sent:", values)
		Logging.Debug("Received:", answer)
	}
	select {}
}
