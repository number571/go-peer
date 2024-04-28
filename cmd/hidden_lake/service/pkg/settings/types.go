package settings

import "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"

type SFriend struct {
	FAliasName string `json:"alias_name"`
	FPublicKey string `json:"public_key"`
}

type SRequest struct {
	FReceiver string            `json:"receiver"` // alias_name
	FReqData  *request.SRequest `json:"req_data"`
}
