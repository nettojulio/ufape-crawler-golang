#!/bin/bash

set -e

# === Config ===
VERSION=$(cat VERSION)
APP_NAME="ufape-crawler-golang"
BIN_NAME=$APP_NAME
OUTPUT_BIN="dist/${BIN_NAME}-linux-static"

echo "🔧 Versão: $VERSION"
echo "🛠️  Nome do binário: $BIN_NAME"
echo "🐳 Nome da imagem: $IMAGE_TAG"

# === Build Go binário ===
echo "📦 Compilando Go..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -ldflags="-s -w -X 'main.Version=$VERSION'" \
-a -installsuffix cgo \
-o $OUTPUT_BIN cmd/main.go
