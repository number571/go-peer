package app

import (
	"github.com/number571/go-peer/internal/pprof"
)

func (p *sApp) initServicePPROF() {
	p.fServicePPROF = pprof.InitPprofService(p.fConfig.GetAddress().GetPPROF())
}
