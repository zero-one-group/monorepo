#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
ROOT_DIR=$(dirname "$SCRIPT_DIR")
TEMPLATE_DIR="$ROOT_DIR/templates"
TEMPLATE_TARGET_NAME="vue-app"
TARGET_PATH="$TEMPLATE_DIR/$TEMPLATE_TARGET_NAME"

echo "Building Vue SPA project templates..."

if [ -d "$TARGET_PATH" ]; then
    rm -rf "$TARGET_PATH/node_modules"
    rm -rf "$TARGET_PATH/dist"
    rm -rf "$TARGET_PATH/tests-results"
    find "$TARGET_PATH" -type f -name ".DS_Store" -delete
fi
