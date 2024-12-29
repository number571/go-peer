// Package message used as a storage and loading of encrypted messages.
//
// The package allows initializing verification of the correctness of
// the message by the hash length and proof of work.
/*
	NETWORK MESSAGE FORMAT

	E( K, P(HM) || HM || M )
	where
		HM = H( K, M )
		where
			H - HMAC
			K - network key
			M - message bytes
			P - proof of work
			E - encrypt

	Scheme: https://github.com/number571/go-peer/blob/master/images/go-peer_layer2_message.jpg
*/
package layer2
