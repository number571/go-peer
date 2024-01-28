package app

import (
	"fmt"
	"os"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
)

func (p *sApp) initStorage() {
	os.MkdirAll(fmt.Sprintf("%s/%s", p.fPathTo, hlf_settings.CPathLoadedSTG), 0o777)
}
