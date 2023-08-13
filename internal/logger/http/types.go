package logbuilder

type IHTTPLogger interface {
	Get(string) string
}

const (
	CLogSuccess    = "_"
	CLogMethod     = "method"
	CLogDecodeBody = "decode_body"
	CLogNotFound   = "not_found"
	CLogRedirect   = "redirect"
)
