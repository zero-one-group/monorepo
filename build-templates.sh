#!/bin/bash

TEMPLATES=(
  ansible
  astro
  fastapi-ai
  gitlab-cicd
  golang
  load-balancer
  monitoring
  nextjs
  phoenix
  postgresql
  react-app
  react-ssr
  shared-ui
  squidproxy
  strapi
  swarm
  terragrunt
)

TEMPLATE_DIR="./templates"
OUTPUT_DIR="./docs/templates"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

for template in "${TEMPLATES[@]}"; do
  template_path="$TEMPLATE_DIR/$template"
  # Skip if directory does not exist
  if [ -d "$template_path" ]; then
    zip -r "$OUTPUT_DIR/$template.zip" "$template_path"
    echo "Zipped $template_path -> $OUTPUT_DIR/$template.zip"
  else
    echo "Directory $template_path not found, skipping."
  fi
done
