package adb

import (
	"os/exec"
	"strings"
)

// ListDevices returns a slice of all connected device IDs
func ListDevices() []string {
	out, err := exec.Command("adb", "devices").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	var devices []string
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == "device" {
			devices = append(devices, fields[0])
		}
	}
	return devices
}
