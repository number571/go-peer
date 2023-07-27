package settings

type SFriend struct {
	FAliasName string `json:"alias_name"`
	FPublicKey string `json:"public_key"`
}

type SRequest struct {
	FReceiver string `json:"receiver"` // public key
	FHexData  string `json:"hex_data"` // data in hex encode
}

type SPrivKey struct {
	FEncPrivKey    string `json:"enc_priv_key"`
	FEncSessionKey string `json:"enc_session_key"`
}
