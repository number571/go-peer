package network

import (
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/settings"
)

// in bytes
const (
	cSizeUint = settings.CSizeUint64
	cSizeHash = hashing.GSHA256Size
	cSizeHead = settings.CSizeUint64
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
