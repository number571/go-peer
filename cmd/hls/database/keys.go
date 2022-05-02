package database

import (
	"fmt"
)

const (
	keyHash = "database.hashes[%s]"
)

func getKeyHash(key []byte) []byte {
	return []byte(fmt.Sprintf(keyHash, key))
}
