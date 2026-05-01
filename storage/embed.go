package storage

import "embed"

//go:embed *.json *.md
var Assets embed.FS
