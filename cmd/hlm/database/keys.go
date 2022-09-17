package database

import "fmt"

const (
	cKeySizeTemplate    = "database.relations[%s-%s].size"
	cKeyMessageTemplate = "database.relations[%s-%s].messages[%d]"
)

func getKeySize(rel IRelation) []byte {
	return []byte(fmt.Sprintf(cKeySizeTemplate, rel.IAm().Address(), rel.Friend().Address()))
}

func getKeyMessage(rel IRelation, i uint64) []byte {
	return []byte(fmt.Sprintf(cKeyMessageTemplate, rel.IAm().Address(), rel.Friend().Address(), i))
}
