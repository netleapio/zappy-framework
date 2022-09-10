//go:build baremetal

package framework

import (
	"runtime"
	"time"
)

func setTimeFromRTC(rtc RTC) error {
	DebugOut("adjusting system time")
	rtcNow, err := rtc.Now()
	if err != nil {
		return err
	}

	// We could loop until the RTC 'ticks over' to be closely aligned to RTC time, but that
	// could delay for a significant period (1s?).  Instead, get within 50% of RTC accuracy as
	// 'good enough'
	runtime.AdjustTimeOffset(int64(rtcNow.Sub(time.Now())) + rtc.Accuracy().Nanoseconds()/2)

	DebugOut("time adjusted")
	return nil
}
