package sensor_simulation

import (
	"bufio"
	"encoding/json"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"net"
)

func AskAndAnswer(conn net.Conn, p []byte) (*messaging.AnswerMessage, error) {
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
