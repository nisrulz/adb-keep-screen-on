package main

import (
	"flag"
	"time"

	"adb-keep-screen-on/internal"
	"adb-keep-screen-on/internal/adb"
)

// main is the entry point for the CLI tool
func main() {
	// Check if adb is available
	if !adb.IsAdbAvailable() {
		return
	}

	// Store the original stay_on_while_plugged_in value
	original, err := adb.GetStayAwakeSetting()
	if err == nil {
		internal.SetOriginalStayAwake(original)
	}

	interval := flag.Int("interval", 10, "Polling interval in seconds")
	flag.Parse()

	println("🔌 Monitoring ADB connection every", *interval, "seconds to keep screen on...\n")

	// Setup signal handler for graceful exit
	internal.HandleInterrupt()

	var lastConnected bool = false

	for {
		// Check device connection status
		connected := adb.IsDeviceConnected()

		// Only print and act if connection state changes
		if connected != lastConnected {
			if connected {
				println("📱 Device connected via ADB. Keeping screen awake.\n")
				println("✅ Device is connected. Screen will stay awake.\n")
				adb.SetStayAwake(true)
			} else {
				println("❌ Device disconnected\n")
				println("⏳ Waiting for device connection...\n")
				adb.SetStayAwake(false)
			}
			lastConnected = connected
		}

		// Wait for the next polling interval
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
