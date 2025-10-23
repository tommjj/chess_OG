package web

// Embed file

import (
	"embed"
	"io/fs"
)

//go:embed static/**
var embeddedFS embed.FS

var (
	StaticFS = mustSub("static")
)

func mustSub(path string) fs.FS {
	sub, err := fs.Sub(embeddedFS, path)
	if err != nil {
		panic(err)
	}
	return sub
}
