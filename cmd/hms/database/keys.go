package database

import "fmt"

const (
	cKeySizeTemplate    = "database.users[%s].size"
	cKeyMessageTemplate = "database.users[%s].messages[%d]"
	cKeyHashTemplate    = "database.hashes[%s]"
)

func getKeySize(key []byte) []byte {
	return []byte(fmt.Sprintf(cKeySizeTemplate, key))
}

func getKeyMessage(key []byte, i uint64) []byte {
	return []byte(fmt.Sprintf(cKeyMessageTemplate, key, i))
}

func getKeyHash(key []byte) []byte {
	return []byte(fmt.Sprintf(cKeyHashTemplate, key))
}
