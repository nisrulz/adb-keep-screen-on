package adb

import "strings"

// ListDevices returns a slice of all connected device IDs
func ListDevices() []string {
	out, err := execCommand("adb", "devices").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	var devices []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == "device" {
			devices = append(devices, fields[0])
		}
	}
	return devices
}
