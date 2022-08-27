package settings

type SStatusNetwork struct {
	PubKey string `json:"pub_key"`
	Online bool   `json:"online"`
}

type SStatusResponse struct {
	PubKey  string           `json:"pub_key"`
	Network []SStatusNetwork `json:"network"`
	SResponse
}

type SResponse struct {
	Result string `json:"result"`
	Return int    `json:"return"`
}

type SRequest struct {
	Receiver string `json:"receiver"` // public key
	HexData  string `json:"hex_data"` // data in hex encode
}
