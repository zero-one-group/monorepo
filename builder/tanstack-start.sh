#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
ROOT_DIR=$(dirname "$SCRIPT_DIR")
TEMPLATE_DIR="$ROOT_DIR/templates"
TEMPLATE_SOURCE_NAME="tanstack-start"
TEMPLATE_TARGET_NAME="tanstack-start"

echo "Building TanStack Start project templates..."

SOURCE_PATH="$TEMPLATE_DIR/$TEMPLATE_SOURCE_NAME"
TARGET_PATH="$TEMPLATE_DIR/$TEMPLATE_TARGET_NAME"

if [ -d "$SOURCE_PATH" ] && [ "$TEMPLATE_SOURCE_NAME" != "$TEMPLATE_TARGET_NAME" ]; then
    mv "$SOURCE_PATH" "$TARGET_PATH"
fi

replace_string() {
    local file="$1"
    local search="$2"
    local replace="$3"
    if command -v sd >/dev/null 2>&1; then
        sd "$search" "$replace" "$file"
    else
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s/$search/$replace/g" "$file"
        else
            sed -i "s/$search/$replace/g" "$file"
        fi
    fi
}

if [ -f "$TARGET_PATH/moon.yml" ]; then
    replace_string "$TARGET_PATH/moon.yml" "$TEMPLATE_SOURCE_NAME" "{{ package_name | kebab_case }}"
    replace_string "$TARGET_PATH/moon.yml" "_CHANGE_ME_DESCRIPTION_" "{{ package_description }}"
fi

# Replace string "3300" with "{{ port_number }}" in any files that contain it except "template.yml"
grep -rl "3300" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "3300" "{{ port_number }}"
done

grep -rl "$TEMPLATE_SOURCE_NAME" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "$TEMPLATE_SOURCE_NAME" "{{ package_name | kebab_case }}"
done

grep -rl "_CHANGE_ME_DESCRIPTION_" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "_CHANGE_ME_DESCRIPTION_" "{{ package_description }}"
done

if [ -f "$TARGET_PATH/package.json" ]; then
    jq 'if .dependencies then .dependencies |= del(."@repo/shared-ui") else . end' "$TARGET_PATH/package.json" > "$TARGET_PATH/package.json.tmp"
    mv "$TARGET_PATH/package.json.tmp" "$TARGET_PATH/package.json"
fi

if [ -f "$TARGET_PATH/tsconfig.json" ]; then
    jq 'if .references then .references |= map(select(.path != "../../packages/shared-ui")) else . end' "$TARGET_PATH/tsconfig.json" > "$TARGET_PATH/tsconfig.json.tmp"
    mv "$TARGET_PATH/tsconfig.json.tmp" "$TARGET_PATH/tsconfig.json"
fi
