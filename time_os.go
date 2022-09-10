//go:build !baremetal

package framework

func setTimeFromRTC(rtc RTC) error {
	DebugOut("ignoring time change from RTC")
	return nil
}
