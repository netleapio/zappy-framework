package protocol

import (
	"errors"
)

var (
	ErrTooSmall = errors.New("too small")
	ErrInvalid  = errors.New("invalid")
)

// SensorReport is a wrapper for a packet.
//
// Sensor reports are sent from devices to the controller indicating
// the current status of all sensors available on the device.
//
// Devices may have different sensing capabilities, so the different
// types of sensor are optional to include.
type SensorReport struct {
	packet *Packet
}

func (r *SensorReport) AttachPacket(packet *Packet) {
	r.packet = packet
}

func (r *SensorReport) Packet() *Packet {
	return r.packet
}

func (r *SensorReport) AddBatteryVoltage(milliVolts uint16) {
	r.packet.WriteUint16(uint16(SensorTypeBattVolts))
	r.packet.WriteUint16(milliVolts)
}

func (r *SensorReport) HasBatteryVoltage() bool {
	return r.HasReadingType(SensorTypeBattVolts)
}

func (r *SensorReport) BatteryVoltage() uint16 {
	return r.GetReadingUint16(SensorTypeBattVolts, 0)
}

func (r *SensorReport) AddTemperature(centiCelcius uint16) {
	r.packet.WriteUint16(uint16(SensorTypeTemperature))
	r.packet.WriteUint16(centiCelcius)
}

func (r *SensorReport) HasTemperature() bool {
	return r.HasReadingType(SensorTypeTemperature)
}

func (r *SensorReport) Temperature() uint16 {
	return r.GetReadingUint16(SensorTypeTemperature, 0)
}

func (r *SensorReport) AddHumidity(centiPercent uint16) {
	r.packet.WriteUint16(uint16(SensorTypeHumidity))
	r.packet.WriteUint16(centiPercent)
}

func (r *SensorReport) HasHumidity() bool {
	return r.HasReadingType(SensorTypeHumidity)
}

func (r *SensorReport) Humidity() uint16 {
	return r.GetReadingUint16(SensorTypeHumidity, 0)
}

func (r *SensorReport) AddPressure(dekaPascal uint16) {
	r.packet.WriteUint16(uint16(SensorTypePressure))
	r.packet.WriteUint16(dekaPascal)
}

func (r *SensorReport) HasPressure() bool {
	return r.HasReadingType(SensorTypePressure)
}

func (r *SensorReport) Pressure() uint16 {
	return r.GetReadingUint16(SensorTypePressure, 0)
}

func (r *SensorReport) AddSupplyVoltage(milliVolts uint16) {
	r.packet.WriteUint16(uint16(SensorTypeSupplyVolts))
	r.packet.WriteUint16(milliVolts)
}

func (r *SensorReport) HasSupplyVoltage() bool {
	return r.HasReadingType(SensorTypeSupplyVolts)
}

func (r *SensorReport) SupplyVoltage() uint16 {
	return r.GetReadingUint16(SensorTypeSupplyVolts, 0)
}

func (r *SensorReport) AddLoadPower(deciVolts uint16) {
	r.packet.WriteUint16(uint16(SensorTypeLoadPower))
	r.packet.WriteUint16(deciVolts)
}

func (r *SensorReport) HasLoadPower() bool {
	return r.HasReadingType(SensorTypeLoadPower)
}

func (r *SensorReport) LoadPower() uint16 {
	return r.GetReadingUint16(SensorTypeLoadPower, 0)
}

func (r *SensorReport) AddCoils(coils uint16) {
	r.packet.WriteUint16(uint16(SensorTypeCoils))
	r.packet.WriteUint16(coils)
}

func (r *SensorReport) HasCoils() bool {
	return r.HasReadingType(SensorTypeCoils)
}

func (r *SensorReport) Coils() uint16 {
	return r.GetReadingUint16(SensorTypeCoils, 0)
}

func (r *SensorReport) AllReadings() map[SensorType]uint16 {
	result := make(map[SensorType]uint16)

	r.packet.ptr = r.packet.HeaderLen()
	for r.packet.Remaining() >= 4 {
		fieldType := SensorType(r.packet.ReadUint16())
		fieldValue := r.packet.ReadUint16()

		result[fieldType] = fieldValue
	}

	return result
}

func (r *SensorReport) HasReadingType(t SensorType) bool {
	r.packet.ptr = r.packet.HeaderLen()
	for r.packet.Remaining() >= 4 {
		fieldType := SensorType(r.packet.ReadUint16())
		if fieldType == t {
			return true
		}

		// skip value
		r.packet.Skip(2)
	}

	return false
}

// Gets the value of a reading or a default value
func (r *SensorReport) GetReadingUint16(t SensorType, def uint16) uint16 {
	r.packet.ptr = r.packet.HeaderLen()
	for r.packet.Remaining() >= 4 {
		fieldType := SensorType(r.packet.ReadUint16())
		fieldValue := r.packet.ReadUint16()
		if fieldType == t {
			return fieldValue
		}
	}

	return def
}
