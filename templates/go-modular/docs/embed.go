package docs

import "embed"

// SwaggerFS contains the embedded swagger.json and swagger.yaml files.
// This will include the generated swagger documentation.
//
//go:embed swagger.json swagger.yaml
var SwaggerFS embed.FS
