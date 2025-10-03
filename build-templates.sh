#!/bin/bash
echo

# Check if jq is installed
if ! command -v jq >/dev/null 2>&1; then
    echo "Error: jq is not installed. Please install jq to continue."
    exit 1
fi

# First, run code formatter
echo "Running code formatter..."
pnpm --silent run format
if [ $? -ne 0 ]; then
    echo "Code formatting failed. Please fix the issues and try again."
    exit 1
fi

TEMPLATE_DIR="./templates"
APPS_DIR="./apps"

mkdir -p $TEMPLATE_DIR

# Mapping apps dir to template dir (parallel arrays)
SRC_DIRS=("astro-web" "fastapi-ai" "go-clean" "go-modular" "nextjs-app" "react-app" "react-ssr" "strapi-cms")
TGT_DIRS=("astro"     "fastapi-ai" "go-clean" "go-modular" "nextjs"     "react-app" "react-ssr" "strapi")

echo "Copying template files with mapping and path validation..."
for i in "${!SRC_DIRS[@]}"; do
    src="${SRC_DIRS[$i]}"
    tgt="${TGT_DIRS[$i]}"
    SRC_PATH="$APPS_DIR/$src"
    TGT_PATH="$TEMPLATE_DIR/$tgt"
    if [ -d "$SRC_PATH" ]; then
        echo "Processing $src -> $tgt..."
        rm -rf "$TGT_PATH"
        mkdir -p "$TGT_PATH"
        cp -r "$SRC_PATH/." "$TGT_PATH/"
    else
        echo "Source directory $SRC_PATH does not exist, skipping."
    fi
done
echo "Template files copied successfully."
echo

echo "Cleaning up unnecessary files..."
for i in "${!TGT_DIRS[@]}"; do
    tgt="${TGT_DIRS[$i]}"
    TGT_PATH="$TEMPLATE_DIR/$tgt"
    if [ -d "$TGT_PATH" ]; then
        echo "Cleaning up $tgt..."
        rm -rf "$TGT_PATH/node_modules"
        rm -rf "$TGT_PATH/vendor"
        rm -rf "$TGT_PATH/build"
        rm -rf "$TGT_PATH/dist"
        rm -rf "$TGT_PATH/uv.lock"
        rm -rf "$TGT_PATH/.react-router"
        rm -rf "$TGT_PATH/.astro"
        rm -rf "$TGT_PATH/.cache"
        rm -rf "$TGT_PATH/.strapi"
        rm -rf "$TGT_PATH/.next"
        rm -rf "$TGT_PATH/.venv"
        rm -rf "$TGT_PATH/.env"
        find "$TGT_PATH" -type f -name ".DS_Store" -delete
    fi
done
echo "Cleanup completed."
echo

echo "Scaffolding templates..."
bash ./builder/astro.sh
bash ./builder/fastapi-ai.sh
bash ./builder/go-clean.sh
bash ./builder/go-modular.sh
bash ./builder/nextjs-app.sh
bash ./builder/react-app.sh
bash ./builder/react-ssr.sh
bash ./builder/strapi-cms.sh
echo
echo "All processes completed successfully."
echo
