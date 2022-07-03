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
	Result []byte `json:"result"`
	Return int    `json:"return"`
}

type SRequest struct {
	Receiver string `json:"receiver"` // public key
	Data     []byte `json:"data"`
}
