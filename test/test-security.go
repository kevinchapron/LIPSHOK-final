package main

import (
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"github.com/kevinchapron/LIPSHOK/security"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)

	var m messaging.Message
	m.AesIV = security.RandomKey()
	m.DataType = constants.MESSAGING_DATATYPE_DATA
	m.From = "TestingSensor-UDP"
	m.Data = []byte("{\"value\":1234}")

	var bytes = m.ToBytes()

	Logging.Info("Sent:", m.From, ";", string(m.Data))
	Logging.Info("Sent:", bytes)

	var m2 messaging.Message
	err := m2.FromBytes(bytes)
	if err != nil {
		Logging.Error(err)
		return
	}

	Logging.Info("Received:", bytes)
	Logging.Info("Received:", m2.From, string(m2.Data))

}
