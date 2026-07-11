package adb

import "fmt"

// SetStayAwakeValue sets stay_on_while_plugged_in to a specific value on the device
func SetStayAwakeValue(deviceID, val string) error {
	cmd := execCommand("adb", "-s", deviceID, "shell", "settings", "put", "global", "stay_on_while_plugged_in", val)
	return cmd.Run()
}

// WakeScreen sends a wake keyevent to the device.
func WakeScreen(deviceID string) error {
	cmd := execCommand("adb", "-s", deviceID, "shell", "input", "keyevent", "KEYCODE_WAKEUP")
	return cmd.Run()
}

// SetStayAwake sets the device's stay awake setting while plugged in. When
// enabling, the screen is also woken if wake is true.
func SetStayAwake(deviceID string, enable, wake bool) {
	value := "0"
	if enable {
		value = "1"
	}
	if err := SetStayAwakeValue(deviceID, value); err != nil {
		fmt.Printf("⚠️  Could not update setting for device %s (is it still connected?)\n", deviceID)
	}
	if enable && wake {
		if err := WakeScreen(deviceID); err != nil {
			fmt.Printf("⚠️  Could not wake device %s\n", deviceID)
		}
	}
}
