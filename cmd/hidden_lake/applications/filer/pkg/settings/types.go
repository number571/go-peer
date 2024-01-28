package settings

type SFileInfo struct {
	FName string `json:"name"`
	FHash string `json:"hash"`
	FSize uint64 `json:"size"`
}
