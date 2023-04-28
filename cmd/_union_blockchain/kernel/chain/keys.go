package chain

import (
	"fmt"
)

const (
	cHeight        = "chain.blocks.height"
	cBlockTemplate = "chain.blocks.block[%d]"
	cTxTemplate    = "chain.txs.tx[%X]"
)

func getKeyHeight() []byte {
	return []byte(cBlockTemplate)
}

func getKeyBlock(height uint64) []byte {
	return []byte(fmt.Sprintf(cBlockTemplate, height))
}

func getKeyTX(hash []byte) []byte {
	return []byte(fmt.Sprintf(cTxTemplate, hash))
}
