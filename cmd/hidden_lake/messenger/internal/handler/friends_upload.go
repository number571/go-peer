package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
)

type sUploadFile struct {
	*state.STemplateState
	FAliasName    string
	FMessageLimit uint64
}

func FriendsUploadPage(pStateManager state.IStateManager, msgLimit uint64) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/friends/upload" {
			NotFoundPage(pStateManager)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		aliasName := pR.URL.Query().Get("alias_name")
		if aliasName == "" {
			fmt.Fprint(pW, "alias name is null")
			return
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"upload.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		res := &sUploadFile{
			STemplateState: pStateManager.GetTemplate(),
			FAliasName:     aliasName,
			FMessageLimit:  msgLimit,
		}
		t.Execute(pW, res)
	}
}
