package sensor_simulation

import (
	"encoding/json"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"github.com/kevinchapron/LIPSHOK/security"
)

func GetDevice() *Device {
	return &Device{
		Name:     "TestingSensor",
		Protocol: "UDP",
		CalibrationSteps: []CalibrationStep{
			{Description: "Itération 1", Value: -1, Completed: false},
			{Description: "Itération 2", Value: -1, Completed: false},
			{Description: "Itération 3", Value: -1, Completed: false},
		},
	}
}

func objectToBytes(a interface{}, dataType byte) ([]byte, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	var m = messaging.Message{
		DataType: dataType,
		AesIV:    security.RandomKey(),
	}
	m.Data = data
	return m.ToBytes(), nil
}

func ObjectToBytesAuth(a interface{}) ([]byte, error) {
	return objectToBytes(a, constants.MESSAGING_DATATYPE_AUTH)
}

func ObjectToBytesData(a interface{}) ([]byte, error) {
	return objectToBytes(a, constants.MESSAGING_DATATYPE_DATA)
}
