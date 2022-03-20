package database

import "fmt"

const (
	KeySize    = "database.users[%s].size"
	KeyMessage = "database.users[%s].messages[%d]"
	KeyHash    = "database.hashes[%s]"
)

func GetKeySize(key []byte) []byte {
	return []byte(fmt.Sprintf(KeySize, key))
}

func GetKeyMessage(key []byte, i uint64) []byte {
	return []byte(fmt.Sprintf(KeyMessage, key, i))
}

func GetKeyHash(key []byte) []byte {
	return []byte(fmt.Sprintf(KeyHash, key))
}
