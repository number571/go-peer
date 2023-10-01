package http

type sLogGetter struct {
	fLogBuilder *sLogBuilder
}

func wrapLogBuilder(pLogBuilder ILogBuilder) ILogGetter {
	return &sLogGetter{
		fLogBuilder: pLogBuilder.(*sLogBuilder),
	}
}

func (p *sLogGetter) GetService() string {
	return p.fLogBuilder.fService
}

func (p *sLogGetter) GetConn() string {
	return p.fLogBuilder.fConn
}

func (p *sLogGetter) GetMethod() string {
	return p.fLogBuilder.fMethod
}

func (p *sLogGetter) GetPath() string {
	return p.fLogBuilder.fPath
}

func (p *sLogGetter) GetMessage() string {
	return p.fLogBuilder.fMessage
}
