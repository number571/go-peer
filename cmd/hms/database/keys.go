package database

import "fmt"

const (
	keySize    = "database.users[%s].size"
	keyMessage = "database.users[%s].messages[%d]"
	keyHash    = "database.hashes[%s]"
)

func getKeySize(key []byte) []byte {
	return []byte(fmt.Sprintf(keySize, key))
}

func getKeyMessage(key []byte, i uint64) []byte {
	return []byte(fmt.Sprintf(keyMessage, key, i))
}

func getKeyHash(key []byte) []byte {
	return []byte(fmt.Sprintf(keyHash, key))
}
