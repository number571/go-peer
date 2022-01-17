package network

import (
	"bytes"
	"net"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"
)

func readMessage(conn net.Conn) *local.Message {
	const (
		SizeUint64 = 8 // bytes
	)

	var (
		pack   []byte
		size   = uint(0)
		buflen = make([]byte, SizeUint64)
		buffer = make([]byte, settings.Get("BUFF_SIZE").(uint))
	)

	length, err := conn.Read(buflen)
	if err != nil {
		return nil
	}
	if length != SizeUint64 {
		return nil
	}

	mustLen := local.Package(buflen).BytesToSize()
	if mustLen > settings.Get("PACK_SIZE").(uint) {
		return nil
	}

	for {
		length, err = conn.Read(buffer)
		if err != nil {
			return nil
		}

		size += uint(length)
		if size > mustLen {
			return nil
		}

		pack = bytes.Join(
			[][]byte{
				pack,
				buffer[:length],
			},
			[]byte{},
		)

		if size == mustLen {
			break
		}
	}

	return initialCheck(local.Package(pack).Deserialize())
}

func initialCheck(msg *local.Message) *local.Message {
	if msg == nil {
		return nil
	}

	if len(msg.Body.Hash) != crypto.HashSize {
		return nil
	}

	diff := uint8(settings.Get("POWS_DIFF").(uint))

	puzzle := crypto.NewPuzzle(diff)
	if !puzzle.Verify(msg.Body.Hash, msg.Body.Npow) {
		return nil
	}

	return msg
}
