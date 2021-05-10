package security

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
)

const (
	BLE_ENCRYPTION_MASTERKEY_SIZE = 10

	BLE_ENCRYPTION_PRESENT_BLOCK_SIZE = 64
	BLE_ENCRYPTION_PRESENT_KEY_SIZE   = 80
	BLE_ENCRYPTION_PRESENT_ROUNDS     = 31

	BLE_ENCRYPTION_PRESENT_MASTERKEY_SIZE = (BLE_ENCRYPTION_PRESENT_KEY_SIZE / 8 * (BLE_ENCRYPTION_PRESENT_ROUNDS + 1))
)

var _sBox4 = []byte{0xC, 0x5, 0x6, 0xB, 0x9, 0x0, 0xA, 0xD, 0x3, 0xE, 0xF, 0x8, 0x4, 0x7, 0x1, 0x2}
var _sInvBox4 = []byte{0x5, 0xE, 0xF, 0x8, 0xC, 0x1, 0x2, 0xD, 0xB, 0x4, 0x6, 0x3, 0x0, 0x7, 0x9, 0xA}

func ror64(x uint64, n uint) uint64 { return x>>n | x<<(64-n) }
func rol32(x uint32, n uint) uint32 { return x<<n | x>>(32-n) }
func ror32(x uint32, n uint) uint32 { return x>>n | x<<(32-n) }
func reverse(arr *[]uint8, n int) {
	var high, low uint8
	low = 0
	high = uint8(n - 1)
	for low < high {
		var tmp uint8 = (*(arr))[low]
		(*(arr))[low] = (*(arr))[high]
		(*(arr))[high] = tmp

		low++
		high--
	}
}

var _ble_encryption_instance = make(map[string]*BLEEncryption)

type BLEEncryption struct {
	masterKey []byte
	fullKeys  []byte

	keyHigh []byte
	keyLow  []byte
}

func GetBLEEncryptionSystem(addr string) *BLEEncryption {
	instance, ok := _ble_encryption_instance[addr]
	if !ok {
		_ble_encryption_instance[addr] = &BLEEncryption{
			masterKey: make([]byte, BLE_ENCRYPTION_MASTERKEY_SIZE),
			fullKeys:  make([]byte, BLE_ENCRYPTION_PRESENT_MASTERKEY_SIZE),
			keyHigh:   make([]byte, 8),
			keyLow:    make([]byte, 8),
		}
		instance = _ble_encryption_instance[addr]
	}
	return instance
}

func (system *BLEEncryption) setMasterKey(bytes_key []byte) {
	system.masterKey = bytes_key

	var key uint64 = binary.LittleEndian.Uint64(bytes_key)

	var keyhigh uint64
	var keylow uint64

	keylow = key

	var highBytes uint16 = binary.LittleEndian.Uint16([]byte{bytes_key[8], bytes_key[9]})
	keyhigh = (uint64(highBytes)<<48 | (keylow >> 16))

	binary.LittleEndian.PutUint64(system.fullKeys[:8], keyhigh)

	var temp uint64
	var i uint8

	var index_min uint
	var index_max uint

	for i = 0; i < BLE_ENCRYPTION_PRESENT_ROUNDS; i++ {
		//	/* 61-bit left shift */
		temp = keyhigh
		keyhigh <<= 61
		keyhigh |= (keylow << 45)
		keyhigh |= (temp >> 19)
		keylow = (temp >> 3) & 0xFFFF

		//	/* S-Box application */
		temp = uint64(_sBox4[keyhigh>>60])
		//	temp = sbox[keyhigh >> 60];

		keyhigh &= 0x0FFFFFFFFFFFFFFF
		keyhigh |= temp << 60

		//	/* round counter addition */
		keylow ^= (((uint64)(i+1) & 0x01) << 15)
		keyhigh ^= ((uint64)(i+1) >> 1)
		//

		index_min = uint((i + 1) * 8)
		index_max = uint((i+1)*8) + 8

		binary.LittleEndian.PutUint64(system.fullKeys[index_min:index_max], keyhigh)
	}

	binary.LittleEndian.PutUint64(system.keyHigh, keyhigh)
	binary.LittleEndian.PutUint64(system.keyLow, keylow)
}

func (system *BLEEncryption) GenerateMasterKey() []byte {
	var key = make([]byte, BLE_ENCRYPTION_MASTERKEY_SIZE)
	rand.Read(key)
	system.setMasterKey(key)
	return key
}

func (system *BLEEncryption) GetMasterKey() []byte {
	return system.masterKey
}

func (system *BLEEncryption) Encrypt(b []byte) ([]byte, error) {
	if len(b) != 20 && len(b) != 8 {
		return nil, errors.New("Must be for 20-bytes array or 8-bytes array")
	}
	var d []byte = make([]byte, len(b))
	copy(d[0:], b)

	if len(b) == 20 {
		copy(d[:8], system.encrypt(d[:8]))
		reverse(&d, len(d))
		copy(d[:8], system.encrypt(d[:8]))
		reverse(&d, len(d))

		copy(d[6:14], system.encrypt(d[6:14]))
	} else {
		copy(d[0:], system.encrypt(d))
	}

	return d, nil
}

func (system *BLEEncryption) Decrypt(b []byte) ([]byte, error) {
	if len(b) != 20 && len(b) != 15 && len(b) != 8 {
		return nil, errors.New("Must be for a bytes-array of size 8,15 or 20.")
	}
	var d = make([]byte, len(b))
	copy(d[0:], b)

	if len(b) > 16 {
		copy(d[6:14], system.decrypt(d[6:14]))
		reverse(&d, len(d))
		copy(d[:8], system.decrypt(d[:8]))
		reverse(&d, len(d))
		copy(d[:8], system.decrypt(d[:8]))
	}
	if len(b) > 8 && len(b) <= 16 {
		reverse(&d, len(d))
		copy(d[:8], system.decrypt(d[:8]))
		reverse(&d, len(d))
		copy(d[:8], system.decrypt(d[:8]))
	}
	if len(b) <= 8 {
		copy(d[0:], system.decrypt(d))
	}
	return d, nil
}

func (system *BLEEncryption) GetFullKeys() []byte {
	return system.fullKeys
}

func (system *BLEEncryption) encrypt(b []byte) []byte {
	var state uint64 = binary.LittleEndian.Uint64(b)

	var rk []uint64 = make([]uint64, len(system.fullKeys)/8)
	for i := 0; i < len(rk); i++ {
		rk[i] = binary.LittleEndian.Uint64(system.fullKeys[i*8 : i*8+8])
	}

	var result uint64
	var sInput uint8
	var pLayerIndex uint8
	var stateBit uint64
	var i uint8
	var k uint16

	for i = 0; i < BLE_ENCRYPTION_PRESENT_ROUNDS; i++ {
		state ^= rk[i]

		//	/* sbox */
		for k = 0; k < BLE_ENCRYPTION_PRESENT_BLOCK_SIZE/4; k++ {
			sInput = uint8(state & 0x0F)
			state &= 0xFFFFFFFFFFFFFFF0
			state |= uint64(_sBox4[sInput])
			state = ror64(state, 4)
		}

		//	/* pLayer */
		result = 0
		for k = 0; k < BLE_ENCRYPTION_PRESENT_BLOCK_SIZE; k++ {
			stateBit = state & 0x1
			state = state >> 1

			if 0 != stateBit {
				pLayerIndex = uint8((16 * k) % 63)
				if 63 == k {
					pLayerIndex = 63
				}
				result |= stateBit << pLayerIndex
			}
		}
		state = result
	}

	state ^= rk[i]
	var return_value = make([]byte, len(b))
	binary.LittleEndian.PutUint64(return_value, state)
	return return_value
}

func (system *BLEEncryption) decrypt(b []byte) []byte {

	var state uint64 = binary.LittleEndian.Uint64(b)

	var rk []uint64 = make([]uint64, len(system.fullKeys)/8)
	for i := 0; i < len(rk); i++ {
		rk[i] = binary.LittleEndian.Uint64(system.fullKeys[i*8 : i*8+8])
	}

	var result uint64
	var sInput uint8
	var pLayerIndex uint8
	var stateBit uint64

	var i uint8
	var k uint16

	for i = BLE_ENCRYPTION_PRESENT_ROUNDS; i > 0; i-- {
		state ^= rk[i]

		/* pLayer */
		result = 0
		for k = 0; k < BLE_ENCRYPTION_PRESENT_BLOCK_SIZE; k++ {
			stateBit = state & 0x1
			state = state >> 1
			if 0 != stateBit {
				pLayerIndex = uint8((4 * k) % 63)
				if 63 == k {
					pLayerIndex = 63
				}
				result |= stateBit << pLayerIndex
			}
		}
		state = result
		/* sbox */
		for k = 0; k < BLE_ENCRYPTION_PRESENT_BLOCK_SIZE/4; k++ {
			sInput = uint8(state & 0xF)
			state &= 0xFFFFFFFFFFFFFFF0
			state |= uint64(_sInvBox4[sInput])
			state = ror64(state, 4)
		}
	}
	state ^= rk[i]

	var return_value = make([]byte, len(b))
	binary.LittleEndian.PutUint64(return_value, state)
	return return_value
}
