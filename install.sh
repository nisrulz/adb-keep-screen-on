#!/bin/bash

# Install script for adb-keep-screen-on
# Symlinks the built binary to ~/bin for global access (no sudo required)

set -e

BINARY_PATH="$(pwd)/dist/adb-keep-screen-on"
LINK_DIR="$HOME/bin"
LINK_PATH="$LINK_DIR/adb-keep-screen-on"

if [ ! -f "$BINARY_PATH" ]; then
  echo ""
  echo "⚙️ Binary not found! Running ./build.sh..."
  ./build.sh
  echo ""
fi

# Create ~/bin if it doesn't exist
if [ ! -d "$LINK_DIR" ]; then
  echo "📁 Creating $LINK_DIR directory..."
  mkdir -p "$LINK_DIR"
  echo ""
fi

if [ -L "$LINK_PATH" ] || [ -f "$LINK_PATH" ]; then
  echo "🔄 Removing existing symlink or binary"
  rm -f "$LINK_PATH"
  echo ""
fi

echo "🔗 Creating symlink"
ln -s "$BINARY_PATH" "$LINK_PATH"
echo ""

echo "✅ Installed adb-keep-screen-on."
echo ""
echo "You can now run 'adb-keep-screen-on' from anywhere"
echo "if $LINK_DIR is in your PATH."
echo ""

if ! echo "$PATH" | grep -q "$LINK_DIR"; then
  echo "=================================================================="
  echo ""
  echo "⚠️ $LINK_DIR is not in your PATH."
  echo ""
  echo "Add the following line to your shell profile (e.g., ~/.bashrc or ~/.zshrc):"
  echo ""
  echo "    export PATH=\"$LINK_DIR:\$PATH\""
  echo ""
  echo "After editing your profile, run:"
  echo "    source ~/.zshrc   # or source ~/.bashrc, depending on your shell"
fi
