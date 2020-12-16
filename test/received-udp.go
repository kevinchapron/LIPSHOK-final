package main

import (
	"fmt"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/internal-connectors"
	"github.com/kevinchapron/FSHK-final/messaging"
	"net"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)
	Logging.Info("Program started. Trying to connect internal websocket ...")

	// created internal connection with main app.
	internConnector := internal_connectors.CreateInternalConnector()
	go internConnector.Connect()
	if a := <-internConnector.IsConnected; !a {
		Logging.Error("Problem while trying to connect.")
		return
	}
	Logging.Info("Successfully connected... Starting UDP socket ...")

	// waiting for the UDP server to be up.
	addr := net.UDPAddr{
		Port: constants.UDP_PORT,
		IP:   net.ParseIP(constants.UDP_ADDR),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		Logging.Error(err)
		return
	}
	Logging.Info("Waiting for UDP sensor connections ...")
	var p = make([]byte, constants.MAX_UDP_PACKET_SIZE)
	for {
		_, remoteAddr, err := ser.ReadFromUDP(p)
		if err != nil {
			Logging.Error(err)
			continue
		}
		go receivedMessage(ser, remoteAddr, p)
	}
}

func receivedMessage(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	var m messaging.Message
	err := m.FromBytes(data)
	if err != nil {
		Logging.Error(err)
		return
	}
	Logging.Info(fmt.Sprintf("Received Message from %s : %s", addr.String(), m.Data))

	//conn.WriteToUDP(m.ToBytes(),addr)
}
