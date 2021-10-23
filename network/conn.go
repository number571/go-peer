package network

import (
	"bytes"
	"net"

	"github.com/number571/gopeer"
	"github.com/number571/gopeer/crypto"
	"github.com/number571/gopeer/local"
)

func readMessage(conn net.Conn) *local.Message {
	const (
		UINT64_SIZE = 8 // bytes
	)

	var (
		pack   []byte
		size   = uint(0)
		buflen = make([]byte, UINT64_SIZE)
		buffer = make([]byte, gopeer.Get("BUFF_SIZE").(uint))
	)

	length, err := conn.Read(buflen)
	if err != nil {
		return nil
	}
	if length != UINT64_SIZE {
		return nil
	}

	mustLen := local.Package(buflen).BytesToSize()
	if mustLen > gopeer.Get("PACK_SIZE").(uint) {
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

	diff := uint8(gopeer.Get("POWS_DIFF").(uint))

	puzzle := crypto.NewPuzzle(diff)
	if !puzzle.Verify(msg.Body.Hash, msg.Body.Npow) {
		return nil
	}

	return msg
}
