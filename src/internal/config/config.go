// Package config loads user configuration from a JSON file, falling back to
// sensible defaults and auto-creating the file on first run.
package config

import (
	"encoding/json"
	"os"
	"slices"
	"time"
)

// Config holds all user-tunable behavior. Fields map to config.json keys.
type Config struct {
	// PollIntervalSeconds is how often `adb devices` is polled. Minimum 1.
	PollIntervalSeconds int `json:"poll_interval_seconds"`
	// WakeOnConnect wakes the screen when a device connects.
	WakeOnConnect bool `json:"wake_on_connect"`
	// RestoreOnDisconnect restores the original setting when a device
	// disconnects or the tool exits.
	RestoreOnDisconnect bool `json:"restore_on_disconnect"`
	// Devices is an optional allowlist of device IDs to manage. Empty means
	// all connected devices are managed.
	Devices []string `json:"devices"`
}

// Default returns the built-in configuration.
func Default() Config {
	return Config{
		PollIntervalSeconds: 1,
		WakeOnConnect:       true,
		RestoreOnDisconnect: true,
		Devices:             []string{},
	}
}

// PollInterval returns the poll interval as a duration, guarding against
// non-positive values.
func (c Config) PollInterval() time.Duration {
	if c.PollIntervalSeconds < 1 {
		return time.Second
	}
	return time.Duration(c.PollIntervalSeconds) * time.Second
}

// Allowed reports whether the given device should be managed under this config.
func (c Config) Allowed(deviceID string) bool {
	if len(c.Devices) == 0 {
		return true
	}
	return slices.Contains(c.Devices, deviceID)
}

// Load reads config from path. If the file does not exist it is created with
// defaults. On a malformed file it returns the defaults along with the error so
// callers can warn but keep running.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg := Default()
		return cfg, Save(path, cfg)
	}
	if err != nil {
		return Default(), err
	}

	cfg := Default()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), err
	}
	return cfg, nil
}

// Save writes config to path as indented JSON.
func Save(path string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}
