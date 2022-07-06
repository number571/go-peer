package settings

import "sync"

const (
	CSizeAkey uint64 = iota + 1
	CSizeTrns
	CSizePayl
	CSizeMemp
	CPathBlck
	CPathTrns
	CPathMemp
	CMaskHeig
	CMaskBlck
	CMaskTrns
	CMaskPref
)

var (
	_ iSettings = &sSettings{}

	// singleton
	GSettings = newSettings()
)

type sSettings struct {
	fMutex   sync.Mutex
	fMapping map[uint64]interface{}
}

func newSettings() iSettings {
	return &sSettings{
		fMapping: defaultSettings(),
	}
}

func (s *sSettings) Set(k uint64, v interface{}) iSettings {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	s.fMapping[k] = v
	return s
}

func (s *sSettings) Get(k uint64) interface{} {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	v, ok := s.fMapping[k]
	if !ok {
		panic("settings: value undefined")
	}

	return v
}

func defaultSettings() map[uint64]interface{} {
	return map[uint64]interface{}{
		CSizeAkey: uint64(1024), // num bits
		CSizeTrns: uint64(32),   // num txs in block
		CSizePayl: uint64(1024), // num bytes in tx.payload
		CSizeMemp: uint64(512),  // max num txs in mempool
		CPathBlck: "blocks.db",
		CPathTrns: "txs.db",
		CPathMemp: "mempool.db",
		CMaskHeig: "chain.blocks.height",
		CMaskBlck: "chain.blocks.block[%d]",
		CMaskTrns: "chain.txs.tx[%X]",
		CMaskPref: "chain.txs.tx[",
	}
}

/*
// CSizeAkey: uint64(4096), // num bits
// CSizeTrns: uint64(1024), // num txs in block
// CSizePayl: uint64(2048), // num bytes in tx.payload
// CSizeMemp: uint64(8192), // max num txs in mempool
// CPathBlck: "blocks.db",
// CPathTrns: "txs.db",
// CPathMemp: "mempool.db",
// CMaskHeig: "chain.blocks.height",
// CMaskBlck: "chain.blocks.block[%d]",
// CMaskTrns: "chain.txs.tx[%X]",
// CMaskPref: "chain.mempool.tx[",
*/
