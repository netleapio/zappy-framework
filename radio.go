package framework

// Radio is the interface to a LoRa peripheral
type Radio interface {
	// Tx transmits a packet
	Tx(pkt []uint8, timeoutMs uint32) error

	// Rx receives a packet, blocking until received
	Rx(timeoutMs uint32, buf []uint8) (int, error)
}
