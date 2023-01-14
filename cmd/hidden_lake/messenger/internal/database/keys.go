package database

import (
	"fmt"
)

const (
	cKeySizeTemplate    = "database[%s].friends[%s].size"
	cKeyMessageTemplate = "database[%s].friends[%s].messages[%d]"
)

func getKeySize(r IRelation) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		r.IAm().Address().String(),
		r.Friend().Address().String(),
	))
}

func getKeyMessage(r IRelation, i uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageTemplate,
		r.IAm().Address().String(),
		r.Friend().Address().String(),
		i,
	))
}
