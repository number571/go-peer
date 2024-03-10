package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
)

const (
	workSize    = 22
	messageSize = (8 << 10)
)

const (
	databasePath = "common_consumer.db"
	dataCountKey = "count_consumer"
)

func initDB() database.IKVDatabase {
	var err error
	db, err := database.NewKVDatabase(
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
	if err := db.Set([]byte(dataCountKey), []byte(strconv.Itoa(0))); err != nil {
		panic(err)
	}
	return db
}

func main() {
	db := initDB()
	defer db.Close()

	if len(os.Args) != 4 {
		panic("./consumer [service-addr] [hlt-addr] [logger]")
	}

	serviceAddr := os.Args[1]
	hltAddr := os.Args[2]
	logger := os.Args[3]

	transferTraffic(db, serviceAddr, hltAddr, logger == "true")
}

func printLog(hasLog bool, msg error) {
	if !hasLog {
		return
	}
	fmt.Println(msg)
}

func transferTraffic(db database.IKVDatabase, serviceAddr, hltAddr string, hasLog bool) {
	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", hltAddr),
			&http.Client{Timeout: time.Minute},
			getMessageSettings(),
		),
	)

	for {
		time.Sleep(time.Second)

		countService, err := loadCountFromService(serviceAddr)
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
			if err := incrementCountInDB(db); err != nil {
				printLog(hasLog, err)
				continue
			}

			msg, err := loadMessageFromService(serviceAddr, i)
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

func loadMessageFromService(serviceAddr string, id uint64) (net_message.IMessage, error) {
	// build request to service
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf("http://%s/load?data_id=%d", serviceAddr, id),
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

	msg, err := net_message.LoadMessage(
		getMessageSettings(),
		string(msgStringAsBytes[1:]),
	)
	if err != nil {
		return nil, fmt.Errorf("message is nil")
	}

	return msg, nil
}

func loadCountFromService(serviceAddr string) (uint64, error) {
	// build request to service
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf("http://%s/size", serviceAddr),
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

	if err := db.Set([]byte(dataCountKey), []byte(strconv.FormatUint(count+1, 10))); err != nil {
		return err
	}

	return nil
}

func getMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: workSize,
	})
}
