package web

// Embed file

import (
	"embed"
	"io/fs"
)

//go:embed static/** client/**
var embeddedFS embed.FS

var (
	StaticFS = mustSub("static")
	ClientFS = mustSub("client")
)

func mustSub(path string) fs.FS {
	sub, err := fs.Sub(embeddedFS, path)
	if err != nil {
		panic(err)
	}
	return sub
}
