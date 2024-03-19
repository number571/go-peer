package adapted

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/chatingar"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/database"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	cPageOffet  = 5
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
		fMessages:   make(chan net_message.IMessage, cPageOffet),
	}
}

func (p *sAdaptedConsumer) Consume(pCtx context.Context) (net_message.IMessage, error) {
	select {
	case msg := <-p.fMessages:
		return msg, nil
	default:
	}

	countComments, err := p.loadCountComments(pCtx)
	if err != nil {
		return nil, err
	}

	// start read from last
	countPages := (countComments / cPageOffet) + 1
	if !p.fEnabled {
		if err := p.setCountPagesDB(countPages); err != nil {
			return nil, err
		}
		p.fEnabled = true
		return nil, nil
	}

	countPagesDB, err := p.getCountPagesDB()
	if err != nil {
		return nil, err
	}

	fmt.Println(countPagesDB, countPages)
	if countPagesDB >= countPages {
		return nil, nil
	}

	return p.loadMessage(pCtx, countPagesDB)
}

// curl 'https://api.chatingar.com/api/comment/65f7214f5b65dcbdedcca3fb?page=1' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://chatingar.com/' -H 'Origin: https://chatingar.com' -H 'Connection: keep-alive' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site'
func (p *sAdaptedConsumer) loadMessage(ctx context.Context, page uint64) (net_message.IMessage, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"https://api.chatingar.com/api/comment/%s?page=%d",
			p.fPostID,
			page,
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

	for _, v := range messagesDTO.Comments {
		if err := p.setCommentHash(v); err != nil {
			continue
		}
		msg, err := net_message.LoadMessage(p.fSettings, v.Body)
		if err != nil {
			continue
		}
		p.fMessages <- msg
	}

	return nil, p.incCountPagesDB()
}

// curl 'https://api.chatingar.com/api/post/65f7214f5b65dcbdedcca3fb' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://chatingar.com/' -H 'Origin: https://chatingar.com' -H 'Connection: keep-alive' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site'
func (p *sAdaptedConsumer) loadCountComments(ctx context.Context) (uint64, error) {
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

func (p *sAdaptedConsumer) getCountPagesDB() (uint64, error) {
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

func (p *sAdaptedConsumer) incCountPagesDB() error {
	count, err := p.getCountPagesDB()
	if err != nil {
		return err
	}
	return p.setCountPagesDB(count + 1)
}

func (p *sAdaptedConsumer) setCountPagesDB(n uint64) error {
	return p.fKVDatabase.Set(
		[]byte(cDBCountKey),
		[]byte(strconv.FormatUint(n, 10)),
	)
}

func (p *sAdaptedConsumer) setCommentHash(b sCommentBlock) error {
	hash := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{[]byte(b.Body), []byte(b.Timestamp)},
		[]byte{},
	)).ToBytes()
	if _, err := p.fKVDatabase.Get(hash); err == nil {
		return errors.New("hash already exist")
	}
	return p.fKVDatabase.Set(hash, []byte{})
}
