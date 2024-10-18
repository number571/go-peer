// Package crypto represents wrapper functions and interfaces over cryptographic primitives
// with hardcode 192-bit security.
//
// A package consists of many sub-packages:
// 1. quantum - asymmetric KEM, message signing,
// 2. hashing - message hashing and message authentication code,
// 3. keybuilder - generating a key from a password,
// 4. puzzle - solving complex mathematical problems,
// 5. random - generation of random bytes, strings, numbers,
// 6. symmetric - symmetric message encryption.
package crypto
