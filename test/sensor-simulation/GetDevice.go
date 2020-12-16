package sensor_simulation

import (
	"encoding/json"
	"github.com/kevinchapron/FSHK-final/messaging"
	"github.com/kevinchapron/FSHK-final/security"
)

func GetDevice() *Device {
	return &Device{
		Name:           "TestingSensor",
		ConnectionType: "UDP",
		CalibrationSteps: []CalibrationStep{
			{Description: "Itération 1", Value: -1, Completed: false},
			{Description: "Itération 2", Value: -1, Completed: false},
			{Description: "Itération 3", Value: -1, Completed: false},
		},
	}
}

func objectToBytes(a interface{}, dataType uint) ([]byte, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	var m = messaging.Message{
		DataType: 0x01,
		AesIV:    security.RandomKey(),
	}
	m.Data = data
	return m.ToBytes(), nil
}

func ObjectToBytesAuth(a interface{}) ([]byte, error) {
	return objectToBytes(a, 0x01)
}
