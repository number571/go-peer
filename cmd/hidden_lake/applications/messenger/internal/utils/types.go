package utils

type SMessageInfo struct {
	FFileName  string `json:"filename"` // can be ""
	FTimestamp string `json:"timestamp"`
	FSenderID  string `json:"sender_id"`
	FMessage   string `json:"message"`
}
