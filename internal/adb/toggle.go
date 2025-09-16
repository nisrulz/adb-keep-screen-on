package adb

import (
	"fmt"
	"os/exec"
)

// SetStayAwake sets the device's stay awake setting while plugged in
func SetStayAwake(enable bool) {
	value := "0"
	if enable {
		value = "1"
	}
	// Attempt to set the stay awake value
	err := SetStayAwakeValue(value)
	if err != nil {
		fmt.Printf("⚠️ Failed to set stay_on_while_plugged_in to %s: %v\n", value, err)
	}
}

// SetStayAwakeValue sets stay_on_while_plugged_in to a specific value on the device
func SetStayAwakeValue(val string) error {
	cmd := exec.Command("adb", "shell", "settings", "put", "global", "stay_on_while_plugged_in", val)
	return cmd.Run()
}
