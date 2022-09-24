package protocol

import (
	"encoding/binary"
	"strings"
)

//
//   Header Format
//
//   31                              15                             0
//   +-------------------------------+------------------------------+
//   | DeviceID                      | NetworkID                    |
//   +-------------------------------+------------------------------+
//   | Alerts                        | Version                      |
//   +-------------------------------+------------------------------+
//   | Type                          |
//   +-------------------------------+

// Packet is a fixed-sized buffer (max LoRa payload)
type Packet struct {
	data [256]byte
	ptr  uint8
	len  uint8
}

const HeaderLen = 10

const CurrentVersion = 1

type Alerts uint16

const (
	AlertNone    Alerts = 0
	AlertBattLow Alerts = 1 << iota
	AlertBattCritical
	AlertRTCFailure
)

func (a Alerts) String() string {
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

	return "[" + strings.Join(flagStrs, ",") + "]"
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
}

func (p *Packet) DeviceID() uint16 {
	return binary.BigEndian.Uint16(p.data[0:])
}

func (p *Packet) SetDeviceID(device uint16) {
	binary.BigEndian.PutUint16(p.data[0:], device)
}

func (p *Packet) NetworkID() uint16 {
	return binary.BigEndian.Uint16(p.data[2:])
}

func (p *Packet) SetNetworkID(network uint16) {
	binary.BigEndian.PutUint16(p.data[2:], network)
}

func (p *Packet) Version() uint16 {
	return binary.BigEndian.Uint16(p.data[4:])
}

func (p *Packet) SetVersion(version uint16) {
	binary.BigEndian.PutUint16(p.data[4:], version)
}

func (p *Packet) Alerts() Alerts {
	return Alerts(binary.BigEndian.Uint16(p.data[6:]))
}

func (p *Packet) SetAlerts(alerts Alerts) {
	binary.BigEndian.PutUint16(p.data[6:], uint16(alerts))
}

func (p *Packet) Type() PacketType {
	return PacketType(binary.BigEndian.Uint16(p.data[8:]))
}

func (p *Packet) SetType(t PacketType) {
	binary.BigEndian.PutUint16(p.data[8:], uint16(t))
}

// Remaining gets the number of unread bytes in the packet
func (p *Packet) Remaining() uint8 {
	return p.len - p.ptr
}

func (p *Packet) WriteUint16(v uint16) {
	if p.ptr < HeaderLen {
		p.ptr = HeaderLen
	}

	binary.BigEndian.PutUint16(p.data[p.ptr:], v)
	p.ptr += 2
	p.len = p.ptr
}

func (p *Packet) ReadUint16() uint16 {
	if p.ptr < HeaderLen {
		p.ptr = HeaderLen
	}

	v := binary.BigEndian.Uint16(p.data[p.ptr:])
	p.ptr += 2
	return v
}

func (p *Packet) Skip(n uint8) {
	if p.ptr < HeaderLen {
		p.ptr = HeaderLen
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
