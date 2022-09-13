// Framework provides the glue logic between board abstractions and apps
package framework

import (
	"tinygo.org/x/drivers"
)

// Board provides the common abstraction for all compatible boards
type Board interface {
	// Initialize the board, configuration of peripherals, etc sufficient
	// to initialize the RTC.
	InitializePreRTC() error

	// Initialize the board, configuration of peripherals, etc.  This is
	// after RTC, so time is established.
	InitializePostRTC() error

	// DeepSleep instructs the board abstraction to:
	//  1. reduce power consumption as far as possible (excepting RTC alarm)
	//  2. put MCU into deepest sleep state permissible
	//  3. wait on RTC alarm / USB plug-in
	//  4. restore power to peripherals
	//
	// If RTC is currently alerting, or USB is currently plugged-in, it should
	// return immediately.
	//
	// It is acceptable for the MCU to loose all state during deep-sleep if
	// the RTC alarm and USB-wake events are handled.  In which case, this
	// method may never return.
	//
	DeepSleep() WakeReason

	// Sensor I2C
	SensorI2C() (drivers.I2C, error)

	// Display access
	Display() (Display, error)

	// RTC access
	RTC() (RTC, error)

	// LoRa radio access
	Radio() (Radio, error)

	// BatteryVoltage gets access to the current battery voltage level (in milliVolts)
	BatteryVoltage() uint16
}

// WakeReason indicates the (one or more) reasons why a board awoke
// from DeepSleep
type WakeReason uint8

const (
	WakeUnknown    WakeReason = 0
	WakeByRTCAlarm WakeReason = 1 << iota
	WakeByUSB
)
