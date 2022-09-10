package protocol

import (
	"errors"
)

var (
	ErrTooSmall = errors.New("too small")
	ErrInvalid  = errors.New("invalid")
)

type SensorReport Packet

type ReadingType uint8

const (
	ReadingTypeBattVolts ReadingType = iota
	ReadingTypeTemperature
	ReadingTypePressure
	ReadingTypeHumidity
)

func (r *SensorReport) AddBatteryVoltage(milliVolts uint16) {
	(*Packet)(r).WriteUint16(uint16(ReadingTypeBattVolts))
	(*Packet)(r).WriteUint16(milliVolts)
}

func (r *SensorReport) AddTemperature(centiCelcius uint16) {
	(*Packet)(r).WriteUint16(uint16(ReadingTypeTemperature))
	(*Packet)(r).WriteUint16(centiCelcius)
}

func (r *SensorReport) AddHumidity(centiPercent uint16) {
	(*Packet)(r).WriteUint16(uint16(ReadingTypeHumidity))
	(*Packet)(r).WriteUint16(centiPercent)
}

func (r *SensorReport) AddPressure(dekaPascal uint16) {
	(*Packet)(r).WriteUint16(uint16(ReadingTypePressure))
	(*Packet)(r).WriteUint16(dekaPascal)
}
