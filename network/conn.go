package network

import (
	"bytes"
	"net"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"
)

/*
	send может передавать информацию фрагментированно?
	recv может принимать информацию фрагментированно?

	A -> 512 [X] -> C ?
	A -> 512 [X] -> B ?
*/

func (node *sNode) readMessage(conn net.Conn) local.IMessage {
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

	mustLen := local.LoadPackage(bufsize).BytesToSize()
	if mustLen > node.fClient.Settings().Get(settings.SizePack) {
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

	return node.initialCheck(local.LoadPackage(pack).ToMessage())
}

func (node *sNode) initialCheck(msg local.IMessage) local.IMessage {
	if msg == nil {
		return nil
	}

	if len(msg.Body().Hash()) != crypto.HashSize {
		return nil
	}

	diff := node.fClient.Settings().Get(settings.SizeWork)
	puzzle := crypto.NewPuzzle(diff)
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		return nil
	}

	return msg
}
