#!/bin/bash

TEMPLATES=(
  template-ansible
  template-astro
  template-fastapi-ai
  template-gitlab-cicd
  template-golang
  template-load-balancer
  template-monitoring
  template-nextjs
  template-phoenix
  template-postgresql
  template-react-app
  template-react-ssr
  template-shared-ui
  template-squidproxy
  template-strapi
  template-swarm
  template-terragrunt
)

OUTPUT_DIR="./public/templates"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

for template in "${TEMPLATES[@]}"; do
  # Skip if directory does not exist
  if [ -d "$template" ]; then
    zip -r "$OUTPUT_DIR/$template.zip" "$template"
    echo "Zipped $template -> $OUTPUT_DIR/$template.zip"
  else
    echo "Directory $template not found, skipping."
  fi
done
