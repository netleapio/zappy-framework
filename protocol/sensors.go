package protocol

type SensorType uint8

const (
	SensorTypeBattVolts SensorType = iota
	SensorTypeTemperature
	SensorTypePressure
	SensorTypeHumidity
)

var AllSensorTypes = []SensorType{
	SensorTypeBattVolts,
	SensorTypeTemperature,
	SensorTypePressure,
	SensorTypeHumidity,
}

// SensorInfo provides meta-data about each type of sensor, including
// it's common name, SI unit and conversion ratios of the protocol value
// to SI unit.
//
// The conversion ratios are modeled as two integers 'Mult' and 'Div'
// enabling an integer representation (if needed).
type SensorInfo struct {
	Name string
	Unit string
	Mult int
	Div  int
}

var SensorMetadata = map[SensorType]*SensorInfo{
	SensorTypeBattVolts:   {Name: "battery", Unit: "volts", Mult: 1, Div: 1000},
	SensorTypeTemperature: {Name: "temperature", Unit: "celcius", Mult: 1, Div: 100},
	SensorTypePressure:    {Name: "pressure", Unit: "pascals", Mult: 10, Div: 1},
	SensorTypeHumidity:    {Name: "humidity", Unit: "ratio", Mult: 1, Div: 10000},
}
