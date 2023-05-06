package settings

// TODO: to string as SMessage
type SPrivKey struct {
	FPrivKey string `json:"priv_key"`
}

// TODO: to string as SMessage
type SConnect struct {
	FConnect string `json:"connect"`
}

type SMessage = []byte

type SFriend struct {
	FAliasName string `json:"alias_name"`
	FPublicKey string `json:"public_key"`
}

type SRequest struct {
	FReceiver string `json:"receiver"` // public key
	FHexData  string `json:"hex_data"` // data in hex encode
}
