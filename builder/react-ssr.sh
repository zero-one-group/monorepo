#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
ROOT_DIR=$(dirname "$SCRIPT_DIR")
APPS_DIR="$ROOT_DIR/apps"
TEMPLATE_DIR="$ROOT_DIR/templates"
TEMPLATE_SOURCE_NAME="react-ssr"
TEMPLATE_TARGET_NAME="react-ssr"

echo "Building React SSR project templates..."

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
    # Replace string "_CHANGE_ME_DESCRIPTION_" with "{{ package_description }}" in moon.yml
    replace_string "$TARGET_PATH/moon.yml" "_CHANGE_ME_DESCRIPTION_" "{{ package_description }}"
fi

# Replace string "3100" with "{{ port_number }}" in any files that contain it except "template.yml"
grep -rl "3100" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "3100" "{{ port_number }}"
done

# Replace string "$TEMPLATE_SOURCE_NAME" with "{{ package_name | kebab_case }}" in any files that contain it except "template.yml"
grep -rl "$TEMPLATE_SOURCE_NAME" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "$TEMPLATE_SOURCE_NAME" "{{ package_name | kebab_case }}"
done

# Replace string "_CHANGE_ME_DESCRIPTION_" with "{{ package_description }}" in any files that contain it except "template.yml"
grep -rl "_CHANGE_ME_DESCRIPTION_" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "_CHANGE_ME_DESCRIPTION_" "{{ package_description }}"
done

# Remove "@repo/shared-ui" dependency from package.json using jq
if [ -f "$TARGET_PATH/package.json" ]; then
    jq 'if .dependencies then .dependencies |= del(.["@repo/shared-ui"]) else . end' "$TARGET_PATH/package.json" > "$TARGET_PATH/package.json.tmp"
    mv "$TARGET_PATH/package.json.tmp" "$TARGET_PATH/package.json"
fi

# Remove reference to "shared-ui" in tsconfig.json using jq
if [ -f "$TARGET_PATH/tsconfig.json" ]; then
    jq 'if .references then .references |= map(select(.path != "../../packages/shared-ui")) else . end' "$TARGET_PATH/tsconfig.json" > "$TARGET_PATH/tsconfig.json.tmp"
    mv "$TARGET_PATH/tsconfig.json.tmp" "$TARGET_PATH/tsconfig.json"
fi
