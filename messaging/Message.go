package messaging

import (
	"encoding/binary"
	"errors"
	"github.com/kevinchapron/FSHK-final/security"
)

type Message struct {
	AesIV    [12]byte
	DataType byte

	Data []byte
}

func (m *Message) FromBytes(s []byte) error {
	if len(s) <= 24 {
		return errors.New("Message not long enough")
	}
	for i, v := range s[2 : len(m.AesIV)+2] {
		m.AesIV[i] = v
	}
	m.DataType = s[14]

	lengthData := binary.LittleEndian.Uint16(s[:2])
	encryptedData := s[24 : lengthData+24]
	msg, err := security.Decrypt(encryptedData, m.AesIV[:])
	if err != nil {
		return err
	}
	for _, v := range msg {
		m.Data = append(m.Data, v)
	}
	return nil
}

func (m *Message) ToBytes() []byte {
	var r = make([]byte, 24)

	encryptedData, _ := security.Encrypt(m.Data, m.AesIV[:])

	binary.LittleEndian.PutUint16(r[:2], uint16(len(encryptedData)))
	for i, v := range m.AesIV {
		r[2+i] = v
	}
	r[14] = m.DataType

	for _, v := range encryptedData {
		r = append(r, v)
	}

	return r
}
