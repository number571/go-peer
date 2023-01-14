package database

import "fmt"

const (
	cKeyPointer         = "database.pointer"
	cKeyHashTemplate    = "database.hashes[%d]"
	cKeyMessageTemplate = "database.messages[%X]"
)

func getKeyPointer() []byte {
	return []byte(cKeyPointer)
}

func getKeyHash(i uint64) []byte {
	return []byte(fmt.Sprintf(cKeyHashTemplate, i))
}

func getKeyMessage(hash []byte) []byte {
	return []byte(fmt.Sprintf(cKeyMessageTemplate, hash))
}
