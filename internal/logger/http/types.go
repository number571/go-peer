package http

type ILogBuilder interface {
	Get() ILogGetter
	WithMessage(string) ILogBuilder
}

type ILogGetter interface {
	GetService() string
	GetMethod() string
	GetPath() string
	GetConn() string
	GetMessage() string
}

const (
	CLogSuccess    = "_"
	CLogMethod     = "method"
	CLogDecodeBody = "decode_body"
	CLogNotFound   = "not_found"
	CLogRedirect   = "redirect"
)
