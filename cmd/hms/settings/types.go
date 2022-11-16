package settings

type SResponse struct {
	FResult []byte `json:"result"`
	FReturn int    `json:"return"`
}

type SSizeRequest struct {
	FReceiver []byte `json:"receiver"`
}

type SLoadRequest struct {
	FReceiver []byte `json:"receiver"`
	FIndex    uint64 `json:"index"`
}

type SPushRequest struct {
	FReceiver []byte `json:"receiver"`
	FPackage  []byte `json:"package"`
}
