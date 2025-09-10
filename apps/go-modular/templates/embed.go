package web

import "embed"

// TemplateDir embeds all files matching the *html pattern into the binary.
// This is not limited to email HTML templates — any file placed under the
// emails/ directory and matching the embed pattern will be included at build
// time and available at runtime via the embed.FS.
//
//go:embed emails/*.html
var TemplateDir embed.FS
