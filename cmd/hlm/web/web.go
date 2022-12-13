package web

import (
	"embed"
	"io/fs"
)

var (
	//go:embed static
	embededStatic embed.FS

	//go:embed template
	embededTemplate embed.FS
)

func GetStaticPath() fs.FS {
	fsys, err := fs.Sub(embededStatic, "static")
	if err != nil {
		panic(err)
	}
	return fsys
}

func GetTemplatePath() fs.FS {
	fsys, err := fs.Sub(embededTemplate, "template")
	if err != nil {
		panic(err)
	}
	return fsys
}
