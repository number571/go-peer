package testutils

import (
	"math/rand"
	"time"
)

func PseudoRandomBytes(pSeed int) []byte {
	r := rand.New(rand.NewSource(int64(pSeed))) //nolint:gosec
	result := make([]byte, 0, 16)
	for i := 0; i < 16; i++ {
		result = append(result, byte(r.Intn(256)))
	}
	return result
}

func TryN(n int, t time.Duration, f func() error) error {
	var resErr error
	for i := 0; i < n; i++ {
		time.Sleep(t)
		if err := f(); err != nil {
			resErr = err
			continue
		}
		return nil
	}
	return resErr
}

var (
	TgAddrs = [20]string{
		"localhost:8000",
		"localhost:8001",
		"localhost:8002",
		"localhost:8003",
		"localhost:8004",
		"localhost:8005",
		"localhost:8006",
		"localhost:8007",
		"localhost:8008",
		"localhost:8009",

		"localhost:8010",
		"localhost:8011",
		"localhost:8012",
		"localhost:8013",
		"localhost:8014",
		"localhost:8015",
		"localhost:8016",
		"localhost:8017",
		"localhost:8018",
		"localhost:8019",
	}
)
