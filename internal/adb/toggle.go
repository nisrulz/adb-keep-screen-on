package adb

import (
	"fmt"
	"os/exec"
)

// SetStayAwakeValue sets stay_on_while_plugged_in to a specific value on the device
func SetStayAwakeValue(deviceID, val string) error {
	cmd := exec.Command("adb", "-s", deviceID, "shell", "settings", "put", "global", "stay_on_while_plugged_in", val)
	return cmd.Run()
}

// WakeScreen sends an ADB command to wake up the device screen
func WakeScreen(deviceID string) {
	cmd := exec.Command("adb", "-s", deviceID, "shell", "input", "keyevent", "KEYCODE_WAKEUP")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("⚠️ Failed to wake device %s screen: %v\n", deviceID, err)
	}
}

// SetStayAwake sets the device's stay awake setting while plugged in
func SetStayAwake(deviceID string, enable bool) {
	value := "0"
	if enable {
		value = "1"
	}
	// Attempt to set the stay awake value
	err := SetStayAwakeValue(deviceID, value)
	if err != nil {
		fmt.Printf("⚠️ Failed to set stay_on_while_plugged_in to %s for device %s: %v\n", value, deviceID, err)
	}
	// Wake the screen if enabling stay awake
	if enable {
		WakeScreen(deviceID)
	}
}

