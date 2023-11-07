// Package storage is designed to securely store the most sensitive information (keys, passwords).
//
/*
	STORAGE ALGORITHM 'CREATE'

	1. 	S = G(N)
		where
			G - generator pseudo random bytes,
			N - count of bytes for generator,
			S - pseudo random bytes (salt).
	2. 	K = KDF(P, S)
		where
			KDF - key derivation function
			K - encryption key,
			P - password,
	3.	EM = E(K, VM)
		where
			E - encryption function,
			EM - encrypted map/storage,
			VM - void map/storage.
*/
/*
	STORAGE ALGORITHM 'SET'

	1. 	M = D(K, EM)
		where
			D - decryption function
			M - map/storage
	2. 	Km = KDF(Ki, S)
		where
			Km - key map/storage
			Ki - input key
	3. 	Vm = E(Km, Vi)
		where
			Vm - value map/storage
			Vi - input value
	4. 	M = SET(H(Km), Vm)
		where
			H - hash function
			SET - set H(Km),Vm to map/storage
	5. 	EM = E(K, M)
*/
/*
	STORAGE ALGORITHM 'GET'

	1. 	M = D(K, EM)
		where
			M - map/storage
	2. 	Km = KDF(Ki, S)
		where
			Km - key map/storage
			Ki - input key
	3. 	Vm = GET(H(Km))
		where
			GET - get Vm from map/storage by Km
	4. 	Vi = D(Km, Vm)
*/
/*
	STORAGE ALGORITHM 'DEL'

	1. 	M = D(K, EM)
		where
			D - decryption function
			M - map/storage
	2. 	Km = KDF(Ki, S)
		where
			Km - key map/storage
			Ki - input key
	3. 	M = DEL(H(Km))
		where
			DEL - delete Vm from map/storage by Km
	4. 	EM = E(K, M)
*/
package storage
