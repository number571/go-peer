package settings

type ISettings interface {
	Set(uint64, uint64) ISettings
	Get(uint64) uint64
}

const (
	CMaskRout uint64 = iota + 1
	CMaskPing
	CTimeWait
	CTimePreq
	CTimePrsp
	CTimePing
	CSizePsdo
	CSizeRtry
	CSizeWork
	CSizeConn
	CSizePack
	CSizeMapp
	CSizeSkey
	CSizeBmsg
)
