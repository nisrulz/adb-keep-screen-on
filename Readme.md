# ADB Keep Screen On

**ADB Keep Screen On** is a lightweight CLI tool written in Go that prevents your Android device from sleeping when connected via ADB over USB.

> It sends a periodic ADB command to maintain screen wakefulness, ideal for development, debugging, testing, and live presentations where screen timeout could interrupt workflow.

## Usage

1. Build the binary:

   ```sh
   ./build.sh
   ```

   The binary will be created in the `dist` directory.

2. Run the tool:
    > Ensure you have ADB installed and your device is connected.

   ```sh
   ./dist/adb-keep-screen-on --interval 10
   ```

   > The optional `--interval` argument sets the polling interval in seconds (default: 10).

## License

[Apache License Version 2.0 © Nishant Srivastava](/LICENSE)
