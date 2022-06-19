package database

import "fmt"

const (
	cKeySizeTemplate    = "database.users[%X].size"
	cKeyMessageTemplate = "database.users[%X].messages[%d]"
	cKeyHashTemplate    = "database.hashes[%X]"
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
