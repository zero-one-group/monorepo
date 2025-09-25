#!/bin/bash

# Define source and destination paths
SOURCE_DIR="templates"
DEST_DIR="docsite/static/templates"
BASE_URL="https://oss.zero-one-group.com/monorepo/templates"
METADATA_FILE="docsite/static/templates.json"

# Create destination directory if it doesn't exist
mkdir -p "$DEST_DIR"

# Initialize the templates.json file
echo "{" > "$METADATA_FILE"
echo "  \"last_updated\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"," >> "$METADATA_FILE"
echo "  \"templates\": [" >> "$METADATA_FILE"

# Loop through each subfolder in the source directory
FIRST=true
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

    # Append the ZIP file metadata to templates.json
    if [ "$FIRST" = true ]; then
      FIRST=false
    else
      echo "    ," >> "$METADATA_FILE"
    fi
    echo "    {" >> "$METADATA_FILE"
    echo "      \"name\": \"$SUBFOLDER_NAME\"," >> "$METADATA_FILE"
    echo "      \"url\": \"$BASE_URL/$ZIP_FILE\"" >> "$METADATA_FILE"
    echo "    }" >> "$METADATA_FILE"

    # Print success message
    echo "ZIP file created at $DEST_DIR/$ZIP_FILE"
  fi
done

# Close the templates.json file
echo "  ]" >> "$METADATA_FILE"
echo "}" >> "$METADATA_FILE"
