package messenger

import (
	"bytes"
	"os"
	"testing"
)

func TestMessenger(t *testing.T) {
	t.Parallel()

	webBytes, err := os.ReadFile("web/web.go")
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Contains(webBytes, []byte("cUsedEmbedFS = true")) {
		t.Error("cUsedEmbedFS should be = true")
		return
	}
}
