package main

import (
	"embed"
	"io/fs"
)

//go:embed ui/dist/*
var rawFs embed.FS
var uiFs, _ = fs.Sub(rawFs, "ui/dist")
