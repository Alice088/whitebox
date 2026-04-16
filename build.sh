#!/bin/bash

APP=whitebox
DIST=dist

rm -rf $DIST
mkdir -p $DIST

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

  # playbox
  cp -r ./playbox "$OUT_DIR/"

  # .env.example (если есть)
  if [ -f ".env.example" ]; then
    cp .env.example "$OUT_DIR/.env.example"
  fi

  # архив
  if [ "$OS" = "windows" ]; then
    (cd $DIST && zip -r "$NAME.zip" "$NAME" >/dev/null)
  else
    tar -czf "$DIST/$NAME.tar.gz" -C "$DIST" "$NAME"
  fi

  rm -rf "$OUT_DIR"
}

build linux amd64 ""
build darwin arm64 ""
build windows amd64 ".exe"

echo "Done → dist/"