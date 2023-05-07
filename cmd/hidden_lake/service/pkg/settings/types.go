package settings

type (
	SPrivKey = string
	SConnect = string
	SMessage = string
)

type SFriend struct {
	FAliasName string `json:"alias_name"`
	FPublicKey string `json:"public_key"`
}

type SRequest struct {
	FReceiver string `json:"receiver"` // public key
	FHexData  string `json:"hex_data"` // data in hex encode
}
