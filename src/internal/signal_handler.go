package internal

import (
	"adb-keep-screen-on/src/internal/adb"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// OriginalSettings stores original stay_on_while_plugged_in values per device.
// main assigns its map reference here so the signal handler sees live updates.
var OriginalSettings map[string]string

// PidPath is the path to the PID file. Set by main before HandleInterrupt is called.
var PidPath string

// RestoreOnExit controls whether original device settings are restored on exit.
var RestoreOnExit = true

// HandleInterrupt sets up a signal handler for SIGINT/SIGTERM to restore device settings and exit
func HandleInterrupt() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		if PidPath != "" {
			os.Remove(PidPath)
		}
		if RestoreOnExit {
			for deviceID, val := range OriginalSettings {
				_ = adb.SetStayAwakeValue(deviceID, val)
			}
			fmt.Println("\n🔌 Monitoring stopped. Device settings restored. Exiting...")
		} else {
			fmt.Println("\n🔌 Monitoring stopped. Exiting...")
		}
		os.Exit(0)
	}()
}
