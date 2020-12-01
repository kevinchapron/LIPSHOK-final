package database

import "gorm.io/gorm"

const DATABASE_DEVICE_CONNECTION_WIFI = "Wi-Fi"
const DATABASE_DEVICE_CONNECTION_BLUETOOTH = "Bluetooth"
const DATABASE_DEVICE_CONNECTION_BLE = "BLE"

type DatabaseDeviceType struct {
	gorm.Model
	Name      string
	Activated bool
}

type DatabaseDevice struct {
	gorm.Model
	PhysicalAddr string
	Name         string

	ConnectionType   DatabaseDeviceType
	ConnectionTypeID uint
}

type DatabaseDeviceCalibrationStep struct {
	gorm.Model
	Description string
	Value       float64
	Completed   bool

	Device   DatabaseDevice
	DeviceID uint
}

func migrations() []interface{} {
	return []interface{}{
		&DatabaseDeviceType{},
		&DatabaseDeviceCalibrationStep{},
		&DatabaseDevice{},
	}
}
