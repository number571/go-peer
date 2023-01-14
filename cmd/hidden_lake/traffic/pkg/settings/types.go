package settings

type SResponse struct {
	FResult string `json:"result"`
	FReturn int    `json:"return"`
}

type SLoadRequest struct {
	FHash string `json:"hash"`
}

type SPushRequest struct {
	FMessage string `json:"message"`
}
