package encoding

import (
	"bytes"
	"testing"
)

func TestBase64(t *testing.T) {
	data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 255, 254, 253, 252, 251, 250, 128, 127, 126, 125}
	edata := Base64Encode(data)
	if !bytes.Equal(data, Base64Decode(edata)) {
		t.Errorf("bytes not equals")
	}
}

func TestBytes(t *testing.T) {
	num := uint64(0xABCDEF0123456789)
	bnum := Uint64ToBytes(num)
	if num != BytesToUint64(bnum) {
		t.Errorf("numbers not equals")
	}
}
