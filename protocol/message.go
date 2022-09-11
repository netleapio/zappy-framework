package protocol

// Message is the abstraction for types that wrap packets providing
// type-safe access to packet contents.
type Message interface {
	AttachPacket(packet *Packet)
	Packet() *Packet
}
