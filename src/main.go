package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"adb-keep-screen-on/src/internal"
	"adb-keep-screen-on/src/internal/adb"
	"adb-keep-screen-on/src/internal/config"
)

const envDaemon = "ADB_KEEP_SCREEN_ON_DAEMON"

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "stop":
		os.Exit(stopDaemon())
	case "status":
		os.Exit(statusDaemon())
	case "help", "--help", "-h":
		printHelp()
		return
	case "version", "--version", "-v":
		fmt.Printf("adb-keep-screen-on %s\n", version)
		return
	case "", "--foreground":
		// fall through to run
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}

	foreground := command == "--foreground"

	if !foreground && os.Getenv(envDaemon) != "1" {
		os.Exit(daemonize())
	}

	run(foreground)
}

func run(foreground bool) {
	if err := adb.CheckAdb(); err != nil {
		fmt.Println("❌ 'adb' is not installed or not found in your PATH.")
		fmt.Println("   Install Android Platform Tools: https://developer.android.com/tools")
		fmt.Println("   macOS/Linux via Homebrew:")
		fmt.Println("     brew install --cask android-platform-tools")
		os.Exit(1)
	}

	os.MkdirAll(dataDir(), 0755)
	internal.PidPath = pidPath()

	cfg, err := config.Load(configPath())
	if err != nil {
		fmt.Printf("⚠️  Could not read config (%v) — using defaults.\n", err)
	}
	internal.RestoreOnExit = cfg.RestoreOnDisconnect

	fmt.Println("🔌 Monitoring ADB connection to keep the screen awake")
	fmt.Printf("   config:   %s\n", configPath())
	fmt.Printf("   interval: %s   wake-on-connect: %v   restore: %v\n", cfg.PollInterval(), cfg.WakeOnConnect, cfg.RestoreOnDisconnect)
	if len(cfg.Devices) > 0 {
		fmt.Printf("   devices:  %s (allowlist)\n", strings.Join(cfg.Devices, ", "))
	}
	fmt.Println()

	internal.HandleInterrupt()

	originalSettings := make(map[string]string)
	internal.OriginalSettings = originalSettings

	handleConnected := func(d string) {
		if !cfg.Allowed(d) {
			return
		}
		storeOriginalSetting(d, originalSettings, adb.GetStayAwakeSetting)
		adb.SetStayAwake(d, true, cfg.WakeOnConnect)
		fmt.Printf("✅ Connected: %s — screen will stay awake\n", d)
	}

	handleDisconnected := func(d string) {
		if !cfg.Allowed(d) {
			return
		}
		if val, ok := originalSettings[d]; ok {
			if cfg.RestoreOnDisconnect {
				adb.SetStayAwakeValue(d, val)
			}
			delete(originalSettings, d)
		}
		fmt.Printf("🔌 Disconnected: %s\n", d)
	}

	for _, d := range adb.ListDevices() {
		handleConnected(d)
	}

	adb.WatchDevices(context.Background(), cfg.PollInterval(), handleConnected, handleDisconnected)
}

func printHelp() {
	fmt.Print(`adb-keep-screen-on — keep Android screens awake while connected via ADB

Usage:
  adb-keep-screen-on [command]

Commands:
  (none)          Start the background daemon
  stop            Stop the running daemon
  status          Show whether the daemon is running
  --foreground    Run in the foreground (debug/logging)
  version         Print the version
  help            Show this help

Config (JSON, auto-created on first run):
  ` + configPath() + `
    poll_interval_seconds   How often to poll adb devices (default 1)
    wake_on_connect         Wake the screen on connect (default true)
    restore_on_disconnect   Restore original setting on disconnect/exit (default true)
    devices                 Optional allowlist of device IDs ([] = all)

Logs: ` + logPath() + `
`)
}

type settingFetcher func(string) (string, error)

func storeOriginalSetting(deviceID string, settings map[string]string, fetch settingFetcher) {
	if _, ok := settings[deviceID]; ok {
		return
	}
	setting, err := fetch(deviceID)
	if err == nil {
		settings[deviceID] = setting
	}
}

func dataDir() string    { return expandPath("~/.adb-keep-screen-on") }
func pidPath() string    { return filepath.Join(dataDir(), "PID") }
func logPath() string    { return filepath.Join(dataDir(), "debug.log") }
func configPath() string { return filepath.Join(dataDir(), "config.json") }

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func readPIDFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func processAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

func daemonize() int {
	os.MkdirAll(dataDir(), 0755)
	pp := pidPath()
	lp := logPath()

	if pid, err := readPIDFile(pp); err == nil && processAlive(pid) {
		fmt.Printf("ℹ️  Already running (PID %d)\n", pid)
		return 1
	}

	logFile, err := os.OpenFile(lp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("❌ Failed to open log file:", err)
		return 1
	}

	exe, err := os.Executable()
	if err != nil {
		fmt.Println("❌ Failed to get executable path:", err)
		return 1
	}

	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Env = append(os.Environ(), envDaemon+"=1")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ Failed to start daemon:", err)
		return 1
	}

	os.WriteFile(pp, []byte(strconv.Itoa(cmd.Process.Pid)+"\n"), 0644)
	fmt.Printf("✅ Started (PID %d)\n", cmd.Process.Pid)
	fmt.Printf("   Logs:   %s\n", lp)
	fmt.Printf("   Config: %s\n", configPath())
	fmt.Println("   Stop with: adb-keep-screen-on stop")
	return 0
}

func stopDaemon() int {
	pp := pidPath()

	pid, err := readPIDFile(pp)
	if err != nil {
		fmt.Println("ℹ️  Not running (no PID file)")
		return 1
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("ℹ️  Not running")
		os.Remove(pp)
		return 1
	}

	if err := proc.Signal(syscall.SIGINT); err != nil {
		fmt.Println("⚠️  Daemon not responding, removing PID file")
		os.Remove(pp)
		return 1
	}

	fmt.Printf("✅ Stopped (PID %d)\n", pid)
	os.Remove(pp)
	return 0
}

func statusDaemon() int {
	pid, err := readPIDFile(pidPath())
	if err != nil {
		fmt.Println("⚪ Not running")
		return 1
	}
	if !processAlive(pid) {
		fmt.Printf("⚪ Not running (stale PID file for %d)\n", pid)
		return 1
	}
	fmt.Printf("🟢 Running (PID %d)\n", pid)
	fmt.Printf("   Logs:   %s\n", logPath())
	fmt.Printf("   Config: %s\n", configPath())
	return 0
}
