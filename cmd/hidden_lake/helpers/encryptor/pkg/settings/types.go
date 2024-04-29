package settings

type SContainer struct {
	FPublicKey string `json:"public_key"`
	FPldHead   uint64 `json:"pld_head"`
	FHexData   string `json:"hex_data"`
}
