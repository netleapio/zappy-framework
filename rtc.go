package framework

import "time"

// RTC provides access to the RTC on the board.  This may be
// part of the MCU or a separate IC.
//
// All boards are expected to have an RTC which at least provides
// the ability to get and set the time as well as at least one
// alarm.
type RTC interface {
	// Now gets the current time according to the RTC
	//
	// Note: RTC accuracy could be as low as 1s
	Now() (time.Time, error)

	// Accuracy gives the accuracy of the RTC (e.g. 1s)
	Accuracy() time.Duration

	// SetTime adjusts the RTC time
	SetTime(time.Time) error

	// SetPeriodicAlarm initializes an alarm that will trigger repeatedly.
	SetPeriodicAlarm(value uint8, field RTCField) error

	// IsAlarmEnabled indicates if an alarm is currently configured.
	IsAlarmEnabled() (bool, error)

	// IsAlarmTriggered indicates if an alarm is currently triggered.
	IsAlarmTriggered() (bool, error)

	// AcknowledgeAlarm clears the interrupt status associated with the alarm.
	AcknowledgeAlarm() error

	// Healthy indicates if the RTC is in a healthy running state.
	IsHealthy() bool

	// Dump displays diagnostic info to console
	Dump()
}

type RTCField uint8

// RTC fields for alarms
const (
	RTCSecs RTCField = iota
	RTCMins
	RTCHours
	RTCWeekDay
	RTCDay
)
