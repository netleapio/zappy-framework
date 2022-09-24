package framework

import (
	"github.com/netleapio/zappy-framework/protocol"
)

type Framework struct {
	Board Board
	App   App
}

func New(board Board, app App) *Framework {
	return &Framework{
		Board: board,
		App:   app,
	}
}

// Run is the function called by the app to hand over control to the framework.
func (f *Framework) Run() error {
	// On any failure, try to enter deep sleep
	// TODO: watchdog
	defer f.Board.DeepSleep()

	err := f.Board.InitializePreRTC()
	if err != nil {
		// Failure to initialize board is fatal
		DebugOut("failed to initialize board (pre RTC)")
		panic(err)
	}

	rtc, err := f.Board.RTC()
	if err != nil || rtc == nil {
		// Failover to get RTC is fatal
		DebugOut("RTC not available")
		panic(err)
	}

	// Set time according to RTC
	err = initTime(rtc)
	if err != nil {
		// Failing to establish time is fatal
		DebugOut("failed to initialize time from RTC")
		panic(err)
	}

	err = f.Board.InitializePostRTC()
	if err != nil {
		// Failure to initialize board is fatal
		DebugOut("failed to initialize board (post RTC)")
		panic(err)
	}

	err = f.App.Initialize(f)
	if err != nil {
		DebugOut("app failed to initialize:", err.Error())
		panic(err)
	}

	for {
		triggered, err := rtc.IsAlarmTriggered()
		if err != nil {
			DebugOut("failed to get RTC alarm status")
			panic(err)
		}

		if triggered {
			rtc.AcknowledgeAlarm()

			err = f.App.Triggered(f)
			if err != nil {
				DebugOut("app error on trigger:", err.Error())
			}
		}

		reason := f.Board.DeepSleep()

		if reason&WakeByUSB != 0 {
			f.App.USBPowered(f)
		}
	}
}

// Send transmits a message using the zappy protocol
func (f *Framework) Send(msg protocol.Message, t protocol.PacketType) error {
	// TODO: async/queuing

	r, err := f.Board.Radio()
	if err != nil {
		return err
	}

	pkt := msg.Packet()

	pkt.SetVersion(protocol.CurrentVersion)
	pkt.SetType(t)
	pkt.SetNetworkID(0)
	pkt.SetDeviceID(0)
	pkt.SetAlerts(f.currentAlerts())

	data := pkt.AsBytes()

	return r.Tx(data, 1000)
}

func (f *Framework) currentAlerts() protocol.Alerts {
	val := protocol.AlertNone

	battV := f.Board.BatteryVoltage()

	if battV < LowBatteryAlertVoltage {
		val |= protocol.AlertBattLow
	}
	if battV < CriticalBatteryAlertVoltage {
		val |= protocol.AlertBattCritical
	}

	rtc, err := f.Board.RTC()
	if err != nil || rtc == nil || !rtc.IsHealthy() {
		val |= protocol.AlertRTCFailure
	}

	return val
}
