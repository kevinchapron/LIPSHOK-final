package database

import "gorm.io/gorm"

type DatabaseDevice struct {
	gorm.Model
	PhysicalAddr string
	Name         string

	ConnectionType string
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
		&DatabaseDeviceCalibrationStep{},
		&DatabaseDevice{},
	}
}
