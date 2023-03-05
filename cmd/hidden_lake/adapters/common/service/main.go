package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage/database"
)

const (
	databasePath = "service.db"
	dataCountKey = "count_service"
)

var (
	mtx sync.Mutex
	db  database.IKeyValueDB
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

	if len(os.Args) != 2 {
		panic("./service [port]")
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/size", sizePage)
	http.HandleFunc("/push", pushPage)
	http.HandleFunc("/load", loadPage)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func pushPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Response(w, 2, "failed: incorrect method")
		return
	}

	res, err := io.ReadAll(r.Body)
	if err != nil {
		api.Response(w, 3, "failed: read body")
		return
	}

	if err := pushDataToDB(res); err != nil {
		api.Response(w, 4, "failed: push to database")
		return
	}

	api.Response(w, 1, "success: push to database")
}

func sizePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.Response(w, 2, "failed: incorrect method")
		return
	}

	count, err := countOfDataInDB()
	if err != nil {
		api.Response(w, 3, "failed: read count of data")
		return
	}

	api.Response(w, 1, fmt.Sprintf("%d", count))
}

func loadPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.Response(w, 2, "failed: incorrect method")
		return
	}

	query := r.URL.Query()
	strDataID := query.Get("data_id")

	dataID, err := strconv.ParseUint(strDataID, 10, 64)
	if err != nil {
		api.Response(w, 3, "failed: decode data_id")
		return
	}

	data, err := loadDataFromDB(dataID)
	if err != nil {
		api.Response(w, 3, "failed: load data by data_id")
		return
	}

	api.Response(w, 1, encoding.HexEncode(data))
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

func loadDataFromDB(dataID uint64) ([]byte, error) {
	data, err := db.Get([]byte(fmt.Sprintf("%d", dataID)))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func pushDataToDB(data []byte) error {
	mtx.Lock()
	defer mtx.Unlock()

	count, err := countOfDataInDB()
	if err != nil {
		return err
	}

	if err := db.Set([]byte(fmt.Sprintf("%d", count)), data); err != nil {
		return err
	}

	if err := db.Set([]byte(dataCountKey), []byte(fmt.Sprintf("%d", count+1))); err != nil {
		return err
	}

	return nil
}
