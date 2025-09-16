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

// HandleInterrupt sets up a signal handler for Ctrl+C (SIGINT) to restore device settings and exit
func HandleInterrupt() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		// Restore original stay_on_while_plugged_in value if available
		if OriginalStayAwake != "" {
			_ = adb.SetStayAwakeValue(OriginalStayAwake)
		}
		println("\n🔌 Monitoring stopped. Device settings restored. Exiting...\n")
		os.Exit(0)
	}()
}
