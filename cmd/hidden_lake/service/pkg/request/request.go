package request

import (
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IRequest = &sRequest{}
)

type sRequest struct {
	fMutex  sync.Mutex
	FMethod string            `json:"method"`
	FHost   string            `json:"host"`
	FPath   string            `json:"path"`
	FHead   map[string]string `json:"head"`
	FBody   []byte            `json:"body"`
}

func NewRequest(method, host, path string) IRequest {
	return &sRequest{
		FMethod: method,
		FHost:   host,
		FPath:   path,
	}
}

func LoadRequest(data []byte) IRequest {
	var request = new(sRequest)
	encoding.Deserialize(data, request)
	return request
}

func (r *sRequest) Bytes() []byte {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()

	return encoding.Serialize(r)
}

func (r *sRequest) WithHead(head map[string]string) IRequest {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()

	r.FHead = make(map[string]string)
	for k, v := range head {
		r.FHead[k] = v
	}
	return r
}

func (r *sRequest) WithBody(body []byte) IRequest {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()

	r.FBody = body
	return r
}

func (r *sRequest) Host() string {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()

	return r.FHost
}

func (r *sRequest) Path() string {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()

	return r.FPath
}

func (r *sRequest) Method() string {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()

	return r.FMethod
}

func (r *sRequest) Head() map[string]string {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()

	headers := make(map[string]string)
	for k, v := range r.FHead {
		headers[k] = v
	}

	return headers
}

func (r *sRequest) Body() []byte {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()

	return r.FBody
}
