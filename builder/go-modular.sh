#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
ROOT_DIR=$(dirname "$SCRIPT_DIR")
APPS_DIR="$ROOT_DIR/apps"
TEMPLATE_DIR="$ROOT_DIR/templates"
TEMPLATE_SOURCE_NAME="go-modular"
TEMPLATE_TARGET_NAME="go-modular"

echo "Building Go Modular project templates..."

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

# Replace string "8000" with "{{ port_number }}" in any files that contain it except "template.yml"
grep -rl "8000" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "8000" "{{ port_number }}"
done

# Replace string "$TEMPLATE_SOURCE_NAME" with "{{ package_name | kebab_case }}" in any files that contain it except "template.yml"
grep -rl "$TEMPLATE_SOURCE_NAME" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "$TEMPLATE_SOURCE_NAME" "{{ package_name | kebab_case }}"
done

# Replace string "_CHANGE_ME_DESCRIPTION_" with "{{ package_description }}" in any files that contain it except "template.yml"
grep -rl "_CHANGE_ME_DESCRIPTION_" "$TARGET_PATH" | grep -v "template.yml" | while read -r file; do
    replace_string "$file" "_CHANGE_ME_DESCRIPTION_" "{{ package_description }}"
done

# Move $TARGET_PATH/.mockery.yml to $TARGET_PATH/.mockery.raw.yml
if [ -f "$TARGET_PATH/.mockery.yml" ]; then
    mv "$TARGET_PATH/.mockery.yml" "$TARGET_PATH/.mockery.raw.yml"
fi

# Move go-modular/internal/notification/mailer_test.go to go-modular/internal/notification/mailer_test.raw.go
if [ -f "$TARGET_PATH/internal/notification/mailer_test.go" ]; then
    mv "$TARGET_PATH/internal/notification/mailer_test.go" "$TARGET_PATH/internal/notification/mailer_test.raw.go"
fi

# Move go-modular/templates/emails/*.html to go-modular/templates/emails/*.raw.html
if [ -d "$TARGET_PATH/templates/emails" ]; then
    find "$TARGET_PATH/templates/emails" -type f -name "*.html" | while read -r file; do
        mv "$file" "${file%.html}.raw.html"
    done
fi

# Remove everything inside "templates/go-modular/docs" except "embed.go"
if [ -d "$TARGET_PATH/docs" ]; then
    find "$TARGET_PATH/docs" -type f ! -name 'embed.go' -delete
    find "$TARGET_PATH/docs" -type d -empty -delete
fi

# Remove everything inside "templates/go-modular/web" except "embed.go" and "static/index.html"
if [ -d "$TARGET_PATH/web" ]; then
    find "$TARGET_PATH/web" -type f ! -name 'embed.go' ! -path "$TARGET_PATH/web/static/index.html" -delete
    find "$TARGET_PATH/web" -type d -empty -delete
fi
