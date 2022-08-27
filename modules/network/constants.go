package network

import (
	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/modules/encoding"
)

// in bytes
const (
	cSizeUint = encoding.CSizeUint64
	cSizeHash = hashing.CSHA256Size
	cSizeHead = encoding.CSizeUint64
)

const (
	cPointSize = 0
	cPointHash = cPointSize + cSizeUint
	cPointHead = cPointHash + cSizeHash
	cPointBody = cPointHead + cSizeHead
)

const (
	// cBeginSize = cPointSize
	// cEndSize   = cPointHash

	cBeginHash = cPointHash
	cEndHash   = cPointHead

	cBeginHead = cPointHead
	cEndHead   = cPointBody

	// cBeginBody = cPointBody
	// cEndBody = [cBeginBody:]
)
