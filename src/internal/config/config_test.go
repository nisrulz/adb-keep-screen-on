package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestPollInterval_defaults_to_one_second_when_non_positive(t *testing.T) {
	for _, v := range []int{0, -5} {
		got := Config{PollIntervalSeconds: v}.PollInterval()
		if got != time.Second {
			t.Errorf("PollIntervalSeconds=%d: expected 1s, got %s", v, got)
		}
	}
}

func TestPollInterval_uses_configured_value(t *testing.T) {
	if got := (Config{PollIntervalSeconds: 5}).PollInterval(); got != 5*time.Second {
		t.Errorf("expected 5s, got %s", got)
	}
}

func TestAllowed_empty_list_allows_all(t *testing.T) {
	if !(Config{}).Allowed("any-device") {
		t.Error("empty allowlist should allow all devices")
	}
}

func TestAllowed_respects_allowlist(t *testing.T) {
	cfg := Config{Devices: []string{"device-a"}}
	if !cfg.Allowed("device-a") {
		t.Error("device-a should be allowed")
	}
	if cfg.Allowed("device-b") {
		t.Error("device-b should not be allowed")
	}
}

func TestLoad_creates_default_file_when_missing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(cfg, Default()) {
		t.Errorf("expected defaults, got %+v", cfg)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected config file to be created: %v", err)
	}
}

func TestLoad_reads_existing_file(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"poll_interval_seconds":3,"wake_on_connect":false}`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollIntervalSeconds != 3 {
		t.Errorf("expected 3, got %d", cfg.PollIntervalSeconds)
	}
	if cfg.WakeOnConnect {
		t.Error("expected wake_on_connect false")
	}
	if !cfg.RestoreOnDisconnect {
		t.Error("omitted restore_on_disconnect should keep its default true")
	}
}

func TestLoad_returns_defaults_on_malformed_file(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("{not json"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err == nil {
		t.Error("expected error for malformed config")
	}
	if !reflect.DeepEqual(cfg, Default()) {
		t.Errorf("expected defaults on error, got %+v", cfg)
	}
}

func TestSaveLoad_roundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	want := Config{PollIntervalSeconds: 2, WakeOnConnect: false, RestoreOnDisconnect: false, Devices: []string{"x"}}

	if err := Save(path, want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.PollIntervalSeconds != want.PollIntervalSeconds || got.WakeOnConnect != want.WakeOnConnect ||
		got.RestoreOnDisconnect != want.RestoreOnDisconnect || len(got.Devices) != 1 || got.Devices[0] != "x" {
		t.Errorf("roundtrip mismatch: got %+v want %+v", got, want)
	}
}
