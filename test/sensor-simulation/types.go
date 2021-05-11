package sensor_simulation

type CalibrationStep struct {
	Description string
	Value       float64
	Completed   bool
}

type Device struct {
	Name             string
	Protocol         string
	CalibrationSteps []CalibrationStep
}
