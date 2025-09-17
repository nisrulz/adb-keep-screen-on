package internal

import (
	"adb-keep-screen-on/internal/adb"
	"os"
	"os/signal"
	"syscall"
)

// OriginalStayAwake stores the original stay_on_while_plugged_in value
var OriginalStayAwake string

// SetOriginalStayAwake sets the original stay_on_while_plugged_in value
func SetOriginalStayAwake(val string) {
	OriginalStayAwake = val
}

// OriginalStayAwakeMap stores the original stay_on_while_plugged_in value for each device
var OriginalStayAwakeMap = make(map[string]string)

// SetOriginalStayAwakeMap sets the original stay awake map from main.go
func SetOriginalStayAwakeMap(m map[string]string) {
	for k, v := range m {
		OriginalStayAwakeMap[k] = v
	}
}

// HandleInterrupt sets up a signal handler for Ctrl+C (SIGINT) to restore device settings and exit
func HandleInterrupt() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		// Restore original stay_on_while_plugged_in value if available
		for deviceID, val := range OriginalStayAwakeMap {
			_ = adb.SetStayAwakeValue(deviceID, val)
		}
		println("\n🔌 Monitoring stopped. Device settings restored. Exiting...\n")
		os.Exit(0)
	}()
}
