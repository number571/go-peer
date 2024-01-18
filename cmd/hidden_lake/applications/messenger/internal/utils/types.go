package utils

type SMessageInfo struct {
	FFileName  string `json:"filename"` // can be ""
	FTimestamp string `json:"timestamp"`
	FPseudonym string `json:"pseudonym"`
	FMessage   string `json:"message"`
}
