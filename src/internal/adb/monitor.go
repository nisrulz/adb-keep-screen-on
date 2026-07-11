package adb

import (
	"context"
	"slices"
	"strings"
	"time"
)

// CheckAdb returns an error if adb is not found in PATH.
func CheckAdb() error {
	_, err := execLookPath("adb")
	return err
}

// GetStayAwakeSetting retrieves the current value of stay_on_while_plugged_in from the device
func GetStayAwakeSetting(deviceID string) (string, error) {
	out, err := execCommand("adb", "-s", deviceID, "shell", "settings", "get", "global", "stay_on_while_plugged_in").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// WatchDevices polls `adb devices` every interval and calls onConnect/onDisconnect
// whenever a device's connection state changes. Stops when ctx is cancelled.
func WatchDevices(ctx context.Context, interval time.Duration, onConnect func(string), onDisconnect func(string)) error {
	prevDevices := ListDevices()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			curDevices := ListDevices()

			for _, d := range curDevices {
				if !slices.Contains(prevDevices, d) {
					onConnect(d)
				}
			}
			for _, d := range prevDevices {
				if !slices.Contains(curDevices, d) {
					onDisconnect(d)
				}
			}

			prevDevices = curDevices
		}
	}
}
