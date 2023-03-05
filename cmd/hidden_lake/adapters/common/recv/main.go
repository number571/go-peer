package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters/common"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage/database"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

const (
	databasePath = "recv.db"
	dataCountKey = "count_recv"
)

var (
	db database.IKeyValueDB
)

func init() {
	db = database.NewLevelDB(
		database.NewSettings(&database.SSettings{
			FPath: databasePath,
		}),
	)
	if _, err := db.Get([]byte(dataCountKey)); err == nil {
		return
	}
	if err := db.Set([]byte(dataCountKey), []byte(fmt.Sprintf("%d", 0))); err != nil {
		panic(err)
	}
}

func main() {
	defer db.Close()

	if len(os.Args) != 3 {
		panic("./receiver [service-port] [hlt-port]")
	}

	portService, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	portHLT, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s:%d", "localhost", portHLT),
			message.NewParams(hls_settings.CMessageSize, hls_settings.CWorkSize),
		),
	)

	for {
		time.Sleep(time.Second)

		strCount, err := api.Request(
			http.MethodGet,
			fmt.Sprintf("%s:%d/size", common.HostService, portService),
			nil,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}

		countService, err := strconv.ParseUint(strCount, 10, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		countDB, err := countOfDataInDB()
		if err != nil {
			fmt.Println(err)
			continue
		}

		for i := countDB; i < countService; i++ {
			data, err := api.Request(
				http.MethodGet,
				fmt.Sprintf("%s:%d/load?data_id=%d", common.HostService, portService, i),
				nil,
			)
			if err != nil {
				fmt.Println(err)
				continue
			}

			msg := message.LoadMessage(
				encoding.HexDecode(data),
				message.NewParams(
					hls_settings.CMessageSize,
					hls_settings.CWorkSize,
				),
			)
			if err := hltClient.PutMessage(msg); err != nil {
				fmt.Println(err)
				continue
			}

			incrementCountInDB()
		}
	}
}

func countOfDataInDB() (uint64, error) {
	res, err := db.Get([]byte(dataCountKey))
	if err != nil {
		return 0, err
	}

	count, err := strconv.ParseUint(string(res), 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func incrementCountInDB() error {
	count, err := countOfDataInDB()
	if err != nil {
		return err
	}

	if err := db.Set([]byte(dataCountKey), []byte(fmt.Sprintf("%d", count+1))); err != nil {
		return err
	}

	return nil
}
