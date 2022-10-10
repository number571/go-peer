package settings

type SResponse struct {
	FResult string `json:"result"`
	FReturn int    `json:"return"`
}

type SBroadcast SRequest
type SRequest struct {
	FReceiver string `json:"receiver"` // public key
	FHexData  string `json:"hex_data"` // data in hex encode
}
