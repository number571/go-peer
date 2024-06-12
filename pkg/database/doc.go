// Package database is a key-value database with integrated password encryption functions.
//
/*
	DATABASE INIT FORMAT

	Ak, Ck <- KDF(S, P)
	S || R || H[Ak](R)
	where
		H  - hmac
		S  - salt value
		R  - random value
		P  - password
		Ak - auth-key
		Ck - cipher-key
*/
/*
	DATABASE SET FORMAT

	{H[Ak](K) : (H[Ak](EV) || EV)}
	where
		EV = E[Ck](V)
		where
			E  - encrypt
			H  - hmac
			K  - database key
			V  - database value
			Ak - auth-key
			Ck - cipher-key
*/
package database
