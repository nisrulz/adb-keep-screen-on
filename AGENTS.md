# AGENTS.md — adb-keep-screen-on

## Project Structure

```
src/
  main.go                         — CLI entry point, daemon + event loop
  main_test.go                    — Unit tests for main helpers
  internal/
    signal_handler.go             — SIGINT/SIGTERM handler, restores settings on exit
    config/
      config.go                   — JSON config load/save + defaults (~/.adb-keep-screen-on/config.json)
    adb/
      exec.go                     — exec.Command/exec.LookPath vars (injectable for tests)
      devices.go                  — ListDevices: parses `adb devices` output
      monitor.go                  — CheckAdb, GetStayAwakeSetting, WatchDevices
      toggle.go                   — SetStayAwakeValue, WakeScreen, SetStayAwake
      adb_test.go                 — Unit tests for adb package
      config_test.go              — Unit tests for config package
Makefile                        — Build, install, test, vet targets
build.sh                        — Cross-platform binary build
install.sh                      — Download-and-install from GitHub releases (Unix + Windows Git Bash)
.goreleaser.yaml                — GoReleaser config for CI releases
.github/workflows/release.yml   — GitHub Actions: release on v* tags
```

## Commands

```
adb-keep-screen-on               — Start daemon (runs in background)
adb-keep-screen-on stop          — Stop the daemon
adb-keep-screen-on status        — Show whether the daemon is running
adb-keep-screen-on --foreground  — Run in foreground (debug)
adb-keep-screen-on version       — Print the version
adb-keep-screen-on help          — Show usage help

make start / make stop / make status / make run
```

Daemon logs to `~/.adb-keep-screen-on/debug.log` and saves PID to `~/.adb-keep-screen-on/PID`.

Behavior is configurable via `~/.adb-keep-screen-on/config.json` (poll interval,
wake-on-connect, restore-on-disconnect, device allowlist). Polls `adb devices`
to detect changes at the configured interval. Already-connected devices are
detected immediately on startup.

## Build & Test

```
make build              — Compile (injects version via -ldflags)
make vet                — Static analysis
make test               — Run all tests
make build-all          — Cross-compile for all platforms
make run                — Build and run in foreground
make status             — Build and show daemon status
```

## Release

Tag and push to trigger the GitHub Actions release workflow:

```
git tag v0.1.0
git push origin v0.1.0
```

This runs GoReleaser, which builds binaries for macOS/Linux/Windows (amd64+arm64)
and creates a GitHub Release with checksums.

## Conventions

- No external dependencies (stdlib only)
- `execCommand`/`execLookPath` package vars in `internal/adb/exec.go` enable mocking in tests via `fakeCmd` helpers
- Pure logic extracted into injectable helpers (e.g. `storeOriginalSetting` accepts a `settingFetcher`)
- Logging uses `fmt.Print`/`fmt.Printf`/`fmt.Println` throughout
