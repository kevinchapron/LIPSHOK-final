package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/internal-connectors"
	"github.com/kevinchapron/FSHK-final/messaging"
	"github.com/kevinchapron/FSHK-final/security"
	"net"
)

var TCPmessagesInternal = make(chan messaging.Message)

func main() {
	Logging.SetLevel(Logging.DEBUG)
	Logging.Info("Program started. Trying to connect internal websocket ...")

	// created internal connection with main app.
	internConnector := internal_connectors.CreateInternalConnector()
	go internConnector.Connect(&TCPmessagesInternal)
	if a := <-internConnector.IsConnected; !a {
		Logging.Error("Problem while trying to connect.")
		return
	}
	Logging.Info("Successfully connected... Starting TCP socket ...")

	// waiting for the TCP server to be up.
	addr := net.TCPAddr{
		Port: constants.TCP_PORT,
		IP:   net.ParseIP(constants.TCP_ADDR),
	}
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		Logging.Error(err)
		return
	}

	defer listener.Close()
	Logging.Info("Waiting for TCP sensor connections ...")

	for {
		c, err := listener.Accept()
		if err != nil {
			Logging.Error(err)
			continue
		}

		go func() {
			var p = make([]byte, constants.MAX_TCP_PACKET_SIZE)
			for {
				_, err := bufio.NewReader(c).Read(p)
				if err != nil {
					Logging.Error(err)
					break
				}
				go receivedTCPMessage(c, p)
			}
		}()
	}
}

func receivedTCPMessage(conn net.Conn, data []byte) {
	var m messaging.Message
	err := m.FromBytes(data)
	if err != nil {
		Logging.Error(err)
		return
	}
	Logging.Info(fmt.Sprintf("Received Message from %s : %s", conn.RemoteAddr().String(), m.Data))
	// forward message to main app
	// change AES IV
	m.AesIV = security.RandomKey()
	TCPmessagesInternal <- m

	answered := messaging.AnswerMessage{Data: "OK"}
	m.Data, err = json.Marshal(answered)
	if err != nil {
		Logging.Error(err)
		return
	}

	m.AesIV = security.RandomKey()
	m.DataType = constants.MESSAGING_DATATYPE_DATA
	n, err := conn.Write(m.ToBytes())
	if err != nil {
		Logging.Error(err)
		return
	}
	if n == 0 {
		Logging.Error("Nothing could be sent...")
		return
	}
}
