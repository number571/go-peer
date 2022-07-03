package settings

type ISettings interface {
	Set(uint64, uint64) ISettings
	Get(uint64) uint64
}

const (
	CMaskNetw uint64 = iota + 1
	CMaskRout
	CMaskPing
	CMaskPsdo
	CMaskPasw
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
	CSizePasw
)

// 3-bit mask for password
const (
	CPaswAplh = 0b100
	CPaswNumr = 0b010
	CPaswSpec = 0b001
)

const (
	CSizeUint64 = 8
)
