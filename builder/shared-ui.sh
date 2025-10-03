#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
ROOT_DIR=$(dirname "$SCRIPT_DIR")
PKG_DIR="$ROOT_DIR/packages"
TEMPLATE_DIR="$ROOT_DIR/templates"
TEMPLATE_SOURCE_NAME="shared-ui"
TEMPLATE_TARGET_NAME="shared-ui"

echo "Building Shared UI project templates..."

SRC_PATH="$PKG_DIR/$TEMPLATE_SOURCE_NAME"
TARGET_PATH="$TEMPLATE_DIR/$TEMPLATE_TARGET_NAME"

# Remove target if exists, then copy source to target
if [ -d "$SRC_PATH" ]; then
    rm -rf "$TARGET_PATH"
    cp -R "$SRC_PATH" "$TARGET_PATH"
else
    echo "Source directory $SRC_PATH does not exist, skipping copy."
fi

# Cleanup unnecessary files in target
if [ -d "$TARGET_PATH" ]; then
    rm -rf "$TARGET_PATH/node_modules"
    rm -rf "$TARGET_PATH/storybook-static"
    rm -rf "$TARGET_PATH/dist"
    rm -rf "$TARGET_PATH/build"
    rm -rf "$TARGET_PATH/.DS_Store"
    find "$TARGET_PATH" -type f -name ".DS_Store" -delete
fi

SOURCE_PATH="$TEMPLATE_DIR/$TEMPLATE_SOURCE_NAME"
TARGET_PATH="$TEMPLATE_DIR/$TEMPLATE_TARGET_NAME"

if [ -d "$SOURCE_PATH" ] && [ "$TEMPLATE_SOURCE_NAME" != "$TEMPLATE_TARGET_NAME" ]; then
    mv "$SOURCE_PATH" "$TARGET_PATH"
fi

# Function to replace string using sd or sed
replace_string() {
    local file="$1"
    local search="$2"
    local replace="$3"
    if command -v sd >/dev/null 2>&1; then
        sd "$search" "$replace" "$file"
    else
        # macOS uses 'sed -i ''', Linux uses 'sed -i'
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s/$search/$replace/g" "$file"
        else
            sed -i "s/$search/$replace/g" "$file"
        fi
    fi
}

if [ -f "$TARGET_PATH/moon.yml" ]; then
    # Replace string "$TEMPLATE_SOURCE_NAME" with "{{ package_name | kebab_case }}" in moon.yml
    replace_string "$TARGET_PATH/moon.yml" "$TEMPLATE_SOURCE_NAME" "{{ package_name | kebab_case }}"
fi

# Replace string "$TEMPLATE_SOURCE_NAME" with "{{ package_name | kebab_case }}" in any files that contain it except "template.yml"
grep -rl "$TEMPLATE_SOURCE_NAME" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "$TEMPLATE_SOURCE_NAME" "{{ package_name | kebab_case }}"
done

# Rename all files *.tsx to *.raw.tsx inside the target directory
find "$TARGET_PATH" -type f -name "*.tsx" | while read -r file; do
    mv "$file" "${file%.tsx}.raw.tsx"
done

# Rename all files *.mdx to *.raw.mdx inside the target directory
find "$TARGET_PATH" -type f -name "*.mdx" | while read -r file; do
    mv "$file" "${file%.mdx}.raw.mdx"
done
