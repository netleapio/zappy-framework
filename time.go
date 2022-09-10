package framework

func initTime(rtc RTC) error {
	err := setTimeFromRTC(rtc)
	if err != nil {
		return err
	}

	// Configure periodic alarm to wake
	alarmEn, err := rtc.IsAlarmEnabled()
	if err != nil {
		return err
	}

	// Default to once per min at 2 secs past
	if !alarmEn {
		DebugOut("setup RTC alarm")
		err = rtc.SetPeriodicAlarm(2, RTCSecs)
		if err != nil {
			return err
		}
	}

	return nil
}
