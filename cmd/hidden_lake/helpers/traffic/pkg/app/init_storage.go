package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/cache"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/storage"
	"github.com/number571/go-peer/pkg/storage/database"
)

func (p *sApp) initStorage(pDatabase database.IKVDatabase) {
	cfgSettings := p.fConfig.GetSettings()
	p.fStorage = storage.NewMessageStorage(
		cfgSettings,
		pDatabase,
		cache.NewLRUCache(cfgSettings.GetMessagesCapacity()),
	)
}
