package adapted

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/chatingar"
	"github.com/number571/go-peer/pkg/cache"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

const (
	cPageOffet = 5
)

var (
	_ adapters.IAdaptedConsumer = &sAdaptedConsumer{}
)

type sAdaptedConsumer struct {
	fEnabled     bool
	fPostID      string
	fSettings    net_message.ISettings
	fMessages    chan net_message.IMessage
	fCacheSetter cache.ICacheSetter
	fCurrPage    uint64
}

func NewAdaptedConsumer(
	pPostID string,
	pSettings net_message.ISettings,
	pCacheSetter cache.ICacheSetter,
) adapters.IAdaptedConsumer {
	return &sAdaptedConsumer{
		fPostID:      pPostID,
		fSettings:    pSettings,
		fCacheSetter: pCacheSetter,
		fMessages:    make(chan net_message.IMessage, cPageOffet),
	}
}

func (p *sAdaptedConsumer) Consume(pCtx context.Context) (net_message.IMessage, error) {
	if !p.fEnabled {
		countComments, err := p.loadCountComments(pCtx)
		if err != nil {
			return nil, err
		}

		p.fCurrPage = (countComments / cPageOffet) + 1
		p.fEnabled = true
	}

	return p.loadMessage(pCtx)
}

// curl 'https://api.chatingar.com/api/comment/65f7214f5b65dcbdedcca3fb?page=1' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://chatingar.com/' -H 'Origin: https://chatingar.com' -H 'Connection: keep-alive' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site'
func (p *sAdaptedConsumer) loadMessage(pCtx context.Context) (net_message.IMessage, error) {
	select {
	case msg := <-p.fMessages:
		return msg, nil
	default:
		// do request
	}

	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodGet,
		fmt.Sprintf(
			"https://api.chatingar.com/api/comment/%s?page=%d",
			p.fPostID,
			p.fCurrPage,
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed: build request")
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(chatingar.EnrichRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed: bad request")
	}
	defer resp.Body.Close()

	var messagesDTO sMessagesDTO
	if err := json.NewDecoder(resp.Body).Decode(&messagesDTO); err != nil {
		return nil, err
	}

	sizeComments := len(messagesDTO.Comments)
	if sizeComments > cPageOffet {
		return nil, errors.New("has limit pages")
	}
	if sizeComments == cPageOffet {
		p.fCurrPage++
	}

	for _, v := range messagesDTO.Comments {
		msg, err := net_message.LoadMessage(p.fSettings, v.Body)
		if err != nil {
			continue
		}
		if ok := p.rememberMessage(msg); !ok {
			continue
		}
		p.fMessages <- msg
	}

	select {
	case msg := <-p.fMessages:
		return msg, nil
	default:
		return nil, nil
	}
}

// curl 'https://api.chatingar.com/api/post/65f7214f5b65dcbdedcca3fb' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://chatingar.com/' -H 'Origin: https://chatingar.com' -H 'Connection: keep-alive' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site'
func (p *sAdaptedConsumer) loadCountComments(pCtx context.Context) (uint64, error) {
	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodGet,
		fmt.Sprintf("https://api.chatingar.com/api/post/%s", p.fPostID),
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("failed: build request")
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(chatingar.EnrichRequest(req))
	if err != nil {
		return 0, fmt.Errorf("failed: bad request")
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return 0, fmt.Errorf("got status code = %d", code)
	}

	var count sCountDTO
	if err := json.NewDecoder(resp.Body).Decode(&count); err != nil {
		return 0, err
	}

	result := count.Post.CommentCount
	if result < 0 {
		return 0, errors.New("got count < 0")
	}

	return uint64(result), nil
}

func (p *sAdaptedConsumer) rememberMessage(pMsg net_message.IMessage) bool {
	hash := hashing.NewSHA256Hasher(pMsg.GetHash()).ToBytes()
	return p.fCacheSetter.Set(hash, []byte{})
}
