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

	interval := flag.Int("interval", 10, "Polling interval in seconds")
	flag.Parse()

	println("🔌 Monitoring ADB connection every", *interval, "seconds to keep screen on...\n")

	// Setup signal handler for graceful exit
	internal.HandleInterrupt()

	var lastConnectedDevices = make(map[string]bool)
	var originalSettings = make(map[string]string)

	for {
		devices := adb.ListDevices()
		if len(devices) == 0 {
			println("⏳ Waiting for device connection...\n")
		} else {
			for _, deviceID := range devices {
				if !adb.IsDeviceConnected(deviceID) {
					continue
				}
				connected := true
				if lastConnectedDevices[deviceID] != connected {
					// Store original setting for this device
					if _, ok := originalSettings[deviceID]; !ok {
						setting, err := adb.GetStayAwakeSetting(deviceID)
						if err == nil {
							originalSettings[deviceID] = setting
						}
					}
					println("📱 Device connected via ADB:", deviceID, ". Keeping screen awake.\n")
					println("✅ Device is connected. Screen will stay awake.\n")
					adb.SetStayAwake(deviceID, true)
					lastConnectedDevices[deviceID] = connected
					internal.SetOriginalStayAwakeMap(originalSettings)
				}
			}
			// Check for disconnected devices
			for deviceID, wasConnected := range lastConnectedDevices {
				stillConnected := false
				for _, d := range devices {
					if d == deviceID {
						stillConnected = true
						break
					}
				}
				if wasConnected && !stillConnected {
					println("❌ Device disconnected:", deviceID, "\n")
					// Restore original setting for this device
					if val, ok := originalSettings[deviceID]; ok {
						adb.SetStayAwake(deviceID, false)
						adb.SetStayAwakeValue(deviceID, val)
						delete(originalSettings, deviceID)
						internal.SetOriginalStayAwakeMap(originalSettings)
					}
					delete(lastConnectedDevices, deviceID)
				}
			}
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
