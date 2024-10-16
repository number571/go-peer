// Package message is used to create network messages for the purpose of confirming integrity and proof of work.
//
// The main purpose of the message is the possibility of retransmission with verification by the network key and hide the structure of the true message.
/*
	NETWORK MESSAGE FORMAT

	E( K, P(HLMR) || HLMR || L(M) || M || L(R) || R )
	where
		HLMR = H( K, L(M) || M || L(R) || R )
		where
			H - HMAC
			K - network key
			M - message bytes
			L - length
			P - proof of work
			E - encrypt
			R - rand bytes

	Scheme: https://github.com/number571/go-peer/blob/master/images/go-peer_layer1_net_message.jpg
*/
package message
