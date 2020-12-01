package security

import (
	"crypto/aes"
	"crypto/cipher"
)

var masterAesKey = []byte{149, 165, 159, 33, 124, 170, 183, 125, 140, 42, 131, 253, 15, 81, 142, 126}

func Encrypt(msg []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(masterAesKey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	bytes := aesgcm.Seal(nil, iv, msg, nil)
	return bytes, nil
}
func Decrypt(msg []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(masterAesKey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	bytes, err := aesgcm.Open(nil, iv, msg, nil)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
