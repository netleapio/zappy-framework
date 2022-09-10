//go:build !nodebug

package framework

import "log"

//
// This implementation of debug output sends output to TinyGo serial.
//
// Configure the TinyGo serial output device using the serial option, eg:
// `tinygo ... -serial uart`
//

const DebugEnabled = true

func DebugOut(msg ...string) {
	log.Println(msg)
}
