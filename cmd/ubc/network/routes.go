package network

import (
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/network"
)

func handlePushBlock(node INode, conn network.IConn, npld payload.IPayload) {

}

func handlePushTransaction(node INode, conn network.IConn, npld payload.IPayload) {

}

func handleLoadHeight(node INode, conn network.IConn, npld payload.IPayload) {
	conn.Write(network.NewMessage(
		payload.NewPayload(
			cMaskLoadHeight,
			encoding.Uint64ToBytes(55),
		),
	))
}

func handleLoadBlock(node INode, conn network.IConn, npld payload.IPayload) {

}

func handleLoadTransaction(node INode, conn network.IConn, npld payload.IPayload) {

}
