package mempool

import (
	"fmt"
)

const (
	cHeight     = "chain.blocks.height"
	cTxTemplate = "chain.txs.tx[%X]"
	cPrefix     = "chain.txs.tx["
)

func getKeyHeight() []byte {
	return []byte(cHeight)
}

func getKeyTX(hash []byte) []byte {
	return []byte(fmt.Sprintf(cTxTemplate, hash))
}
