package utils

type SSubscribe struct {
	FAddress string `json:"address"`
}

type SMessage struct {
	FFileName  string `json:"filename"` // can be ""
	FTimestamp string `json:"timestamp"`
	FPseudonym string `json:"pseudonym"`
	FMainData  string `json:"maindata"`
}
