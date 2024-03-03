//go:build tinygo

package framework

type App interface {
	Initialize(fwk *Framework) error
	Triggered(fwk *Framework) error
	USBPowered(fwk *Framework) error
}
