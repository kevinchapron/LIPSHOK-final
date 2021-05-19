package main

import (
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/test/sensor-simulation"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)

	const NB_SENSORS_SIMULATED = 10
	var mutex sync.Mutex

	for i := 0; i <= NB_SENSORS_SIMULATED; i++ {
		go func(delay int) {
			time.Sleep(time.Duration(5*delay) * time.Second)
			mutex.Lock()
			conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", constants.UDP_ADDR, constants.UDP_PORT))
			if err != nil {
				Logging.Error(err)
				return
			}
			mutex.Unlock()
			defer conn.Close()

			// Default comportement
			device := sensor_simulation.GetDevice("Sensor on UDP #"+strconv.Itoa(rand.Intn(500)), "UDP")

			data, err := sensor_simulation.ObjectToBytesAuth(device)
			if err != nil {
				Logging.Error(err)
				return
			}

			mutex.Lock()
			authAnswer, err := sensor_simulation.AskAndAnswer(conn, data)
			mutex.Unlock()
			if err != nil {
				Logging.Error(err)
				return
			}
			Logging.Debug("Sent:", string(data))
			Logging.Debug("Answer:", authAnswer)

			m := map[string]interface{}{
				"value": 0,
			}

			for {
				time.Sleep(time.Second)

				m["value"] = rand.Intn(1000)

				p, _ := sensor_simulation.ObjectToBytesData(m)

				mutex.Lock()
				answer, err := sensor_simulation.AskAndAnswer(conn, p)
				if err != nil {
					Logging.Error(err)
					continue
				}
				mutex.Unlock()
				Logging.Debug("Sent:", m)
				Logging.Debug("Received:", answer)
			}
		}(i)
	}

	select {}
}
