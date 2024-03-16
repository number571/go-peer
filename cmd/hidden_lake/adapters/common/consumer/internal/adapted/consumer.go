package adapted

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	"github.com/number571/go-peer/pkg/database"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	cDBCountKey = "db_count_key"
)

var (
	_ adapters.IAdaptedConsumer = &sAdaptedConsumer{}
)

type sAdaptedConsumer struct {
	fServiceAddr string
	fSettings    net_message.ISettings
	fKVDatabase  database.IKVDatabase
}

func NewAdaptedConsumer(
	pServiceAddr string,
	pSettings net_message.ISettings,
	pKVDatabase database.IKVDatabase,
) adapters.IAdaptedConsumer {
	return &sAdaptedConsumer{
		fServiceAddr: pServiceAddr,
		fSettings:    pSettings,
		fKVDatabase:  pKVDatabase,
	}
}

func (p *sAdaptedConsumer) Consume(pCtx context.Context) (net_message.IMessage, error) {
	time.Sleep(time.Second)

	countService, err := p.loadCountFromService(pCtx)
	if err != nil {
		return nil, err
	}

	countDB, err := p.loadCountFromDB()
	if err != nil {
		return nil, err
	}

	if countDB >= countService {
		return nil, nil
	}

	msg, err := p.loadMessageFromService(pCtx, countDB)
	if err != nil {
		return nil, err
	}

	if err := p.incrementCountInDB(); err != nil {
		return nil, err
	}

	return msg, nil
}

func (p *sAdaptedConsumer) loadMessageFromService(ctx context.Context, id uint64) (net_message.IMessage, error) {
	// build request to service
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/load?data_id=%d", p.fServiceAddr, id),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed: build request")
	}

	// send request to service
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed: bad request")
	}
	defer resp.Body.Close()

	// read response from service
	msgStringAsBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed: read body from service")
	}

	// read body of response
	if len(msgStringAsBytes) <= 1 || msgStringAsBytes[0] == '!' {
		return nil, fmt.Errorf("failed: incorrect response from service")
	}

	msg, err := net_message.LoadMessage(p.fSettings, string(msgStringAsBytes[1:]))
	if err != nil {
		return nil, fmt.Errorf("message is nil")
	}

	return msg, nil
}

func (p *sAdaptedConsumer) loadCountFromService(ctx context.Context) (uint64, error) {
	// build request to service
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/size", p.fServiceAddr),
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("failed: build request")
	}

	// send request to service
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed: bad request")
	}
	defer resp.Body.Close()

	// read response from service
	bytesCount, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed: read body from service")
	}

	// read body of response
	if len(bytesCount) <= 1 || bytesCount[0] == '!' {
		return 0, fmt.Errorf("failed: incorrect response from service")
	}

	strCount := string(bytesCount[1:])
	countService, err := strconv.ParseUint(strCount, 10, 64)
	if err != nil {
		return 0, err
	}

	return countService, nil
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

func (p *sAdaptedConsumer) incrementCountInDB() error {
	count, err := p.loadCountFromDB()
	if err != nil {
		return err
	}

	return p.fKVDatabase.Set(
		[]byte(cDBCountKey),
		[]byte(strconv.FormatUint(count+1, 10)),
	)
}
