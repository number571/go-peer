package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"hash"
	"os"

	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func main() {
	saltParam := flag.String("salt", "_salt_", "default salt value")
	workParam := flag.Uint("work", 24, "default work value")
	flag.Parse()

	fmt.Println(base64.URLEncoding.EncodeToString(Key(
		[]byte(readUntilEOL()),
		[]byte(*saltParam),
		1<<(*workParam),
		symmetric.CAESKeySize,
		sha256.New,
	)))
}

func readUntilEOL() string {
	res, _, err := bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		panic(err)
	}
	return string(res)
}

// FROM: https://cs.opensource.google/go/x/crypto/+/refs/tags/v0.25.0:pbkdf2/pbkdf2.go;l=42
func Key(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}
