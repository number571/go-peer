package settings

type ISettings interface {
	Set(uint64, uint64) ISettings
	Get(uint64) uint64
}
