package handler

import "github.com/number571/go-peer/cmd/hlm/internal/database"

type sTemplateData struct {
	Authorized bool
}

func newTemplateData(db database.IKeyValueDB) *sTemplateData {
	return &sTemplateData{
		Authorized: db != nil,
	}
}
