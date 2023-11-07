// Package database is a key-value database with integrated password encryption functions.
//
/*
	DATABASE ENCRYPTION FORMAT

	H(EM) || EM
	where
		EM = E(M)
		where
			E - encrypt (cipher-key)
			H - hmac (auth-key)
			M - message bytes
*/
package database
