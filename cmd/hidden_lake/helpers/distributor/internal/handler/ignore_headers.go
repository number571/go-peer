package handler

import (
	"strings"
)

var gIgnoreHeaders = initIgnoreHeaders()

func initIgnoreHeaders() map[string]struct{} {
	headers := map[string]struct{}{
		"Date":           {}, // ignore due to deanonymization
		"Content-Length": {}, // ignore redundant header
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
