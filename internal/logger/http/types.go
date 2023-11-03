package http

type ILogBuilder interface {
	ILogGetterFactory

	WithMessage(string) ILogBuilder
}

type ILogGetterFactory interface {
	Get() ILogGetter
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
