#!/bin/bash

# Define source and destination paths
SOURCE_DIR="templates"
DEST_DIR="static/templates"
BASE_URL="https://oss.zero-one-group.com/monorepo/templates"
INFO_FILE="$DEST_DIR/info.txt"

# Create destination directory if it doesn't exist
mkdir -p "$DEST_DIR"

# Initialize the info.txt file
echo "Last updated: $(date -u +"%Y-%m-%dT%H:%M:%SZ")" > "$INFO_FILE"
echo "" >> "$INFO_FILE"
echo "Available templates:" >> "$INFO_FILE"

# Loop through each subfolder in the source directory
for SUBFOLDER in "$SOURCE_DIR"/*; do
  if [ -d "$SUBFOLDER" ]; then
    # Get the name of the subfolder
    SUBFOLDER_NAME=$(basename "$SUBFOLDER")
    ZIP_FILE="$SUBFOLDER_NAME.zip"

    # Change to the parent directory of the subfolder
    cd "$(dirname "$SUBFOLDER")"

    # Create the ZIP file without including the parent directory
    zip -r "../$DEST_DIR/$ZIP_FILE" "$SUBFOLDER_NAME"

    # Go back to the original directory
    cd - > /dev/null

    # Append the ZIP file URL to info.txt
    echo "- $BASE_URL/$ZIP_FILE" >> "$INFO_FILE"

    # Print success message
    echo "ZIP file created at $DEST_DIR/$ZIP_FILE"
  fi
done
