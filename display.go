package framework

import "image/color"

type Display interface {
	// Size returns the current size of the display.
	Size() (x, y int16)

	// SetPizel modifies the internal buffer.
	SetPixel(x, y int16, c color.RGBA)

	// ClearBuffer resets all pixels
	ClearBuffer()

	// Display sends the buffer (if any) to the screen.
	Display() error
}

// NullDisplay can be used on devices with no attached display
type NullDisplay struct {
	width  int16
	height int16
}

func NewNullDisplay(width, height int16) *NullDisplay {
	return &NullDisplay{width: width, height: height}
}

func (d *NullDisplay) Size() (x, y int16) {
	return d.width, d.height
}

func (d *NullDisplay) SetPixel(x, y int16, c color.RGBA) {
}

func (d *NullDisplay) ClearBuffer() {
}

func (d *NullDisplay) Display() error {
	return nil
}

// WaitUntilIdle helps emulate e-Paper type displays
func (d *NullDisplay) WaitUntilIdle() {
}

// DeepSleep helps emulate e-Paper type displays
func (d *NullDisplay) DeepSleep() {
}
