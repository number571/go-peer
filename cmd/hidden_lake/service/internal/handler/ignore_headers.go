package handler

import (
	"strings"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

var gIgnoreHeaders = initIgnoreHeaders()

func initIgnoreHeaders() map[string]struct{} {
	headers := map[string]struct{}{
		hls_settings.CHeaderResponseMode: {}, // ignore HLS header
		"Date":                           {}, // ignore due to deanonymization
		"Content-Length":                 {}, // ignore redundant header
	}

	lcHeaders := make([]string, 0, len(headers))
	for h := range headers {
		lcHeaders = append(lcHeaders, strings.ToLower(h))
	}
	for _, h := range lcHeaders {
		headers[h] = struct{}{}
	}

	return headers
}
