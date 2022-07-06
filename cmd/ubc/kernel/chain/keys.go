package chain

import (
	"fmt"

	"github.com/number571/go-peer/cmd/ubc/kernel/settings"
)

func getKeyHeight() []byte {
	key := settings.GSettings.Get(settings.CMaskHeig).(string)
	return []byte(key)
}

func getKeyBlock(height uint64) []byte {
	key := settings.GSettings.Get(settings.CMaskBlck).(string)
	return []byte(fmt.Sprintf(key, height))
}

func getKeyTX(hash []byte) []byte {
	key := settings.GSettings.Get(settings.CMaskTrns).(string)
	return []byte(fmt.Sprintf(key, hash))
}
