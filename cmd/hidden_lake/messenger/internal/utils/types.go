package utils

type ILanguage int

const (
	CLangENG = 0
	CLangRUS = 1
	CLangESP = 2
)

type SMessageInfo struct {
	FFileName  string `json:"filename"` // can be ""
	FMessage   string `json:"message"`
	FTimestamp string `json:"timestamp"`
}
