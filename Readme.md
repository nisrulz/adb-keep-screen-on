# 📱 ADB Keep Screen On

**ADB Keep Screen On** is a lightweight CLI tool written in Go that prevents your Android device from sleeping when connected via ADB over USB.

It’s ideal for developers, testers, and presenters who need the screen to stay awake during debugging, automation, or live demos without relying on charging state or manually tweaking device settings.

## ❓ Why This Tool Exists

Android’s developer setting **Stay awake while charging** is too broad and doesn’t cover modern workflows like wireless ADB or USB only connections. Developers have been requesting a more precise option **Stay awake while ADB is connected** since **2016**, but the issue remains unresolved:

🔗 [Google Issue Tracker #37094654](https://issuetracker.google.com/issues/37094654)

This tool fills that gap by monitoring ADB connection status and toggling the screen-on setting only when needed.

## ✨ Features

- Checks and stores the original stay awake setting before starting, restores it on exit.
- Wakes up the screen when enabling the stay awake debug setting.
- Monitors ADB connection and toggles settings automatically.
- Provides clear CLI feedback and graceful exit (restores settings on Ctrl+C).
- Works with USB ADB connections only.
- Lightweight, fast, and easy to use.

## ⚙️ Prerequisites

- ✅ [Go](https://golang.org/dl/) installed (for building the binary)
- ✅ [ADB (Android Debug Bridge)](https://developer.android.com/tools/adb) installed and accessible in your terminal

### 🧪 Install ADB via Homebrew (macOS)

If ADB is not installed, you can install it using Homebrew:

```sh
brew install --cask android-platform-tools
```

## 🚀 Usage

### 1. Build the binary

```sh
./build.sh
```

The binary will be created in the `dist` directory.

### 2. (Optional) Install globally

To run the tool from anywhere, use the provided install script:

```sh
./install.sh
```

This will symlink the binary to `~/bin/adb-keep-screen-on`. If `~/bin` is not in your PATH, follow the instructions printed by the script to add it.

### 3. Run the tool

Ensure your Android device is connected and USB debugging is enabled.

To run the built binary directly:

```sh
./dist/adb-keep-screen-on
```

Or specify a custom polling interval:

```sh
./dist/adb-keep-screen-on --interval 10
```

If you used the install script and symlinked the tool, simply open a new terminal and run:

```sh
adb-keep-screen-on
```

Or with a custom interval:

```sh
adb-keep-screen-on --interval 10
```

If you do not pass the --interval argument, the default polling interval of 10 seconds will be used.

`--interval` (optional): polling interval in seconds (default: `10`)

### Demo

<video src="./assets/demo.mp4" controls autoplay width="600">
  Your browser does not support the video tag. You can download the demo <a href="./assets/demo.mp4">here</a>.
</video>

## 📄 License

[Apache License Version 2.0 © Nishant Srivastava](/LICENSE)
