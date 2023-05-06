package settings

type SPrivKey struct {
	FPrivKey string `json:"priv_key"`
}

type SConnect struct {
	FConnect string `json:"connect"`
}

type SMessage struct {
	FHexMessage string `json:"hex_message"`
}

type SFriend struct {
	FAliasName string `json:"alias_name"`
	FPublicKey string `json:"public_key"`
}

type SRequest struct {
	FReceiver string `json:"receiver"` // public key
	FHexData  string `json:"hex_data"` // data in hex encode
}
