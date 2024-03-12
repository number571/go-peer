package adapted

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/adapters"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ adapters.IAdaptedProducer = &sAdaptedProducer{}
)

type sAdaptedProducer struct {
	fServiceAddr string
}

func NewAdaptedProducer(pServiceAddr string) adapters.IAdaptedProducer {
	return &sAdaptedProducer{
		fServiceAddr: pServiceAddr,
	}
}

func (p *sAdaptedProducer) Produce(pCtx context.Context, pMsg net_message.IMessage) error {
	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodPost,
		fmt.Sprintf("http://%s/push", p.fServiceAddr),
		bytes.NewBuffer([]byte(pMsg.ToString())),
	)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return fmt.Errorf("got status code = %d", code)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(res) == 0 || res[0] == '!' {
		return errors.New("got invalid resp")
	}
	return nil
}
