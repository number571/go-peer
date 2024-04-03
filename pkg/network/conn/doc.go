// Package conn makes it possible to hide the structure of the true message by encrypting and adding random bytes.
//
/*
	NETWORK CONN FORMAT

	L(M) || M
	where
		L - length (uint64)
		M - message bytes
*/
package conn
