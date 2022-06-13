package network

import (
	"bytes"
	"net"

	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/crypto/puzzle"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/settings"
)

func (node *sNode) readMessage(conn net.Conn) message.IMessage {
	const (
		SizeUint64 = 8 // bytes
	)

	var (
		pack    []byte
		size    = uint64(0)
		bufsize = make([]byte, SizeUint64)
	)

	length, err := conn.Read(bufsize)
	if err != nil {
		return nil
	}
	if length != SizeUint64 {
		return nil
	}

	mustLen := message.LoadPackage(bufsize).BytesToSize()
	if mustLen > node.fClient.Settings().Get(settings.CSizePack) {
		return nil
	}

	buffer := make([]byte, mustLen)
	for {
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
		if size >= mustLen {
			break
		}
	}

	return node.initialCheck(message.LoadPackage(pack).ToMessage())
}

func (node *sNode) initialCheck(msg message.IMessage) message.IMessage {
	if msg == nil {
		return nil
	}

	if len(msg.Body().Hash()) != hashing.HashSize {
		return nil
	}

	diff := node.fClient.Settings().Get(settings.CSizeWork)
	puzzle := puzzle.NewPoWPuzzle(diff)
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		return nil
	}

	return msg
}
