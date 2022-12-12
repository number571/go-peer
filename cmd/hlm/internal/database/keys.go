package database

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	cKeySizeTemplate    = "database.friends[%s].size"
	cKeyMessageTemplate = "database.friends[%s].messages[%d]"
)

func getKeySize(pubKey asymmetric.IPubKey) []byte {
	return []byte(fmt.Sprintf(cKeySizeTemplate, pubKey.Address().String()))
}

func getKeyMessage(pubKey asymmetric.IPubKey, i uint64) []byte {
	return []byte(fmt.Sprintf(cKeyMessageTemplate, pubKey.Address().String(), i))
}
