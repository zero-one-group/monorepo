package web

import "embed"

// WebDir contains SPA output files and static assets like images, fonts, etc.
// This will include the SPA output and static assets and Vite build output.
//
//go:embed static/**
var WebDir embed.FS
