package protocol

type SensorType uint8

const (
	SensorTypeBattVolts SensorType = iota
	SensorTypeTemperature
	SensorTypePressure
	SensorTypeHumidity
	SensorTypeSupplyVolts
	SensorTypeLoadPower
	SensorTypeCoils
)

var AllSensorTypes = []SensorType{
	SensorTypeBattVolts,
	SensorTypeTemperature,
	SensorTypePressure,
	SensorTypeHumidity,
	SensorTypeSupplyVolts,
	SensorTypeLoadPower,
	SensorTypeCoils,
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
	SensorTypeTemperature: {Name: "temperature", Unit: "celsius", Mult: 1, Div: 100},
	SensorTypePressure:    {Name: "pressure", Unit: "pascals", Mult: 10, Div: 1},
	SensorTypeHumidity:    {Name: "humidity", Unit: "percent", Mult: 1, Div: 100},
	SensorTypeSupplyVolts: {Name: "supply", Unit: "volts", Mult: 1, Div: 1000},
	SensorTypeLoadPower:   {Name: "load", Unit: "watts", Mult: 1, Div: 10},
	SensorTypeCoils:       {Name: "coils", Unit: "", Mult: 1, Div: 1},
}
