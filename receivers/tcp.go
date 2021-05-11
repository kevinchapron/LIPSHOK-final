package receivers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/internal-connectors"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"github.com/kevinchapron/LIPSHOK/security"
	"net"
)

func CreateTCPReceiver(logPrefix string) {

	// created internal connection with main app.
	internConnector := internal_connectors.CreateInternalConnector("TCP Manager", "TCP")
	go internConnector.Connect(&TCPmessagesInternal)
	if a := <-internConnector.IsConnected; !a {
		Logging.Error(logPrefix, "Problem while trying to connect.")
		return
	}
	Logging.Info(logPrefix, "Successfully connected... Starting TCP socket ...")

	// waiting for the TCP server to be up.
	addr := net.TCPAddr{
		Port: constants.TCP_PORT,
		IP:   net.ParseIP(constants.TCP_ADDR),
	}
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		Logging.Error(logPrefix, err)
		return
	}

	defer listener.Close()
	Logging.Info(logPrefix, fmt.Sprintf("Waiting for TCP sensor connections on port %d...", constants.TCP_PORT))

	for {
		c, err := listener.Accept()
		if err != nil {
			Logging.Error(logPrefix, err)
			continue
		}

		go func() {
			var p = make([]byte, constants.MAX_TCP_PACKET_SIZE)
			for {
				_, err := bufio.NewReader(c).Read(p)
				if err != nil {
					Logging.Error(logPrefix, err)
					break
				}
				go receivedTCPMessage(logPrefix, c, p)
			}
		}()
	}
}

func receivedTCPMessage(logPrefix string, conn net.Conn, data []byte) {
	var m messaging.Message
	err := m.FromBytes(data)
	if err != nil {
		Logging.Error(logPrefix, err)
		return
	}
	Logging.Info(logPrefix, fmt.Sprintf("Received Message from %s : %s", conn.RemoteAddr().String(), m.Data))
	// forward message to main app
	// change AES IV
	m.AesIV = security.RandomKey()
	TCPmessagesInternal <- m

	answered := messaging.AnswerMessage{Data: "OK"}
	m.Data, err = json.Marshal(answered)
	if err != nil {
		Logging.Error(logPrefix, err)
		return
	}

	m.AesIV = security.RandomKey()
	m.DataType = constants.MESSAGING_DATATYPE_DATA
	n, err := conn.Write(m.ToBytes())
	if err != nil {
		Logging.Error(logPrefix, err)
		return
	}
	if n == 0 {
		Logging.Error(logPrefix, "Nothing could be sent...")
		return
	}
}
