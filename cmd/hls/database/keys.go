package database

import "fmt"

const (
	KeyHash = "database.hashes[%s]"
)

func GetKeyHash(key []byte) []byte {
	return []byte(fmt.Sprintf(KeyHash, key))
}
