package database

import (
	"fmt"
)

const (
	cKeySizeTemplate          = "database[%s].friends[%s].size"
	cKeyMessageByUIDTemplate  = "database[%s].friends[%s].messages[uid=%x]"
	cKeyMessageByEnumTemplate = "database[%s].friends[%s].messages[enum=%d]"
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

func getKeyMessageByUID(pR IRelation, pBlockUID [cBlockUIDSize]byte) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByUIDTemplate,
		pR.IAm().GetAddress().ToString(),
		pR.Friend().GetAddress().ToString(),
		pBlockUID[:],
	))
}
