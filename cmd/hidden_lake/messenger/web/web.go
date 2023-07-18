package web

import (
	"embed"
	"io/fs"
	"os"
)

const (
	// cUsedEmbedFS=false used for debug
	// not recommended use in real enviroment
	cUsedEmbedFS = true
)

const (
	cStaticPath   = "static"
	cTemplatePath = "template"
)

var (
	//go:embed static/json/emoji.json
	GEmojiJSON []byte

	//go:embed static/json/emoji_simple.json
	GEmojiSimpleJSON []byte

	//go:embed static
	gEmbededStatic embed.FS

	//go:embed template
	gEmbededTemplate embed.FS
)

func GetStaticPath() fs.FS {
	if !cUsedEmbedFS {
		return os.DirFS("./web/" + cStaticPath)
	}
	fsys, err := fs.Sub(gEmbededStatic, cStaticPath)
	if err != nil {
		panic(err)
	}
	return fsys
}

func GetTemplatePath() fs.FS {
	if !cUsedEmbedFS {
		return os.DirFS("./web/" + cTemplatePath)
	}
	fsys, err := fs.Sub(gEmbededTemplate, cTemplatePath)
	if err != nil {
		panic(err)
	}
	return fsys
}
