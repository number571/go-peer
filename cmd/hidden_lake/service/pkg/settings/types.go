package settings

type SFriend struct {
	FAliasName string `json:"alias_name"`
	FPublicKey string `json:"public_key"`
}

type SRequest struct {
	FReceiver string `json:"receiver"` // alias_name
	FReqData  string `json:"req_data"` // data in hex encode
}
