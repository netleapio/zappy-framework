package protocol

import "encoding/binary"

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
	ptr  int
}

const HeaderLen = 10

const CurrentVersion = 1

type Alerts uint16

const (
	AlertNone    Alerts = 0
	AlertBattLow Alerts = 1 << iota
	AlertBattCritical
)

type PacketType uint16

const (
	TypeAnnounce        PacketType = 0x0000
	TypeSensorReport    PacketType = 0x0001
	TypeConfigureDevice PacketType = 0x8000
)

func (p *Packet) SetDeviceID(device uint16) {
	binary.BigEndian.PutUint16(p.data[0:], device)
}

func (p *Packet) SetNetworkID(network uint16) {
	binary.BigEndian.PutUint16(p.data[2:], network)
}

func (p *Packet) SetVersion(version uint16) {
	binary.BigEndian.PutUint16(p.data[4:], version)
}

func (p *Packet) SetAlerts(alerts Alerts) {
	binary.BigEndian.PutUint16(p.data[6:], uint16(alerts))
}

func (p *Packet) SetType(t PacketType) {
	binary.BigEndian.PutUint16(p.data[8:], uint16(t))
}

func (p *Packet) WriteUint16(v uint16) {
	if p.ptr < HeaderLen {
		p.ptr = HeaderLen
	}

	binary.BigEndian.PutUint16(p.data[p.ptr:], v)
	p.ptr += 2
}

func (p *Packet) AsBytes() []byte {
	return p.data[:p.ptr]
}
