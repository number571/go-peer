package settings

type SResponse struct {
	Result []byte `json:"result"`
	Return int    `json:"return"`
}

type SSizeRequest struct {
	Receiver []byte `json:"receiver"`
}

type SLoadRequest struct {
	Receiver []byte `json:"receiver"`
	Index    uint64 `json:"index"`
}

type SPushRequest struct {
	Receiver []byte `json:"receiver"`
	Package  []byte `json:"package"`
}
