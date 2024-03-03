package protocol

import (
	"encoding/binary"
	"strings"
)

//
//   Packet Format
//
//   31                              15                             0
//   +-------------------------------+------------------------------+
//   | NetworkID                     | DeviceID                     |
//   +-------------------------------+------------------------------+
//   | Alerts                        | Version                      |
//   +-------------------------------+------------------------------+
//   | CRC (Version 2+)              | Type                         |
//   +-------------------------------+------------------------------+
//

// Packet is a fixed-sized buffer (max LoRa payload)
type Packet struct {
	data [256]byte
	ptr  uint8
	len  uint8
}

const (
	HeaderLenV1 = 10
	HeaderLenV2 = 12

	crcKey = 0xA152

	deviceIDHeaderOffset  = 0
	networkIDHeaderOffset = 2
	versionHeaderOffset   = 4
	alertsHeaderOffset    = 6
	typeHeaderOffset      = 8
	crcHeaderOffset       = 10
)

const CurrentVersion = 2

type Alerts uint16

const (
	AlertNone    Alerts = 0
	AlertBattLow Alerts = 1 << iota
	AlertBattCritical
	AlertRTCFailure
)

func (a Alerts) Strings() []string {
	flagStrs := make([]string, 0, 16)

	if a&AlertBattLow != 0 {
		flagStrs = append(flagStrs, "BattLow")
	}
	if a&AlertBattCritical != 0 {
		flagStrs = append(flagStrs, "BattCritical")
	}
	if a&AlertRTCFailure != 0 {
		flagStrs = append(flagStrs, "RTCFailure")
	}

	return flagStrs
}

func (a Alerts) String() string {

	return "[" + strings.Join(a.Strings(), ",") + "]"
}

type PacketType uint16

const (
	TypeAnnounce        PacketType = 0x0000
	TypeSensorReport    PacketType = 0x0001
	TypeConfigureDevice PacketType = 0x8000
)

func (p *Packet) Reset() {
	p.len = 0
	p.ptr = 0
	p.SetVersion(CurrentVersion)
}

func (p *Packet) DeviceID() uint16 {
	return binary.BigEndian.Uint16(p.data[deviceIDHeaderOffset:])
}

func (p *Packet) SetDeviceID(device uint16) {
	binary.BigEndian.PutUint16(p.data[deviceIDHeaderOffset:], device)
}

func (p *Packet) NetworkID() uint16 {
	return binary.BigEndian.Uint16(p.data[networkIDHeaderOffset:])
}

func (p *Packet) SetNetworkID(network uint16) {
	binary.BigEndian.PutUint16(p.data[networkIDHeaderOffset:], network)
}

func (p *Packet) Version() uint16 {
	return binary.BigEndian.Uint16(p.data[versionHeaderOffset:])
}

func (p *Packet) SetVersion(version uint16) {
	binary.BigEndian.PutUint16(p.data[versionHeaderOffset:], version)
}

func (p *Packet) Alerts() Alerts {
	return Alerts(binary.BigEndian.Uint16(p.data[alertsHeaderOffset:]))
}

func (p *Packet) SetAlerts(alerts Alerts) {
	binary.BigEndian.PutUint16(p.data[alertsHeaderOffset:], uint16(alerts))
}

func (p *Packet) Type() PacketType {
	return PacketType(binary.BigEndian.Uint16(p.data[typeHeaderOffset:]))
}

func (p *Packet) SetType(t PacketType) {
	binary.BigEndian.PutUint16(p.data[typeHeaderOffset:], uint16(t))
}

func (p *Packet) UpdateCRC() {
	// Set CRC field to 'unique' value 0xA152 just in case
	// any other protocols happen to use same CRC in same place
	// (highly unlikely, but...)
	binary.BigEndian.PutUint16(p.data[crcHeaderOffset:], crcKey)

	crc := calcCrc(p.data[:p.len])

	binary.BigEndian.PutUint16(p.data[crcHeaderOffset:], crc)
}

func (p *Packet) CRCValid() bool {
	// V1 has no CRC, so indicate always valid
	if p.Version() == 1 {
		return true
	}

	storedCRC := binary.BigEndian.Uint16(p.data[crcHeaderOffset:])

	// Temporarily over-write the packet CRC with the 'key' and
	// ensure it's restored before the function exits.  This is ugly,
	// but...
	defer func() {
		binary.BigEndian.PutUint16(p.data[crcHeaderOffset:], storedCRC)
	}()
	binary.BigEndian.PutUint16(p.data[crcHeaderOffset:], crcKey)

	calcCrc := calcCrc(p.data[:p.len])

	return calcCrc == storedCRC
}

// Remaining gets the number of unread bytes in the packet
func (p *Packet) Remaining() uint8 {
	return p.len - p.ptr
}

func (p *Packet) WriteUint16(v uint16) {
	if p.ptr < p.HeaderLen() {
		p.ptr = p.HeaderLen()
	}

	binary.BigEndian.PutUint16(p.data[p.ptr:], v)
	p.ptr += 2
	p.len = p.ptr
}

func (p *Packet) ReadUint16() uint16 {
	if p.ptr < p.HeaderLen() {
		p.ptr = p.HeaderLen()
	}

	v := binary.BigEndian.Uint16(p.data[p.ptr:])
	p.ptr += 2
	return v
}

func (p *Packet) Skip(n uint8) {
	if p.ptr < p.HeaderLen() {
		p.ptr = p.HeaderLen()
	}

	p.ptr += n
}

func (p *Packet) AsBytes() []byte {
	return p.data[:p.len]
}

// SetLength sets (and resets) the size of the packet
func (p *Packet) SetLength(len uint8) {
	p.len = len
	if p.ptr > p.len {
		p.ptr = p.len
	}
}

// HeaderLen gets the size of the header for this packet
func (p *Packet) HeaderLen() uint8 {
	if p.Version() < 2 {
		return HeaderLenV1
	}

	return HeaderLenV2
}

// DetectMessage scans a packet to determine the message encoded
func DetectMessage(pkt *Packet) Message {

	var msg Message

	switch pkt.Type() {
	case TypeSensorReport:
		msg = &SensorReport{}
	}

	if msg != nil {
		msg.AttachPacket(pkt)
	}

	return msg
}

// calcCrc implements MODBUS-16 CRC
func calcCrc(buf []byte) uint16 {
	crc := uint16(0xFFFF)

	for _, b := range buf {
		crc ^= uint16(b)

		for bit := 0; bit < 8; bit++ {
			if (crc & 0x0001) != 0 {
				crc >>= 1
				crc ^= 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
