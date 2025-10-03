#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
ROOT_DIR=$(dirname "$SCRIPT_DIR")
APPS_DIR="$ROOT_DIR/apps"
TEMPLATE_DIR="$ROOT_DIR/templates"
TEMPLATE_SOURCE_NAME="astro-web"
TEMPLATE_TARGET_NAME="astro"

echo "Building Astro project templates..."

SOURCE_PATH="$TEMPLATE_DIR/$TEMPLATE_SOURCE_NAME"
TARGET_PATH="$TEMPLATE_DIR/$TEMPLATE_TARGET_NAME"

if [ -d "$SOURCE_PATH" ] && [ "$TEMPLATE_SOURCE_NAME" != "$TEMPLATE_TARGET_NAME" ]; then
    mv "$SOURCE_PATH" "$TARGET_PATH"
fi

# Rename all files *.astro to *.raw.astro inside the target directory
find "$TARGET_PATH" -type f -name "*.astro" | while read -r file; do
    mv "$file" "${file%.astro}.raw.astro"
done

# Inject "---\nforce: true\n---" into *.astro files
find "$TARGET_PATH" -type f -name "*.astro" | while read -r file; do
    { echo -e "---\nforce: true\n---\n"; cat "$file"; } > "${file}.tmp" && mv "${file}.tmp" "$file"
done

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

# Replace string "astro-web" with "{{ package_name | kebab_case }}" in moon.yml
if [ -f "$TARGET_PATH/moon.yml" ]; then
    replace_string "$TARGET_PATH/moon.yml" "astro-web" "{{ package_name | kebab_case }}"
fi

# Replace string "astro-web" with "{{ package_name | kebab_case }}" in Dockerfile
if [ -f "$TARGET_PATH/Dockerfile" ]; then
    replace_string "$TARGET_PATH/Dockerfile" "astro-web" "{{ package_name | kebab_case }}"
fi

# Replace string "astro-web" with "{{ package_name | kebab_case }}" in package.json
if [ -f "$TARGET_PATH/package.json" ]; then
    replace_string "$TARGET_PATH/package.json" "astro-web" "{{ package_name | kebab_case }}"
fi

# Replace string "4321" with "{{ port_number }}" in any files that contain it except "template.yml"
grep -rl "4321" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "4321" "{{ port_number }}"
done

# Replace string "_CHANGE_ME_DESCRIPTION_" with "{{ package_description }}" in moon.yml and package.json
if [ -f "$TARGET_PATH/moon.yml" ]; then
    replace_string "$TARGET_PATH/moon.yml" "_CHANGE_ME_DESCRIPTION_" "{{ package_description }}"
fi

if [ -f "$TARGET_PATH/package.json" ]; then
    replace_string "$TARGET_PATH/package.json" "_CHANGE_ME_DESCRIPTION_" "{{ package_description }}"
fi
