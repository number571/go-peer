package adapted

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	"github.com/number571/go-peer/pkg/database"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/utils"
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
	countService, err := p.loadCountFromService(pCtx)
	if err != nil {
		return nil, utils.MergeErrors(ErrLoadCountService, err)
	}

	countDB, err := p.loadCountFromDB()
	if err != nil {
		return nil, utils.MergeErrors(ErrLoadCountDB, err)
	}

	if countDB >= countService {
		return nil, nil
	}

	msg, err := p.loadMessageFromService(pCtx, countDB)
	if err != nil {
		return nil, utils.MergeErrors(ErrLoadMessage, err)
	}

	if err := p.incrementCountInDB(); err != nil {
		return nil, utils.MergeErrors(ErrIncrementCount, err)
	}

	return msg, nil
}

func (p *sAdaptedConsumer) loadMessageFromService(pCtx context.Context, pID uint64) (net_message.IMessage, error) {
	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodGet,
		fmt.Sprintf("http://%s/load?data_id=%d", p.fServiceAddr, pID),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBuildRequest, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}
	defer resp.Body.Close()

	msgStringAsBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.MergeErrors(ErrReadResponse, err)
	}

	if len(msgStringAsBytes) <= 1 || msgStringAsBytes[0] == '!' {
		return nil, ErrInvalidResponse
	}

	msg, err := net_message.LoadMessage(p.fSettings, string(msgStringAsBytes[1:]))
	if err != nil {
		return nil, utils.MergeErrors(ErrDecodeMessage, err)
	}

	return msg, nil
}

func (p *sAdaptedConsumer) loadCountFromService(pCtx context.Context) (uint64, error) {
	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodGet,
		fmt.Sprintf("http://%s/size", p.fServiceAddr),
		nil,
	)
	if err != nil {
		return 0, utils.MergeErrors(ErrBuildRequest, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, utils.MergeErrors(ErrBadRequest, err)
	}
	defer resp.Body.Close()

	bytesCount, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, utils.MergeErrors(ErrReadResponse, err)
	}

	if len(bytesCount) <= 1 || bytesCount[0] == '!' {
		return 0, ErrInvalidResponse
	}

	strCount := string(bytesCount[1:])
	countService, err := strconv.ParseUint(strCount, 10, 64)
	if err != nil {
		return 0, utils.MergeErrors(ErrParseCount, err)
	}

	return countService, nil
}

func (p *sAdaptedConsumer) loadCountFromDB() (uint64, error) {
	res, err := p.fKVDatabase.Get([]byte(cDBCountKey))
	if err != nil {
		if !errors.Is(err, database.ErrGetNotFound) {
			return 0, utils.MergeErrors(ErrGetCount, err)
		}
		res = []byte(strconv.Itoa(0))
		if err := p.fKVDatabase.Set([]byte(cDBCountKey), res); err != nil {
			return 0, utils.MergeErrors(ErrInitCountKey, err)
		}
	}

	count, err := strconv.ParseUint(string(res), 10, 64)
	if err != nil {
		return 0, utils.MergeErrors(ErrParseCount, err)
	}
	return count, nil
}

func (p *sAdaptedConsumer) incrementCountInDB() error {
	count, err := p.loadCountFromDB()
	if err != nil {
		return utils.MergeErrors(ErrLoadCountDB, err)
	}

	err = p.fKVDatabase.Set(
		[]byte(cDBCountKey),
		[]byte(strconv.FormatUint(count+1, 10)),
	)
	if err != nil {
		return utils.MergeErrors(ErrSetNewCount, err)
	}
	return nil
}
