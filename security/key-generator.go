package security

import (
	"math/rand"
	"time"
)

func RandomKey() [12]byte {
	var r [12]byte
	rand.Seed(time.Now().UnixNano())
	for i, _ := range r {
		r[i] = (byte)(rand.Intn(0xFF))
	}
	return r
}
