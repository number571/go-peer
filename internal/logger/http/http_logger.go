package logbuilder

import (
	"fmt"
	"net/http"
)

const (
	cLogTemplate = "service=%s method=%s path=%s conn=%s message=%s"
)

type sLogger struct {
	fService string
	fMethod  string
	fPath    string
	fConn    string
}

func NewHTTPLogger(pService string, pR *http.Request) IHTTPLogger {
	return &sLogger{
		fService: pService,
		fMethod:  pR.Method,
		fPath:    pR.URL.Path,
		fConn:    pR.RemoteAddr,
	}
}

func (p *sLogger) Get(pMessage string) string {
	return fmt.Sprintf(cLogTemplate, p.fService, p.fMethod, p.fPath, p.fConn, pMessage)
}
