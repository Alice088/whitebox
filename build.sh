#!/bin/bash

set -e

APP=whitebox
DIST=dist
INSTALL_DIR="$HOME/.local/bin"

rm -rf "$DIST"
mkdir -p "$DIST"

echo "Building..."

build() {
  OS=$1
  ARCH=$2
  EXT=$3

  NAME="$APP-$OS"
  OUT_DIR="$DIST/$NAME"

  mkdir -p "$OUT_DIR"

  echo "→ $OS/$ARCH"

  GOOS=$OS GOARCH=$ARCH go build -o "$OUT_DIR/$APP$EXT" ./cmd/whitebox/main.go

  cp -r ./playbox "$OUT_DIR/" 2>/dev/null || true

  if [ -f ".env.example" ]; then
    cp .env.example "$OUT_DIR/.env.example"
  fi

  if [ "$OS" = "windows" ]; then
    (cd "$DIST" && zip -r "$NAME.zip" "$NAME" >/dev/null)
  else
    tar -czf "$DIST/$NAME.tar.gz" -C "$DIST" "$NAME"
  fi

  rm -rf "$OUT_DIR"
}

install_local() {
  echo "Installing local binary..."

  mkdir -p "$INSTALL_DIR"

  go build -o "$INSTALL_DIR/$APP" ./cmd/whitebox/main.go
  chmod +x "$INSTALL_DIR/$APP"

  SHELL_RC=""

  if [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="$HOME/.zshrc"
  else
    SHELL_RC="$HOME/.bashrc"
  fi

  if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
    echo "" >> "$SHELL_RC"
    echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$SHELL_RC"
    export PATH="$HOME/.local/bin:$PATH"
    echo "Added $INSTALL_DIR to PATH in $SHELL_RC"
  fi

  echo "Installed: $INSTALL_DIR/$APP"
  echo "Run:"
  echo "  whitebox"
}

build linux amd64 ""
build darwin arm64 ""
build windows amd64 ".exe"

install_local

echo "Done → dist/"