package adapted

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/chatingar"
	"github.com/number571/go-peer/pkg/database"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	cPage       = 5
	cDBCountKey = "db_count_key"
)

var (
	_ adapters.IAdaptedConsumer = &sAdaptedConsumer{}
)

type sAdaptedConsumer struct {
	fEnabled    bool
	fPostID     string
	fSettings   net_message.ISettings
	fKVDatabase database.IKVDatabase
	fMessages   chan net_message.IMessage
}

func NewAdaptedConsumer(
	pPostID string,
	pSettings net_message.ISettings,
	pKVDatabase database.IKVDatabase,
) adapters.IAdaptedConsumer {
	return &sAdaptedConsumer{
		fPostID:     pPostID,
		fSettings:   pSettings,
		fKVDatabase: pKVDatabase,
		fMessages:   make(chan net_message.IMessage, cPage),
	}
}

func (p *sAdaptedConsumer) Consume(pCtx context.Context) (net_message.IMessage, error) {
	select {
	case msg := <-p.fMessages:
		return msg, nil
	default:
	}

	countService, err := p.loadCountFromService(pCtx)
	if err != nil {
		return nil, err
	}

	// start read from last
	if !p.fEnabled {
		if err := p.setCountInDB(countService); err != nil {
			return nil, err
		}
		p.fEnabled = true
		return nil, nil
	}

	countDB, err := p.loadCountFromDB()
	if err != nil {
		return nil, err
	}

	if countDB >= countService {
		return nil, nil
	}

	return p.loadMessageFromService(pCtx, countDB)
}

// curl 'https://api.chatingar.com/api/comment/65f7214f5b65dcbdedcca3fb?page=1' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://chatingar.com/' -H 'Origin: https://chatingar.com' -H 'Connection: keep-alive' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site' -H '
func (p *sAdaptedConsumer) loadMessageFromService(ctx context.Context, id uint64) (net_message.IMessage, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"https://api.chatingar.com/api/comment/%s?page=%d",
			p.fPostID,
			(id/cPage)+1,
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed: build request")
	}

	resp, err := http.DefaultClient.Do(chatingar.EnrichRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed: bad request")
	}
	defer resp.Body.Close()

	var messagesDTO sMessagesDTO
	if err := json.NewDecoder(resp.Body).Decode(&messagesDTO); err != nil {
		return nil, err
	}

	count := uint64(0)
	offset := (id % cPage)

	for i := offset; i < uint64(len(messagesDTO.Comments)); i++ {
		count++
		v := messagesDTO.Comments[i]
		msg, err := net_message.LoadMessage(p.fSettings, v.Body)
		if err != nil {
			continue
		}
		p.fMessages <- msg
	}

	if err := p.addCountInDB(count); err != nil {
		return nil, err
	}

	return nil, nil
}

// curl 'https://api.chatingar.com/api/post/65f7214f5b65dcbdedcca3fb' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://chatingar.com/' -H 'Origin: https://chatingar.com' -H 'Connection: keep-alive' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site'
func (p *sAdaptedConsumer) loadCountFromService(ctx context.Context) (uint64, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.chatingar.com/api/post/%s", p.fPostID),
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("failed: build request")
	}

	resp, err := http.DefaultClient.Do(chatingar.EnrichRequest(req))
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

func (p *sAdaptedConsumer) loadCountFromDB() (uint64, error) {
	res, err := p.fKVDatabase.Get([]byte(cDBCountKey))
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			return 0, err
		}
		res = []byte(strconv.Itoa(0))
		if err := p.fKVDatabase.Set([]byte(cDBCountKey), res); err != nil {
			return 0, err
		}
	}
	return strconv.ParseUint(string(res), 10, 64)
}

func (p *sAdaptedConsumer) addCountInDB(n uint64) error {
	count, err := p.loadCountFromDB()
	if err != nil {
		return err
	}
	return p.setCountInDB(count + n)
}

func (p *sAdaptedConsumer) setCountInDB(n uint64) error {
	return p.fKVDatabase.Set(
		[]byte(cDBCountKey),
		[]byte(strconv.FormatUint(n, 10)),
	)
}
