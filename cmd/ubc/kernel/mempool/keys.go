package mempool

import (
	"fmt"

	"github.com/number571/go-peer/cmd/ubc/kernel/settings"
)

func getKeyHeight() []byte {
	key := settings.GSettings.Get(settings.CMaskHeig).(string)
	return []byte(key)
}

func getKeyTX(hash []byte) []byte {
	key := settings.GSettings.Get(settings.CMaskTrns).(string)
	return []byte(fmt.Sprintf(key, hash))
}
