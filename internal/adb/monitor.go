package adb

import (
	"os/exec"
	"strings"
)

// IsDeviceConnected checks if an Android device is connected via ADB
func IsDeviceConnected(deviceID string) bool {
	out, err := exec.Command("adb", "-s", deviceID, "get-state").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "device"
}

// IsAdbAvailable checks if adb is available in PATH and prints instructions if not
func IsAdbAvailable() bool {
	_, err := exec.LookPath("adb")
	if err != nil {
		println("❌ 'adb' is not installed or not found in your PATH.\n")
		println("Please install Android Platform Tools from https://developer.android.com/tools and ensure 'adb' is available in your system PATH.\n")
		println("If you are on macOS or Linux, you can use Homebrew:\n")
		println("  brew install --cask android-platform-tools\n")
		return false
	}
	return true
}

// GetStayAwakeSetting retrieves the current value of stay_on_while_plugged_in from the device
func GetStayAwakeSetting(deviceID string) (string, error) {
	out, err := exec.Command("adb", "-s", deviceID, "shell", "settings", "get", "global", "stay_on_while_plugged_in").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
