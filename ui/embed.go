package ui

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var files embed.FS

func Files() fs.FS {
	fsys, _ := fs.Sub(files, "dist")
	return fsys
}
