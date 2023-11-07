// Package message is used to create network messages for the purpose of confirming integrity and proof of work.
//
/*
	MESSAGE FORMAT

	P(HM) || HM || M
	where
		HM = H(M)
		where
			P - proof of work
			H - hmac (auth-key)
			M - message bytes
*/
package message
