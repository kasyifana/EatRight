#!/bin/bash

# EatRight Backend Build Script
# This script builds the production binary

set -e

echo "🔨 Building EatRight backend..."

# Set build variables
APP_NAME="eatright-server"
BUILD_DIR="bin"
MAIN_FILE="cmd/server/main.go"

# Create build directory if it doesn't exist
mkdir -p $BUILD_DIR

# Build for Linux (production)
echo "📦 Building for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $BUILD_DIR/$APP_NAME $MAIN_FILE

# Make binary executable
chmod +x $BUILD_DIR/$APP_NAME

# Get binary size
SIZE=$(du -h $BUILD_DIR/$APP_NAME | cut -f1)

echo "✅ Build complete!"
echo "📍 Binary location: $BUILD_DIR/$APP_NAME"
echo "📏 Binary size: $SIZE"
echo ""
echo "To run the server:"
echo "  ./$BUILD_DIR/$APP_NAME"
