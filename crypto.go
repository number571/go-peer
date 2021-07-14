package gopeer

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"math"
	"math/big"
)

// Generates a cryptographically strong pseudo-random sequence.
func GenerateBytes(max uint) []byte {
	var slice []byte = make([]byte, max)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}

// Create private key by number of bits.
func GenerateKey(bits uint) *rsa.PrivateKey {
	priv, err := rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return nil
	}
	return priv
}

// Used SHA256.
func HashSum(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// HashPublicKey(x) = Base64Encode(HashSum(PublicKeyToBytes(x))).
func HashPublicKey(pub *rsa.PublicKey) string {
	return Base64Encode(HashSum(PublicKeyToBytes(pub)))
}

// StringToPrivateKey(x) = BytesToPrivateKey(Base64Decode(x)).
func StringToPrivateKey(privData string) *rsa.PrivateKey {
	return BytesToPrivateKey(Base64Decode(privData))
}

// StringToPublicKey(x) = BytesToPublicKey(Base64Decode(x)).
func StringToPublicKey(pubData string) *rsa.PublicKey {
	return BytesToPublicKey(Base64Decode(pubData))
}

// PrivateKeyToString(x) = Base64Encode(PrivateKeyToBytes(x)).
func PrivateKeyToString(priv *rsa.PrivateKey) string {
	return Base64Encode(PrivateKeyToBytes(priv))
}

// PublicKeyToString(x) = Base64Encode(PublicKeyToBytes(pub)).
func PublicKeyToString(pub *rsa.PublicKey) string {
	return Base64Encode(PublicKeyToBytes(pub))
}

// Used PKCS1.
func BytesToPrivateKey(privData []byte) *rsa.PrivateKey {
	priv, err := x509.ParsePKCS1PrivateKey(privData)
	if err != nil {
		return nil
	}
	return priv
}

// Used PKCS1.
func BytesToPublicKey(pubData []byte) *rsa.PublicKey {
	pub, err := x509.ParsePKCS1PublicKey(pubData)
	if err != nil {
		return nil
	}
	return pub
}

// Used PKCS1.
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(priv)
}

// Used PKCS1.
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(pub)
}

// Used RSA(OAEP).
func EncryptRSA(pub *rsa.PublicKey, data []byte) []byte {
	data, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
	if err != nil {
		return nil
	}
	return data
}

// Used RSA(OAEP).
func DecryptRSA(priv *rsa.PrivateKey, data []byte) []byte {
	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, data, nil)
	if err != nil {
		return nil
	}
	return data
}

// Used RSA(PSS).
func Sign(priv *rsa.PrivateKey, data []byte) []byte {
	signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, data, nil)
	if err != nil {
		return nil
	}
	return signature
}

// Used RSA(PSS).
func Verify(pub *rsa.PublicKey, data, sign []byte) error {
	return rsa.VerifyPSS(pub, crypto.SHA256, data, sign, nil)
}

// AES with CBC-mode.
func EncryptAES(key, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	blockSize := block.BlockSize()
	data = paddingPKCS5(data, blockSize)
	cipherText := make([]byte, blockSize+len(data))
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], data)
	return cipherText
}

// AES with CBC-mode.
func DecryptAES(key, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	blockSize := block.BlockSize()
	if len(data) < blockSize {
		return nil
	}
	iv := data[:blockSize]
	data = data[blockSize:]
	if len(data)%blockSize != 0 {
		return nil
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	return unpaddingPKCS5(data)
}

// Increase entropy by multiple hashing.
func RaiseEntropy(info, salt []byte, bits int) []byte {
	lim := uint64(1 << uint(bits))
	for i := uint64(0); i < lim; i++ {
		info = HashSum(bytes.Join(
			[][]byte{
				info,
				salt,
			},
			[]byte{},
		))
	}
	return info
}

// Proof of work by the method of finding the desired hash.
// Hash must start with 'diff' number of zero bits.
func ProofOfWork(packHash []byte, diff uint) uint64 {
	var (
		Target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = uint64(0)
		hash    []byte
	)
	Target.Lsh(Target, 256-diff)
	for nonce < math.MaxUint64 {
		hash = HashSum(bytes.Join(
			[][]byte{
				packHash,
				ToBytes(nonce),
			},
			[]byte{},
		))
		intHash.SetBytes(hash)
		if intHash.Cmp(Target) == -1 {
			return nonce
		}
		nonce++
	}
	return nonce
}

// Verifies the work of the proof of work function.
func ProofIsValid(packHash []byte, diff uint, nonce uint64) bool {
	intHash := big.NewInt(1)
	Target := big.NewInt(1)
	hash := HashSum(bytes.Join(
		[][]byte{
			packHash,
			ToBytes(nonce),
		},
		[]byte{},
	))
	intHash.SetBytes(hash)
	Target.Lsh(Target, 256-diff)
	return intHash.Cmp(Target) == -1
}

// Uint64 to slice of bytes by big endian.
func ToBytes(num uint64) []byte {
	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		return nil
	}
	return data.Bytes()
}

// Serialize with JSON format.
func SerializePackage(pack *Package) []byte {
	jsonData, err := json.MarshalIndent(pack, "", "\t")
	if err != nil {
		return nil
	}
	return jsonData
}

// Deserialize with JSON format.
func DeserializePackage(jsonData []byte) *Package {
	var pack = new(Package)
	err := json.Unmarshal(jsonData, pack)
	if err != nil {
		return nil
	}
	return pack
}

// Standart encoding in package.
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Standart decoding in package.
func Base64Decode(data string) []byte {
	result, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return result
}

func paddingPKCS5(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func unpaddingPKCS5(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	}
	unpadding := int(origData[length-1])
	if length < unpadding {
		return nil
	}
	return origData[:(length - unpadding)]
}

// func printJSON(data interface{}) {
// 	jsonData, _ := json.MarshalIndent(data, "", "\t")
// 	fmt.Println(string(jsonData))
// }
