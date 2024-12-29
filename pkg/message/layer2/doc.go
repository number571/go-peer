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

	More information in pkg/client
*/
package layer2
