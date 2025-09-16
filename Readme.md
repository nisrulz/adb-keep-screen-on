# 📱 ADB Keep Screen On

**ADB Keep Screen On** is a lightweight CLI tool written in Go that prevents your Android device from sleeping when connected via ADB over USB.

It’s ideal for developers, testers, and presenters who need the screen to stay awake during debugging, automation, or live demos without relying on charging state or manually tweaking device settings.

## ❓ Why This Tool Exists

Android’s developer setting **Stay awake while charging** is too broad and doesn’t cover modern workflows like wireless ADB or USB only connections. Developers have been requesting a more precise option **Stay awake while ADB is connected** since **2016**, but the issue remains unresolved:

🔗 [Google Issue Tracker #37094654](https://issuetracker.google.com/issues/37094654)

This tool fills that gap by monitoring ADB connection status and toggling the screen-on setting only when needed.

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

### 2. Run the tool

Ensure your Android device is connected and USB debugging is enabled.

```sh
./dist/adb-keep-screen-on --interval 10
```

`--interval` (optional): polling interval in seconds (default: `10`)

## 📄 License

[Apache License Version 2.0 © Nishant Srivastava](/LICENSE)
