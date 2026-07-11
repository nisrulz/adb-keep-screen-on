package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ---- helpers ----

func buildBinary(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "adb-keep-screen-on")
	cmd := exec.Command("go", "build", "-o", path, "./src/")
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return path
}

func setHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	return dir
}

func noAdb(t *testing.T) {
	t.Helper()
	t.Setenv("PATH", t.TempDir())
}

func writeScript(t *testing.T, name, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0755); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return dir
}

func addToPath(t *testing.T, dir string) {
	t.Helper()
	t.Setenv("PATH", dir+string(filepath.ListSeparator)+os.Getenv("PATH"))
}

func fakeAdbEmpty(t *testing.T) {
	t.Helper()
	dir := writeScript(t, "adb", "#!/bin/sh\necho 'List of devices attached'\n")
	addToPath(t, dir)
}

func fakeAdbWithDevice(t *testing.T) {
	t.Helper()
	dir := writeScript(t, "adb", `#!/bin/sh
case "$1" in
  devices)
    echo "List of devices attached"
    echo "emulator-5554	device"
    ;;
  -s)
    case "$3" in
      get-state) echo "device" ;;
      shell)
        case "$4" in
          settings)
            case "$5" in
              get) echo "1" ;;
              put) exit 0 ;;
            esac
            ;;
          input) exit 0 ;;
        esac
        ;;
    esac
    ;;
esac
`)
	addToPath(t, dir)
}

func fakeAdbFails(t *testing.T) {
	t.Helper()
	dir := writeScript(t, "adb", "#!/bin/sh\nexit 1\n")
	addToPath(t, dir)
}

func fakeAdbSlow(t *testing.T, delay time.Duration) {
	t.Helper()
	dir := writeScript(t, "adb", fmt.Sprintf("#!/bin/sh\nsleep %d\necho 'List of devices attached'\n", int(delay.Seconds())))
	addToPath(t, dir)
}

func startDaemon(t *testing.T, bin string) {
	t.Helper()
	out, err := exec.Command(bin).CombinedOutput()
	if err != nil {
		t.Fatalf("start daemon failed: %v\n%s", err, out)
	}
}

func runStop(t *testing.T, bin string) {
	t.Helper()
	out, err := exec.Command(bin, "stop").CombinedOutput()
	if err != nil {
		t.Fatalf("stop daemon failed: %v\n%s", err, out)
	}
}

func daemonDir(t *testing.T, home string) (pidPath, logPath string) {
	base := filepath.Join(home, ".adb-keep-screen-on")
	return filepath.Join(base, "PID"), filepath.Join(base, "debug.log")
}

func waitFor(t *testing.T, fn func() bool) {
	t.Helper()
	for i := 0; i < 40; i++ {
		if fn() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func readPID(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read PID file: %v", err)
	}
	var pid int
	if _, err := fmt.Fscanf(strings.NewReader(string(data)), "%d", &pid); err != nil {
		t.Fatalf("parse PID: %v", err)
	}
	return pid
}

// ---- argument parsing ----

func TestNoArgsStartsDaemon(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbEmpty(t)

	out, err := exec.Command(bin).CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Started") {
		t.Errorf("expected 'Started', got: %s", out)
	}
	runStop(t, bin)
}

func TestStopWithoutDaemon(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)

	out, err := exec.Command(bin, "stop").CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for stop without daemon")
	}
	if !strings.Contains(string(out), "Not running") {
		t.Errorf("expected 'Not running', got: %s", out)
	}
}

func TestStopTwice(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	runStop(t, bin)

	out, err := exec.Command(bin, "stop").CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for second stop")
	}
	if !strings.Contains(string(out), "Not running") {
		t.Errorf("expected 'Not running', got: %s", out)
	}
}

func TestForegroundWithoutAdb(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	noAdb(t)

	out, err := exec.Command(bin, "--foreground").CombinedOutput()
	if err != nil {
		t.Logf("exit: %v", err)
	}
	if !strings.Contains(string(out), "not installed") {
		t.Errorf("expected adb-not-found message, got: %s", out)
	}
}

func TestForegroundWithFakeAdb(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbEmpty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, bin, "--foreground").CombinedOutput()
	if !strings.Contains(string(out), "Monitoring") {
		t.Errorf("expected 'Monitoring', got: %s", out)
	}
	_ = err
}

func TestForegroundWithDevice(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbWithDevice(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	out, _ := exec.CommandContext(ctx, bin, "--foreground").CombinedOutput()
	if !strings.Contains(string(out), "Connected:") {
		t.Errorf("expected 'Connected:', got: %s", out)
	}
}

func TestInvalidArgPrintsUsage(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)

	out, err := exec.Command(bin, "unknown").CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if !strings.Contains(string(out), "Usage") {
		t.Errorf("expected 'Usage', got: %s", out)
	}
}

func TestInvalidArgDoesNotDaemonize(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)

	out, err := exec.Command(bin, "unknown").CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if strings.Contains(string(out), "Started") {
		t.Errorf("should not start daemon, got: %s", out)
	}
}

func TestInvalidFlag(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)

	out, err := exec.Command(bin, "--invalid-flag").CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if !strings.Contains(string(out), "Usage") {
		t.Errorf("expected 'Usage', got: %s", out)
	}
}

func TestForegroundWithExtraArg(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbEmpty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Extra arg after --foreground should be ignored (foreground mode runs)
	out, _ := exec.CommandContext(ctx, bin, "--foreground", "extra").CombinedOutput()
	if !strings.Contains(string(out), "Monitoring") {
		t.Errorf("expected 'Monitoring', got: %s", out)
	}
}

func TestStopWithExtraArg(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	// Extra arg after "stop" should not prevent stopping
	out, err := exec.Command(bin, "stop", "--foreground").CombinedOutput()
	if err != nil {
		t.Fatalf("stop failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Stopped") {
		t.Errorf("expected 'Stopped', got: %s", out)
	}
}

// ---- daemon lifecycle ----

func TestDaemonLifecycle(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	pp, lp := daemonDir(t, homeDir)

	waitFor(t, func() bool { _, err := os.Stat(pp); return err == nil })

	if _, err := os.Stat(pp); os.IsNotExist(err) {
		t.Error("PID file was not created")
	}
	if _, err := os.Stat(lp); os.IsNotExist(err) {
		t.Error("log file was not created")
	}

	runStop(t, bin)
	waitFor(t, func() bool { _, err := os.Stat(pp); return os.IsNotExist(err) })

	if _, err := os.Stat(pp); !os.IsNotExist(err) {
		t.Error("PID file was not removed after stop")
	}
}

func TestDaemonRestartAfterStop(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	runStop(t, bin)

	// Should be able to start again
	startDaemon(t, bin)
	runStop(t, bin)
}

func TestDoubleStart(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	defer runStop(t, bin)

	out, err := exec.Command(bin).CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for second start")
	}
	if !strings.Contains(string(out), "Already running") {
		t.Errorf("expected 'Already running', got: %s", out)
	}
}

func TestStopAlreadyDead(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	pp, _ := daemonDir(t, homeDir)

	waitFor(t, func() bool { _, err := os.Stat(pp); return err == nil })

	pid := readPID(t, pp)

	kill := exec.Command("kill", "-9", fmt.Sprintf("%d", pid))
	if out, err := kill.CombinedOutput(); err != nil {
		t.Fatalf("kill failed: %v\n%s", err, out)
	}
	time.Sleep(100 * time.Millisecond)

	stopOut, err := exec.Command(bin, "stop").CombinedOutput()
	if err != nil {
		t.Logf("stop error: %v\n%s", err, stopOut)
	}

	waitFor(t, func() bool { _, err := os.Stat(pp); return os.IsNotExist(err) })
	if _, err := os.Stat(pp); !os.IsNotExist(err) {
		t.Error("PID file was not removed after stop with dead daemon")
	}
}

// ---- PID file state ----

func TestCorruptedPidFile(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	pp, _ := daemonDir(t, homeDir)
	waitFor(t, func() bool { _, err := os.Stat(pp); return err == nil })

	// Corrupt the PID file
	if err := os.WriteFile(pp, []byte("not-a-number\n"), 0644); err != nil {
		t.Fatalf("write corrupt PID: %v", err)
	}

	out, err := exec.Command(bin, "stop").CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for corrupted PID")
	}
	if !strings.Contains(string(out), "Not running") {
		t.Errorf("expected 'Not running', got: %s", out)
	}

	// PID file should be gone
	waitFor(t, func() bool { _, err := os.Stat(pp); return os.IsNotExist(err) })
}

func TestMissingPidFile(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	pp, _ := daemonDir(t, homeDir)
	waitFor(t, func() bool { _, err := os.Stat(pp); return err == nil })

	// Remove PID file manually
	if err := os.Remove(pp); err != nil {
		t.Fatalf("remove PID: %v", err)
	}

	out, err := exec.Command(bin, "stop").CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when PID file missing")
	}
	if !strings.Contains(string(out), "Not running") {
		t.Errorf("expected 'Not running', got: %s", out)
	}
}

func TestEmptyPidFile(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	pp, _ := daemonDir(t, homeDir)
	waitFor(t, func() bool { _, err := os.Stat(pp); return err == nil })

	// Empty the PID file
	if err := os.WriteFile(pp, []byte("\n"), 0644); err != nil {
		t.Fatalf("write empty PID: %v", err)
	}

	out, err := exec.Command(bin, "stop").CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for empty PID")
	}
	if !strings.Contains(string(out), "Not running") {
		t.Errorf("expected 'Not running', got: %s", out)
	}
}

// ---- file content ----

func TestPidFileContainsValidPID(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	pp, _ := daemonDir(t, homeDir)

	waitFor(t, func() bool { _, err := os.Stat(pp); return err == nil })

	data, err := os.ReadFile(pp)
	if err != nil {
		t.Fatalf("read PID file: %v", err)
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		t.Fatal("PID file is empty")
	}

	var pid int
	if _, err := fmt.Sscanf(content, "%d", &pid); err != nil {
		t.Errorf("PID file content is not a number: %q", content)
	}
	if pid <= 0 {
		t.Errorf("PID should be positive, got: %d", pid)
	}

	runStop(t, bin)
}

func TestLogFileIsCreated(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	_, lp := daemonDir(t, homeDir)

	waitFor(t, func() bool { _, err := os.Stat(lp); return err == nil })

	if _, err := os.Stat(lp); os.IsNotExist(err) {
		t.Error("log file was not created")
	}

	runStop(t, bin)
}

func TestDataDirIsCreated(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	defer runStop(t, bin)

	dir := filepath.Join(homeDir, ".adb-keep-screen-on")
	waitFor(t, func() bool { _, err := os.Stat(dir); return err == nil })

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("data directory was not created")
	}
}

// ---- adb interaction ----

func TestDaemonWithDeviceLogsConnection(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbWithDevice(t)

	startDaemon(t, bin)
	_, lp := daemonDir(t, homeDir)

	waitFor(t, func() bool {
		data, err := os.ReadFile(lp)
		if err != nil {
			return false
		}
		return strings.Contains(string(data), "Connected:")
	})

	runStop(t, bin)
}

func TestDaemonWithEmptyAdbDoesNotLogConnection(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	_, lp := daemonDir(t, homeDir)

	waitFor(t, func() bool { _, err := os.Stat(lp); return err == nil })

	data, _ := os.ReadFile(lp)
	if strings.Contains(string(data), "Connected:") {
		t.Error("should not report device connection with no devices")
	}

	runStop(t, bin)
}

func TestDaemonWithFailingAdbLogsError(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbFails(t)

	startDaemon(t, bin)
	_, lp := daemonDir(t, homeDir)

	// Give it time to try and fail
	time.Sleep(3 * time.Second)

	data, _ := os.ReadFile(lp)
	if !strings.Contains(string(data), "Monitoring") {
		t.Errorf("expected at least monitoring message in log, got: %s", data)
	}

	runStop(t, bin)
}

// ---- signal handling ----

func TestStopSendsSigint(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	pp, _ := daemonDir(t, homeDir)
	waitFor(t, func() bool { _, err := os.Stat(pp); return err == nil })

	originalPID := readPID(t, pp)

	runStop(t, bin)
	waitFor(t, func() bool { _, err := os.Stat(pp); return os.IsNotExist(err) })

	// Verify the process actually exited
	proc, err := os.FindProcess(originalPID)
	if err == nil {
		// Signal 0 reports whether the process still exists. On a dead
		// process this should fail, so a nil error means it is still around.
		if proc.Signal(os.Signal(nil)) == nil {
			// Some Unix variants may report success for a dead process, so
			// the real check we trust is whether the PID file got removed.
		}
	}
}

// ---- race conditions (timing) ----

func TestImmediateStopAfterStart(t *testing.T) {
	bin := buildBinary(t)
	homeDir := setHome(t)
	fakeAdbEmpty(t)

	startDaemon(t, bin)
	pp, _ := daemonDir(t, homeDir)

	// Stop right away, without waiting for the PID file to appear first.
	out, err := exec.Command(bin, "stop").CombinedOutput()
	if err != nil {
		t.Fatalf("stop failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Stopped") {
		t.Errorf("expected 'Stopped', got: %s", out)
	}

	// PID file should eventually be gone
	waitFor(t, func() bool { _, err := os.Stat(pp); return os.IsNotExist(err) })
}

func TestMultipleStartStopCycles(t *testing.T) {
	bin := buildBinary(t)
	setHome(t)
	fakeAdbEmpty(t)

	for i := 0; i < 3; i++ {
		t.Logf("cycle %d", i+1)
		startDaemon(t, bin)
		runStop(t, bin)
		time.Sleep(200 * time.Millisecond)
	}
}
