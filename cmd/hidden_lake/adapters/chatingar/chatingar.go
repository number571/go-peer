package chatingar

import "net/http"

func EnrichRequest(pReq *http.Request) *http.Request {
	pReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0")
	pReq.Header.Set("Accept", "*/*")
	pReq.Header.Set("Accept-Language", "en-US,en;q=0.5")
	pReq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	pReq.Header.Set("Referer", "https://chatingar.com/")
	pReq.Header["content-type"] = []string{"application/json"}
	pReq.Header.Set("Origin", "https://chatingar.com")
	pReq.Header.Set("Connection", "keep-alive")
	pReq.Header.Set("Sec-Fetch-Dest", "empty")
	pReq.Header.Set("Sec-Fetch-Mode", "cors")
	pReq.Header.Set("Sec-Fetch-Site", "same-site")
	return pReq
}
