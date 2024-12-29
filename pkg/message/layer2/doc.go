// Package message used as an encapsulated ciphertext.
//
/*
	MESSAGE FORMAT

	E( PubK, K ) || E( K, M )
	where
		PubK - public key
		K - secret key
		M - message bytes
		E - encrypt

	Scheme: https://github.com/number571/go-peer/blob/master/images/go-peer_layer2_message.jpg
*/
package layer2
