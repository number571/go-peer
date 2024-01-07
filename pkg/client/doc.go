// Package client makes it possible to encrypt and decrypt messages using a monolithic cryptographic protocol.
//
/*
	CLIENT MESSAGE PROTOCOL

	Protocol participants:
		A - sender,
		B - receiver.

	Steps of participant A:
	1. 	K = G( N ), R = G( N ),
		where
			G - generator pseudo random bytes,
			N - count of bytes for generator,
			K - encryption session key,
			R - pseudo random bytes (salt).
	2. 	HP = H( R || P || PubKA || PubKB ),
		where
			HP - message hash,
			H - hash function,
			P - plaintext,
			PubKX - public key of X participant.
	3. 	CP = [ E( PubKB, K ), E( K, PubKA ), E( K, R ), E( K, P ), E( K, HP ), E( K, S( PrivKA, HP ) ) ],
		where
			CP - encrypted message,
			E - encryption function,
			S - sign function,
			PrivKX - private key of X participant.

	Steps of participant B:
	4. 	K = D( PrivKB, E( PubKB, K ) ),
		where
			D - decryption function.
		IF ≠, than protocol is interrupted.
	5. 	PubKA = D( K, E( K, PubKA ) ).
		IF ≠, than protocol is interrupted.
	6. 	HP = V( PubKA, D( K, E( K, S( PrivKA, D( K, E( K, HP) ) ) ) ) ),
		where
			V - signature verification function.
		IF ≠, than protocol is interrupted.
	7. 	HP = H( D( K, E( K, R ) ) || D( K, E( K, P ) ) || PubKA || PubKB ),
		IF ≠, than protocol is interrupted.

	More information in article: https://github.com/number571/go-peer/blob/master/docs/monolithic_cryptographic_protocol.pdf
	Scheme: https://github.com/number571/go-peer/blob/master/images/go-peer_layer2_message.jpg
*/
package client
