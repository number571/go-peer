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

func getKeyHash(pI uint64) []byte {
	return []byte(fmt.Sprintf(cKeyHashTemplate, pI))
}

func getKeyMessage(pHash []byte) []byte {
	return []byte(fmt.Sprintf(cKeyMessageTemplate, pHash))
}
