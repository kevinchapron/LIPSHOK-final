package messaging

import (
	"encoding/binary"
	"errors"
	"github.com/kevinchapron/LIPSHOK/security"
)

type Message struct {
	AesIV    [12]byte
	DataType byte
	Data     []byte

	From string
}

func (m *Message) FromBytes(s []byte) error {
	if len(s) <= 24 {
		return errors.New("Message not long enough")
	}
	for i, v := range s[2 : len(m.AesIV)+2] {
		m.AesIV[i] = v
	}
	m.DataType = s[14]

	lengthFrom := uint16(s[15])
	if lengthFrom != 0 {
		encryptedFrom := s[24 : lengthFrom+24]
		sFrom, err := security.Decrypt(encryptedFrom, m.AesIV[:])
		if err != nil {
			return err
		}
		m.From = string(sFrom)
	}

	lengthData := binary.LittleEndian.Uint16(s[:2])
	encryptedData := s[24+lengthFrom : lengthData+24]
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
	encryptedFrom, _ := security.Encrypt([]byte(m.From), m.AesIV[:])

	binary.LittleEndian.PutUint16(r[:2], uint16(len(encryptedData)+len(encryptedFrom)))
	for i, v := range m.AesIV {
		r[2+i] = v
	}
	r[14] = m.DataType
	r[15] = byte(len(encryptedFrom))

	r = append(r, encryptedFrom...)
	r = append(r, encryptedData...)

	return r
}
