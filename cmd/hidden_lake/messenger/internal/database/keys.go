package database

import (
	"fmt"
)

const (
	cKeySizeTemplate          = "database[%s].friends[%s].size"
	cKeyMessageByEnumTemplate = "database[%s].friends[%s].messages[enum=%d]"
	cKeyMessageByHashTemplate = "database[%s].friends[%s].messages[hash=%s]"
)

func getKeySize(r IRelation) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		r.IAm().Address().String(),
		r.Friend().Address().String(),
	))
}

func getKeyMessageByEnum(r IRelation, i uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByEnumTemplate,
		r.IAm().Address().String(),
		r.Friend().Address().String(),
		i,
	))
}

func getKeyMessageByHash(r IRelation, hash string) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByHashTemplate,
		r.IAm().Address().String(),
		r.Friend().Address().String(),
		hash,
	))
}
