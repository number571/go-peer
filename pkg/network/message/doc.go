// Package message is used to create network messages for the purpose of confirming integrity and proof of work.
//
// The main purpose of the message is the possibility of retransmission with verification by the network key and hide the structure of the true message.
/*
	NETWORK MESSAGE FORMAT

	Sc || Sa || E( KDF(K,Sc), P(HLMV) || HLMV || LM || M || V )
	where
		HLMV = H( KDF(K,Sa), LM || M || V )
		LM   = L(M)
		Sc   = G(N)
		Sa   = G(N)
		where
			KDF - key derivation function
			H   - hmac
			K   - network key
			Sa  - auth salt
			M   - message bytes
			L   - length
			G   - prng
			Sc  - cipher salt
			P   - proof of work
			E   - encrypt
			V   - void bytes
			N   - num random bytes

	Scheme: https://github.com/number571/go-peer/blob/master/images/go-peer_layer1_net_message.jpg
*/
package message
