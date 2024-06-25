package adapted

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	autils "github.com/number571/go-peer/cmd/hidden_lake/adapters/chatingar/internal/utils"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ adapters.IAdaptedProducer = &sAdaptedProducer{}
)

type sAdaptedProducer struct {
	fPostID string
}

// curl 'https://api.chatingar.com/api/comment' -X POST -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://chatingar.com/' -H 'content-type: application/json' -H 'Origin: https://chatingar.com' -H 'Connection: keep-alive' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site' --data-raw '{"postId":"65f7214f5b65dcbdedcca3fb","body":"\"123\""}'
func NewAdaptedProducer(pPostID string) adapters.IAdaptedProducer {
	return &sAdaptedProducer{
		fPostID: pPostID,
	}
}

func (p *sAdaptedProducer) Produce(pCtx context.Context, pMsg net_message.IMessage) error {
	reqStr := fmt.Sprintf(
		`{"postId":"%s","body":"%s"}`,
		p.fPostID,
		pMsg.ToString(),
	)

	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodPost,
		"https://chatingar-api.onrender.com/api/comment",
		bytes.NewBuffer([]byte(reqStr)),
	)
	if err != nil {
		return utils.MergeErrors(ErrBuildRequest, err)
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(autils.EnrichRequest(req))
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusCreated {
		return ErrBadStatusCode
	}

	return nil
}
