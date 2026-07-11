# DEV.md - adb-keep-screen-on

## Prerequisites

- Go 1.25+
- ADB installed and in PATH

## Development

```sh
# Build
make build

# Run in foreground (with connected device)
make run

# Start daemon / stop daemon
make start
make stop

# Cross-platform release binaries
make build-all

# Run tests
make test

# Static analysis
make vet
```

## Architecture

Uses `adb devices` to detect device plug/unplug events. The poll interval,
wake-on-connect, restore-on-disconnect, and device allowlist are read from
`~/.adb-keep-screen-on/config.json` (auto-created with defaults on first run).
Already-connected devices are detected immediately on startup.

On SIGINT (Ctrl+C) in foreground mode or when `adb-keep-screen-on stop` is
called, the signal handler restores all tracked original settings (unless
`restore_on_disconnect` is disabled) and cleans up the PID file.

## Releasing

Tag a commit and push to trigger the release workflow:

```sh
git tag v0.1.0
git push origin v0.1.0
```

This runs GoReleaser via GitHub Actions, which builds binaries for macOS/Linux/Windows
(amd64+arm64) and creates a GitHub Release with checksums.

## Adding tests

- Functions calling `exec.Command` use the injectable `execCommand` var: assign `fakeCmd(output)` in tests
- `exec.LookPath` uses `execLookPath` var: assign a mock in tests
