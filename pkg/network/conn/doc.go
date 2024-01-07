// Package conn makes it possible to hide the structure of the true message by encrypting and adding random bytes.
//
/*
	NETWORK MESSAGE FORMAT

	E( LEM || LV || H(LEM||LV) || H(EM||V) ) || EM || V
	where
		LEM = L(EM)
		LV  = L(V)
		EM  = E(M)
		where
			E - encrypt (cipher-key)
			H - hmac (auth-key)
			L - length
			M - message bytes
			V - void bytes

	Scheme: https://github.com/number571/go-peer/blob/master/images/go-peer_layer1_net_message.jpg
*/
package conn
