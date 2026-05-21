// Package client makes it possible to encrypt and decrypt messages using a monolithic cryptographic protocol.
//
/*
	CLIENT MESSAGE PROTOCOL

	Protocol participants:
		A - sender,
		B - receiver.

	Steps of participant A:
	1. 	S = G( N ),
		where
			S - pseudo random bytes (salt).
			G - generator pseudo random bytes,
			N - count of bytes for generator,
	2. 	K' = HMAC( K, S ),
		where
			K' - ...,
			K  - session key,
			HMAC - message authentication code function,
	3. 	C = [S, C' = E( K', P )],
		where
			P - plaintext,
			C - encrypted message,
			E - encryption function (AEAD).

	Steps of participant B:
	4. 	FOR K IN KEY_LIST:
			P = D( HMAC( K, S ), C' ),
			where
				D - decryption function.
	5. 	IF P = NULL, than NEXT ITER (4),
	6. 	ELSE RETURN P.
*/
package symmetric
