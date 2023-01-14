package settings

type SResponse struct {
	FResult string `json:"result"`
	FReturn int    `json:"return"`
}

type SPrivKey struct {
	FPrivKey string `json:"priv_key"`
}

type SConnect struct {
	FConnect string `json:"connect"`
}

type SFriend struct {
	FAliasName string `json:"alias_name"`
	FPublicKey string `json:"public_key"`
}

type SPush struct {
	FReceiver string `json:"receiver"` // public key
	FHexData  string `json:"hex_data"` // data in hex encode
}
