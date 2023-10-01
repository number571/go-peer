package http

import (
	"net/http"
)

type sLogBuilder struct {
	fService   string
	fMethod    string
	fPath      string
	fConn      string
	fMessage   string
	fLogGetter ILogGetter
}

func NewLogBuilder(pService string, pR *http.Request) ILogBuilder {
	logBuilder := &sLogBuilder{
		fService: pService,
		fMethod:  pR.Method,
		fPath:    pR.URL.Path,
		fConn:    pR.RemoteAddr,
	}
	logBuilder.fLogGetter = wrapLogBuilder(logBuilder)
	return logBuilder
}

func (p *sLogBuilder) Get() ILogGetter {
	return p.fLogGetter
}

func (p *sLogBuilder) WithMessage(pMsg string) ILogBuilder {
	p.fMessage = pMsg
	return p
}
