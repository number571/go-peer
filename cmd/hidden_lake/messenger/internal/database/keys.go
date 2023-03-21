package database

import (
	"fmt"
)

const (
	cKeySizeTemplate          = "database[%s].friends[%s].size"
	cKeyMessageByEnumTemplate = "database[%s].friends[%s].messages[enum=%d]"
	cKeyMessageByHashTemplate = "database[%s].friends[%s].messages[hash=%s]"
)

func getKeySize(pR IRelation) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		pR.IAm().GetAddress().ToString(),
		pR.Friend().GetAddress().ToString(),
	))
}

func getKeyMessageByEnum(pR IRelation, pI uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByEnumTemplate,
		pR.IAm().GetAddress().ToString(),
		pR.Friend().GetAddress().ToString(),
		pI,
	))
}

func getKeyMessageByHash(pR IRelation, pHash string) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByHashTemplate,
		pR.IAm().GetAddress().ToString(),
		pR.Friend().GetAddress().ToString(),
		pHash,
	))
}
