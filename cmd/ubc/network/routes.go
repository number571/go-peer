package network

import (
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/payload"
)

func handlePushBlock(node INode, conn network.IConn, npld payload.IPayload) {

}

func handlePushTransaction(node INode, conn network.IConn, npld payload.IPayload) {

}

func handleLoadHeight(node INode, conn network.IConn, npld payload.IPayload) {
	res := encoding.Uint64ToBytes(55)
	conn.Write(network.NewMessage(
		payload.NewPayload(
			cMaskLoadHeight,
			res[:],
		),
		[]byte{},
	))
}

func handleLoadBlock(node INode, conn network.IConn, npld payload.IPayload) {

}

func handleLoadTransaction(node INode, conn network.IConn, npld payload.IPayload) {

}
