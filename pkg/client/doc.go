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
	3. 	CP = [ E( PubKB, K ), E( K, PubKA ), E( K, R ), E( K, P ), HP, E( K, S( PrivKA, HP ) ), W( C, HP ) ],
		where
			CP - encrypted message,
			E - encryption function,
			S - sign function,
			W - work confirmation function,
			C - the complexity of the work,
			PrivKX - private key of X participant.

	Steps of participant B:
	4. 	W( C, HP ) = PW( C, W( C, HP ) ),
		where
			PW - function of work checking.
		IF ≠, than protocol is interrupted.
	5. 	K = D( PrivKB, E( PubKB, K ) ),
		where
			D - decryption function.
		IF ≠, than protocol is interrupted.
	6. 	PubKA = D( K, E( K, PubKA ) ).
		IF ≠, than protocol is interrupted.
	7. 	HP = V( PubKA, D( K, E( K, S( PrivKA, HP ) ) ) ),
		where
			V - signature verification function.
		IF ≠, than protocol is interrupted.
	8. 	HP = H( D( K, E( K, R ) ) || D( K, E( K, P ) ) || PubKA || PubKB ),
		IF ≠, than protocol is interrupted.

	More information in article: https://github.com/number571/go-peer/blob/master/docs/monolithic_cryptographic_protocol.pdf
*/
package client
