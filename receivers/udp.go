package receivers

import (
	"encoding/json"
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/internal-connectors"
	"github.com/kevinchapron/FSHK-final/messaging"
	"github.com/kevinchapron/FSHK-final/security"
	"net"
)

func CreateUDPReceiver(logPrefix string) {
	// created internal connection with main app.
	internConnector := internal_connectors.CreateInternalConnector()
	go internConnector.Connect(&UDPmessagesInternal)
	if a := <-internConnector.IsConnected; !a {
		Logging.Error(logPrefix, "Problem while trying to connect.")
		return
	}
	Logging.Info(logPrefix, "Successfully connected... Starting UDP socket ...")

	// waiting for the UDP server to be up.
	addr := net.UDPAddr{
		Port: constants.UDP_PORT,
		IP:   net.ParseIP(constants.UDP_ADDR),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		Logging.Error(logPrefix, err)
		return
	}
	Logging.Info(logPrefix, "Waiting for UDP sensor connections ...")

	var p = make([]byte, constants.MAX_UDP_PACKET_SIZE)
	for {
		_, remoteAddr, err := ser.ReadFromUDP(p)
		if err != nil {
			Logging.Error(logPrefix, err)
			continue
		}
		go receivedUDPMessage(logPrefix, ser, remoteAddr, p)
	}
}

func receivedUDPMessage(logPrefix string, conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	var m messaging.Message
	err := m.FromBytes(data)
	if err != nil {
		Logging.Error(logPrefix, err)
		return
	}
	Logging.Info(logPrefix, fmt.Sprintf("Received Message from %s : %s", addr.String(), m.Data))
	// forward message to main app
	// change AES IV
	m.AesIV = security.RandomKey()
	UDPmessagesInternal <- m

	answered := messaging.AnswerMessage{Data: "OK"}
	m.Data, err = json.Marshal(answered)
	if err != nil {
		Logging.Error(logPrefix, err)
		return
	}

	m.AesIV = security.RandomKey()
	m.DataType = constants.MESSAGING_DATATYPE_DATA
	_, err = conn.WriteToUDP(m.ToBytes(), addr)
	if err != nil {
		Logging.Error(logPrefix, err)
		return
	}
}
