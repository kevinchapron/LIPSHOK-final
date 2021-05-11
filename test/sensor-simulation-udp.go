package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"github.com/kevinchapron/LIPSHOK/test/sensor-simulation"
	"math/rand"
	"net"
	"time"
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

	authAnswer, err := askAndAnswer(conn, data)
	if err != nil {
		Logging.Error(err)
		return
	}
	Logging.Debug("Answer:", authAnswer)

	m := map[string]interface{}{
		"value": 0,
	}

	for {
		time.Sleep(time.Second)

		m["value"] = rand.Intn(1000)

		p, _ := sensor_simulation.ObjectToBytesData(m)

		answer, err := askAndAnswer(conn, p)
		if err != nil {
			Logging.Error(err)
			continue
		}
		Logging.Debug("Sent:", m)
		Logging.Debug("Received:", answer)
	}
	select {}
}

func askAndAnswer(conn net.Conn, p []byte) (*messaging.AnswerMessage, error) {
	// send depending to the conn.
	_, err := conn.Write(p)
	if err != nil {
		return nil, err
	}

	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		return nil, err
	}
	var m messaging.Message
	m.FromBytes(p)

	var answer messaging.AnswerMessage
	json.Unmarshal(m.Data, &answer)
	return &answer, nil
}
