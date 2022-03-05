package network

import (
	"bytes"
	"net"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"
)

func (node *NodeT) readMessage(conn net.Conn) local.Message {
	const (
		SizeUint64 = 8 // bytes
	)

	var (
		pack   []byte
		size   = uint64(0)
		buflen = make([]byte, SizeUint64)
	)

	length, err := conn.Read(buflen)
	if err != nil {
		return nil
	}
	if length != SizeUint64 {
		return nil
	}

	mustLen := local.Package(buflen).BytesToSize()
	if mustLen > node.Client().Settings().Get(settings.SizePack) {
		return nil
	}

	for {
		buffer := make([]byte, mustLen-size)

		length, err = conn.Read(buffer)
		if err != nil {
			return nil
		}

		pack = bytes.Join(
			[][]byte{
				pack,
				buffer[:length],
			},
			[]byte{},
		)

		size += uint64(length)
		if size == mustLen {
			break
		}
	}

	return node.initialCheck(local.Package(pack).Deserialize())
}

func (node *NodeT) initialCheck(msg local.Message) local.Message {
	if msg == nil {
		return nil
	}

	if len(msg.Body.Hash) != crypto.HashSize {
		return nil
	}

	diff := node.Client().Settings().Get(settings.SizeWork)
	puzzle := crypto.NewPuzzle(diff)
	if !puzzle.Verify(msg.Body.Hash, msg.Body.Proof) {
		return nil
	}

	return msg
}
