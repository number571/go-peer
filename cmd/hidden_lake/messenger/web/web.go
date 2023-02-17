package web

import (
	"embed"
	"io/fs"
	"os"
)

const (
	cUsedEmbedFS = true
)

const (
	staticPath   = "static"
	templatePath = "template"
)

var (
	//go:embed static
	gEmbededStatic embed.FS

	//go:embed template
	gEmbededTemplate embed.FS
)

func GetStaticPath() fs.FS {
	if !cUsedEmbedFS {
		return os.DirFS("./web/" + staticPath)
	}
	fsys, err := fs.Sub(gEmbededStatic, staticPath)
	if err != nil {
		panic(err)
	}
	return fsys
}

func GetTemplatePath() fs.FS {
	if !cUsedEmbedFS {
		return os.DirFS("./web/" + templatePath)
	}
	fsys, err := fs.Sub(gEmbededTemplate, templatePath)
	if err != nil {
		panic(err)
	}
	return fsys
}
