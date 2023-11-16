package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters/common"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

const (
	workSize    = 20
	messageSize = (8 << 10)
)

const (
	databasePath = "common_recv.db"
	dataCountKey = "count_recv"
)

func initDB() database.IKVDatabase {
	var err error
	db, err := database.NewKeyValueDB(
		storage.NewSettings(&storage.SSettings{
			FPath: databasePath,
		}),
	)
	if err != nil {
		panic(err)
	}
	if _, err := db.Get([]byte(dataCountKey)); err == nil {
		return db
	}
	if err := db.Set([]byte(dataCountKey), []byte(fmt.Sprintf("%d", 0))); err != nil {
		panic(err)
	}
	return db
}

func main() {
	db := initDB()
	defer db.Close()

	if len(os.Args) != 4 {
		panic("./receiver [service-port] [hls-port] [logger]")
	}

	portService, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	portHLT, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	logger, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	transferTraffic(db, portService, portHLT, logger == 1)
}

func printLog(hasLog bool, msg error) {
	if hasLog {
		return
	}
	fmt.Println(msg)
}

func transferTraffic(db database.IKVDatabase, portService, portHLT int, hasLog bool) {
	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s:%d", "localhost", portHLT),
			&http.Client{Timeout: time.Minute},
			getMessageSettings(),
		),
	)

	for {
		time.Sleep(time.Second)

		countService, err := loadCountFromService(portService)
		if err != nil {
			printLog(hasLog, err)
			continue
		}

		countDB, err := loadCountFromDB(db)
		if err != nil {
			printLog(hasLog, err)
			continue
		}

		for i := countDB; i < countService; i++ {
			incrementCountInDB(db)

			msg, err := loadMessageFromService(portService, i)
			if err != nil {
				printLog(hasLog, err)
				continue
			}

			if err := hltClient.PutMessage(msg); err != nil {
				printLog(hasLog, err)
				continue
			}
		}
	}
}

func loadMessageFromService(portService int, id uint64) (net_message.IMessage, error) {
	// build request to service
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s:%d/load?data_id=%d", common.HostService, portService, id),
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

	msg := net_message.LoadMessage(
		getMessageSettings(),
		string(msgStringAsBytes[1:]),
	)
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	return msg, nil
}

func loadCountFromService(portService int) (uint64, error) {
	// build request to service
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s:%d/size", common.HostService, portService),
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

func loadCountFromDB(db database.IKVDatabase) (uint64, error) {
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

func incrementCountInDB(db database.IKVDatabase) error {
	count, err := loadCountFromDB(db)
	if err != nil {
		return err
	}

	if err := db.Set([]byte(dataCountKey), []byte(fmt.Sprintf("%d", count+1))); err != nil {
		return err
	}

	return nil
}

func getMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: workSize,
	})
}
