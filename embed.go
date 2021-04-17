package main

import (
	"embed"
	"io/fs"
)

//go:embed ui/build/*
var rawFs embed.FS
var uiFs, _ = fs.Sub(rawFs, "ui/build")
