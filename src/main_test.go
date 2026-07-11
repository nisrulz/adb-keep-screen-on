package main

import (
	"errors"
	"testing"

	"adb-keep-screen-on/src/internal"
)

func TestStoreOriginalSetting_stores_if_not_present(t *testing.T) {
	settings := make(map[string]string)
	fetch := func(id string) (string, error) {
		return "1", nil
	}
	storeOriginalSetting("device-1", settings, fetch)
	if settings["device-1"] != "1" {
		t.Errorf("expected \"1\", got %q", settings["device-1"])
	}
}

func TestStoreOriginalSetting_skips_if_already_present(t *testing.T) {
	settings := map[string]string{"device-1": "0"}
	called := false
	fetch := func(id string) (string, error) {
		called = true
		return "1", nil
	}
	storeOriginalSetting("device-1", settings, fetch)
	if called {
		t.Error("fetch should not be called when setting already exists")
	}
	if settings["device-1"] != "0" {
		t.Errorf("expected \"0\", got %q", settings["device-1"])
	}
}

func TestStoreOriginalSetting_skips_on_fetch_error(t *testing.T) {
	settings := make(map[string]string)
	fetch := func(id string) (string, error) {
		return "", errors.New("adb error")
	}
	storeOriginalSetting("device-1", settings, fetch)
	if _, ok := settings["device-1"]; ok {
		t.Error("should not store setting on fetch error")
	}
}

func TestStoreOriginalSetting_stores_multiple_devices(t *testing.T) {
	settings := make(map[string]string)
	storeOriginalSetting("device-a", settings, func(id string) (string, error) {
		return "1", nil
	})
	storeOriginalSetting("device-b", settings, func(id string) (string, error) {
		return "0", nil
	})
	if settings["device-a"] != "1" {
		t.Errorf("expected \"1\", got %q", settings["device-a"])
	}
	if settings["device-b"] != "0" {
		t.Errorf("expected \"0\", got %q", settings["device-b"])
	}
}

func TestStoreOriginalSetting_passes_device_id(t *testing.T) {
	var calledWith string
	storeOriginalSetting("my-device", make(map[string]string), func(id string) (string, error) {
		calledWith = id
		return "1", nil
	})
	if calledWith != "my-device" {
		t.Errorf("expected \"my-device\", got %q", calledWith)
	}
}

func TestStoreOriginalSetting_retries_after_failure(t *testing.T) {
	settings := make(map[string]string)
	attempts := 0
	fetch := func(id string) (string, error) {
		attempts++
		if attempts == 1 {
			return "", errors.New("adb error")
		}
		return "1", nil
	}

	storeOriginalSetting("device-1", settings, fetch)
	if _, ok := settings["device-1"]; ok {
		t.Error("should not store on first failure")
	}

	storeOriginalSetting("device-1", settings, fetch)
	if settings["device-1"] != "1" {
		t.Errorf("expected \"1\" after retry, got %q", settings["device-1"])
	}
}

func TestOriginalSettings_starts_nil(t *testing.T) {
	if internal.OriginalSettings != nil {
		t.Error("OriginalSettings should be nil before assignment by main")
	}
}
