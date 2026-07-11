package adb

import (
	"os/exec"
	"slices"
	"strings"
	"testing"
)

func fakeCmd(output string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cmd := exec.Command("cat")
		cmd.Stdin = strings.NewReader(output)
		return cmd
	}
}

func fakeCmdSuccess() func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
}

func fakeCmdFailure() func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
}

type cmdCall struct {
	name string
	args []string
}

func recordingCmd(calls *[]cmdCall, success bool) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		*calls = append(*calls, cmdCall{name, args})
		if success {
			return exec.Command("true")
		}
		return exec.Command("false")
	}
}

func resetCommand() {
	execCommand = exec.Command
	execLookPath = exec.LookPath
}

func TestListDevices_returns_devices(t *testing.T) {
	execCommand = fakeCmd("List of devices attached\nemulator-5554\tdevice\n")
	t.Cleanup(resetCommand)

	devices := ListDevices()
	if len(devices) != 1 || devices[0] != "emulator-5554" {
		t.Errorf("expected [emulator-5554], got %v", devices)
	}
}

func TestListDevices_multiple(t *testing.T) {
	execCommand = fakeCmd("List of devices attached\nemulator-5554\tdevice\nemulator-5556\tdevice\n")
	t.Cleanup(resetCommand)

	devices := ListDevices()
	if len(devices) != 2 {
		t.Fatalf("expected 2 devices, got %v", devices)
	}
	if devices[0] != "emulator-5554" || devices[1] != "emulator-5556" {
		t.Errorf("expected [emulator-5554 emulator-5556], got %v", devices)
	}
}

func TestListDevices_empty(t *testing.T) {
	execCommand = fakeCmd("List of devices attached\n\n")
	t.Cleanup(resetCommand)

	devices := ListDevices()
	if len(devices) != 0 {
		t.Errorf("expected empty list, got %v", devices)
	}
}

func TestListDevices_skips_unauthorized(t *testing.T) {
	execCommand = fakeCmd("List of devices attached\nemulator-5554\tunauthorized\n")
	t.Cleanup(resetCommand)

	devices := ListDevices()
	if len(devices) != 0 {
		t.Errorf("expected [], got %v", devices)
	}
}

func TestListDevices_skips_offline(t *testing.T) {
	execCommand = fakeCmd("List of devices attached\nemulator-5554\toffline\n")
	t.Cleanup(resetCommand)

	devices := ListDevices()
	if len(devices) != 0 {
		t.Errorf("expected [], got %v", devices)
	}
}

func TestListDevices_mixed_states(t *testing.T) {
	execCommand = fakeCmd("List of devices attached\ndevice1\tdevice\ndevice2\toffline\ndevice3\tunauthorized\n")
	t.Cleanup(resetCommand)

	devices := ListDevices()
	if len(devices) != 1 || devices[0] != "device1" {
		t.Errorf("expected [device1], got %v", devices)
	}
}

func TestListDevices_returns_nil_on_error(t *testing.T) {
	execCommand = fakeCmdFailure()
	t.Cleanup(resetCommand)

	devices := ListDevices()
	if devices != nil {
		t.Errorf("expected nil, got %v", devices)
	}
}

func TestCheckAdb_true(t *testing.T) {
	t.Cleanup(resetCommand)

	if err := CheckAdb(); err != nil {
		t.Errorf("expected nil when adb in PATH, got %v", err)
	}
}

func TestCheckAdb_false(t *testing.T) {
	execLookPath = func(name string) (string, error) {
		return "", exec.ErrNotFound
	}
	t.Cleanup(resetCommand)

	if err := CheckAdb(); err == nil {
		t.Error("expected error when adb not in PATH")
	}
}

func TestGetStayAwakeSetting_value_1(t *testing.T) {
	execCommand = fakeCmd("1\n")
	t.Cleanup(resetCommand)

	val, err := GetStayAwakeSetting("emulator-5554")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "1" {
		t.Errorf("expected \"1\", got %q", val)
	}
}

func TestGetStayAwakeSetting_value_0(t *testing.T) {
	execCommand = fakeCmd("0\n")
	t.Cleanup(resetCommand)

	val, err := GetStayAwakeSetting("emulator-5554")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "0" {
		t.Errorf("expected \"0\", got %q", val)
	}
}

func TestGetStayAwakeSetting_empty(t *testing.T) {
	execCommand = fakeCmd("\n")
	t.Cleanup(resetCommand)

	val, err := GetStayAwakeSetting("emulator-5554")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty, got %q", val)
	}
}

func TestGetStayAwakeSetting_error(t *testing.T) {
	execCommand = fakeCmdFailure()
	t.Cleanup(resetCommand)

	_, err := GetStayAwakeSetting("emulator-5554")
	if err == nil {
		t.Error("expected error")
	}
}

func TestSetStayAwakeValue_success(t *testing.T) {
	execCommand = fakeCmdSuccess()
	t.Cleanup(resetCommand)

	if err := SetStayAwakeValue("emulator-5554", "1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetStayAwakeValue_error(t *testing.T) {
	execCommand = fakeCmdFailure()
	t.Cleanup(resetCommand)

	if err := SetStayAwakeValue("emulator-5554", "1"); err == nil {
		t.Error("expected error")
	}
}

func TestSetStayAwake_enable_calls_setting_and_wake(t *testing.T) {
	var calls []cmdCall
	execCommand = recordingCmd(&calls, true)
	t.Cleanup(resetCommand)

	SetStayAwake("device-1", true, true)

	if len(calls) != 2 {
		t.Fatalf("expected 2 exec calls, got %d", len(calls))
	}

	// First call: SetStayAwakeValue with "1"
	if calls[0].name != "adb" {
		t.Errorf("expected command 'adb', got %q", calls[0].name)
	}
	if !slices.Contains(calls[0].args, "1") {
		t.Errorf("expected setting value '1' in args, got %v", calls[0].args)
	}

	// Second call: WakeScreen (KEYCODE_WAKEUP)
	if calls[1].name != "adb" {
		t.Errorf("expected command 'adb', got %q", calls[1].name)
	}
	if !slices.Contains(calls[1].args, "KEYCODE_WAKEUP") {
		t.Errorf("expected KEYCODE_WAKEUP in args, got %v", calls[1].args)
	}
}

func TestSetStayAwake_enable_without_wake_skips_wake(t *testing.T) {
	var calls []cmdCall
	execCommand = recordingCmd(&calls, true)
	t.Cleanup(resetCommand)

	SetStayAwake("device-1", true, false)

	if len(calls) != 1 {
		t.Fatalf("expected 1 exec call when wake disabled, got %d", len(calls))
	}
	if slices.Contains(calls[0].args, "KEYCODE_WAKEUP") {
		t.Error("should not wake screen when wake is false")
	}
}

func TestSetStayAwake_disable_skips_wake(t *testing.T) {
	var calls []cmdCall
	execCommand = recordingCmd(&calls, true)
	t.Cleanup(resetCommand)

	SetStayAwake("device-1", false, true)

	if len(calls) != 1 {
		t.Fatalf("expected 1 exec call when disabling, got %d", len(calls))
	}

	if !slices.Contains(calls[0].args, "0") {
		t.Errorf("expected setting value '0' in args, got %v", calls[0].args)
	}
	if slices.Contains(calls[0].args, "KEYCODE_WAKEUP") {
		t.Error("should not wake screen when disabling")
	}
}

func TestSetStayAwake_continues_on_setting_error(t *testing.T) {
	var calls []cmdCall
	execCommand = recordingCmd(&calls, false)
	t.Cleanup(resetCommand)

	SetStayAwake("device-1", true, true)

	// Even though first command fails, it should still attempt WakeScreen
	if len(calls) != 2 {
		t.Fatalf("expected 2 exec calls even on error, got %d", len(calls))
	}
}
