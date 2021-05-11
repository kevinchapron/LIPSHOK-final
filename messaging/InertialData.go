package messaging

import "math"

type InertialData struct {
	AX float64
	AY float64
	AZ float64

	GX float64
	GY float64
	GZ float64
}

func (data *InertialData) HasWeirdData() bool {
	return math.Abs(data.AX) > 4 || math.Abs(data.AY) > 4 || math.Abs(data.AZ) > 4
}

func (data *InertialData) RoundUpData() {
	data.AX = math.Round(data.AX*1000) / 1000
	data.AY = math.Round(data.AY*1000) / 1000
	data.AZ = math.Round(data.AZ*1000) / 1000

	data.GX = math.Round(data.GX*1000) / 1000
	data.GY = math.Round(data.GY*1000) / 1000
	data.GZ = math.Round(data.GZ*1000) / 1000
}
