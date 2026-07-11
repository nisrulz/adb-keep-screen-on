#!/bin/sh
set -eu

REPO="nisrulz/adb-keep-screen-on"
BIN="adb-keep-screen-on"

arch=$(uname -m)
os_raw=$(uname -s)
os=$(echo "$os_raw" | tr '[:upper:]' '[:lower:]')

case "$arch" in
  x86_64 | amd64) arch="amd64" ;;
  aarch64 | arm64) arch="arm64" ;;
  *)
    echo "Unsupported architecture: $arch"
    exit 1
    ;;
esac

case "$os" in
  darwin | linux) ;;
  mingw* | msys* | cygwin*) os="windows" ;;
  *)
    echo "Unsupported OS: $os_raw"
    exit 1
    ;;
esac

tag=$(curl -sfL "https://api.github.com/repos/$REPO/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
[ -z "$tag" ] && { echo "Could not fetch latest release"; exit 1; }

if [ "$os" = "windows" ]; then
  archive="${BIN}_${os}_${arch}.zip"
  url="https://github.com/$REPO/releases/download/$tag/$archive"
  echo "Downloading $BIN $tag ($os/$arch)..."
  curl -sfL "$url" -o "${BIN}.zip"
  unzip -qo "${BIN}.zip"
  rm -f "${BIN}.zip"
  bin_name="${BIN}.exe"
  dst="$HOME/bin/$BIN.exe"
  mkdir -p "$HOME/bin"
else
  archive="${BIN}_${os}_${arch}.tar.gz"
  url="https://github.com/$REPO/releases/download/$tag/$archive"
  echo "Downloading $BIN $tag ($os/$arch)..."
  curl -sfL "$url" | tar xz
  bin_name="$BIN"
  dst="/usr/local/bin/$BIN"
fi

if [ -w "$(dirname "$dst")" ]; then
  mv "$bin_name" "$dst"
else
  echo "  Installing to $dst (requires sudo)..."
  sudo mv "$bin_name" "$dst"
fi
echo "  ✓ Installed $BIN to $dst"
if [ "$os" = "windows" ]; then
  echo "  Ensure $HOME/bin is in your PATH to run '$BIN' from anywhere."
fi
echo "  ➜ Run '$BIN' to start"
