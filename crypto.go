package gopeer

import (
    "io"
    "math"
    "bytes"
    "crypto"
    "math/big"
    "crypto/aes"
    "crypto/rsa"
    "crypto/rand"
    "crypto/x509"
    "crypto/sha256"
    "crypto/cipher"
    "encoding/pem"
)

// Generate private key and save in object node.
func (node *Node) GeneratePrivate(bits int) *Node {
    priv, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil
    }
    return node.SetPrivate(priv)
}

// Translate public key as *rsa.PublicKey (in node object) to string.
func (node *Node) StringPublic() string {
    if node.Keys.Public == nil {
        return ""
    }
    return StringPublic(node.Keys.Public)
}

// Translate private key as *rsa.PrivateKey (in node object) to string.
func (node *Node) StringPrivate() string {
    if node.Keys.Private == nil {
        return ""
    }
    return StringPrivate(node.Keys.Private)
}

// Translate private key as string to *rsa.PrivateKey (in node object).
func (node *Node) ParsePrivate(privData string) *Node {
    priv := ParsePrivate(privData)
    if priv == nil {
        return nil
    }
    return node.SetPrivate(priv)
}

// Set private key in node object, calculate public and hashname.
func (node *Node) SetPrivate(priv *rsa.PrivateKey) *Node {
    node.Keys.Private = priv
    node.Keys.Public = &priv.PublicKey
    node.Hashname = md5HashName(node.StringPublic())
    return node
}

// Decrypt data by private key in object node.
func (node *Node) DecryptRSA(data []byte) []byte {
    return DecryptRSA(node.Keys.Private, data)
}

// Sign data by private key in object node.
func (node *Node) Sign(data []byte) []byte {
    return Sign(node.Keys.Private, data)
}

// Translate private key as string to *rsa.PrivateKey.
func ParsePrivate(privData string) *rsa.PrivateKey {
    block, _ := pem.Decode([]byte(privData))
    if block == nil {
        return nil
    }
    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil
    }
    return priv
}

// Translate public key as string to *rsa.PublicKey.
func ParsePublic(pubData string) *rsa.PublicKey {
    block, _ := pem.Decode([]byte(pubData))
    if block == nil {
        return nil
    }
    pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
    if err != nil {
        return nil
    }
    return pub
}

// Sign data by private key.
func Sign(priv *rsa.PrivateKey, data []byte) []byte {
    signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, data, nil)
    if err != nil {
        return nil 
    }
    return signature
}

// Verify data and signature by public key.
func Verify(pub *rsa.PublicKey, data, sign []byte) error {
    return rsa.VerifyPSS(pub, crypto.SHA256, data, sign, nil)
}

// SHA256 sum.
// Secure by attacks 'message extension'.
func HashSum(data []byte) []byte {
    if setting.CRYPTO_SPEED {
        // (n/2) bits
        return sumSHA256(sumSHA256(data))
    } else {
        // ~(n) bits
        return sumSHA256(bytes.Join(
            [][]byte{sumSHA256(data), data},
            []byte{},
        ))
    }
}

// Generate integers in range [0:MaxInt64).
func GenerateRandomIntegers(max int) []uint64 {
    var list = make([]uint64, max)
    var maxNum = big.NewInt(math.MaxInt64)
    for i := 0; i < max; i++ {
        nBig, err := rand.Int(rand.Reader, maxNum)
        if err != nil {
            list[i] = 0
            continue
        }
        list[i] = nBig.Uint64()
    }
    return list
}

// Generate bytes in range [33:127).
func GenerateRandomBytes(max int) []byte {
    var slice []byte = make([]byte, max)
    _, err := rand.Read(slice)
    if err != nil { return nil }
    for max = max - 1; max >= 0; max-- {
        slice[max] = slice[max] % 94 + 33
    }
    return slice
}

// Encrypt data by public key.
func EncryptRSA(pub *rsa.PublicKey, data []byte) []byte {
    data, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
    if err != nil {
        return nil 
    }
    return data
}

// Decrypt data by private key.
func DecryptRSA(priv *rsa.PrivateKey, data []byte) []byte {
    data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, data, nil)
    if err != nil {
        return nil 
    }
    return data
}

// Encrypt data by session key.
func EncryptAES(key, data []byte) []byte {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil
    }

    blockSize := block.BlockSize()
    data = paddingPKCS5(data, blockSize)

    cipherText := make([]byte, blockSize + len(data))

    iv := cipherText[:blockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil
    }

    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(cipherText[blockSize:], data)

    return cipherText
}

// Decrypt data by session key.
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

    if len(data) % blockSize != 0 {
        return nil
    }

    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(data, data)

    return unpaddingPKCS5(data)
}

// Translate public key as *rsa.PublicKey to string.
func StringPublic(pub *rsa.PublicKey) string {
    return string(pem.EncodeToMemory(
        &pem.Block{
            Type: "RSA PUBLIC KEY",
            Bytes: x509.MarshalPKCS1PublicKey(pub),
        },
    ))
}

// Translate private key as *rsa.PrivateKey to string.
func StringPrivate(priv *rsa.PrivateKey) string {
    return string(pem.EncodeToMemory(
        &pem.Block{
            Type: "RSA PRIVATE KEY",
            Bytes: x509.MarshalPKCS1PrivateKey(priv),
        },
    ))
}
