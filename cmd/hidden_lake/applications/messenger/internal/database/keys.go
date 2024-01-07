package database

import (
	"fmt"
)

const (
	cKeySizeTemplate          = "database[%s].friends[%s].size"
	cKeyMessageByEnumTemplate = "database[%s].friends[%s].messages[enum=%d]"
)

func getKeySize(pR IRelation) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		pR.IAm().GetHasher().ToString(),
		pR.Friend().GetHasher().ToString(),
	))
}

func getKeyMessageByEnum(pR IRelation, pI uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByEnumTemplate,
		pR.IAm().GetHasher().ToString(),
		pR.Friend().GetHasher().ToString(),
		pI,
	))
}
