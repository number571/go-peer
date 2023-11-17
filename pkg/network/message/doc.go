// Package message is used to create network messages for the purpose of confirming integrity and proof of work.
//
// The main purpose of the message is the possibility of retransmission with verification by the network key.
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
