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
    "crypto/x509/pkix"
    "encoding/base64"
    "encoding/binary"
    "encoding/json"
    "encoding/pem"
    "io"
    "os"
    "math"
    "math/big"
    "time"
)

// Encrypt file with algorithm AES256-OFB.
func FileEncryptAES(key []byte, input string, output string) error { 
    const BEGGINING = 0 
    in, err := os.Open(input)  
    if err != nil {
        return err
    }  
    defer in.Close()
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    blockSize := block.BlockSize()
    iv := make([]byte, blockSize)

    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return err
    }

    stream := cipher.NewOFB(block, iv)  
    out, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        return err
    }  
    defer out.Close()
  
    out.Write(iv)
    writer := &cipher.StreamWriter{S: stream, W: out}
    if _, err := io.Copy(writer, in); err != nil {  
        return err
    }  

    return nil
}  

// Decrypt file with algorithm AES256-OFB.
func FileDecryptAES(key []byte, input string, output string) error {  
    const BEGGINING = 0
    in, err := os.Open(input)
    if err != nil {
        return err
    }
    defer in.Close()  
    block, err := aes.NewCipher([]byte(key))  
    if err != nil {
        return err
    }

    blockSize := block.BlockSize()
    iv := readFileBytes(input, blockSize)

    stream := cipher.NewOFB(block, iv)  
    out, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        return err
    }  
    defer out.Close()

    _, err = in.Seek(int64(blockSize), BEGGINING)
    if err != nil {
        return nil
    }

    reader := &cipher.StreamReader{S: stream, R: in}
    if _, err := io.Copy(out, reader); err != nil {
        return err
    }

    return nil
}

func readFileBytes(input string, max int) []byte {
    file, err := os.Open(input)
    if err != nil {
        return nil
    }
    defer file.Close()
    var buffer = make([]byte, max)
    length, err := file.Read(buffer)
    if err != nil {
        return nil
    }
    return buffer[:length]
}

// Generate certificate by server name and number bits of private key.
func GenerateCertificate(name string, bits uint16) (string, string) {
    ca := &x509.Certificate{
        SerialNumber: big.NewInt(int64(GenerateRandomIntegers(1)[0])),
        Subject: pkix.Name{
            CommonName:    name,
            Organization:  []string{name},
            Country:       []string{name},
            Province:      []string{name},
            Locality:      []string{name},
            StreetAddress: []string{name},
            PostalCode:    []string{name},
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years
        IsCA:                  true,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
        KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
        BasicConstraintsValid: true,
    }
    priv := GeneratePrivate(bits)
    cert, err := x509.CreateCertificate(rand.Reader, ca, ca, &priv.PublicKey, priv)
    if err != nil {
        return "", ""
    }
    return StringPrivate(priv), StringCertificate(cert)
}

// Create private key by size bits.
func GeneratePrivate(bits uint16) *rsa.PrivateKey {
    priv, err := rsa.GenerateKey(rand.Reader, int(bits))
    if err != nil {
        return nil
    }
    return priv
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

// Translate certificate as string to *x509.Certificate.
func ParseCertificate(certData string) *x509.Certificate {
    block, _ := pem.Decode([]byte(certData))
    if block == nil {
        return nil
    }
    cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        return nil
    }
    return cert
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

// Hash string by public key.
func HashPublic(pub *rsa.PublicKey) string {
    return Base64Encode(HashSum([]byte(StringPublic(pub))))
}

// Hash sum by HMAC(SHA256, HMACKEY).
func HashSum(data []byte) []byte {
    return HMAC(func(data []byte) []byte {
        hash := sha256.Sum256(data)
        return hash[:]
    }, data, []byte(settings.HMAC_KEY))
}

// MAC by cryptographic hash function.
func HMAC(fHash func([]byte) []byte, data []byte, key []byte) []byte {
    const (
        a = 0x5c
        b = 0x36
    )
    var (
        length = len(key)
        outer  = make([]byte, length)
        inner  = make([]byte, length)
    )
    for index, byte := range key {
        outer[index] = byte ^ a
        inner[index] = byte ^ b
    }
    return fHash(bytes.Join(
        [][]byte{outer, fHash(bytes.Join(
            [][]byte{inner, data},
            []byte{},
        ))},
        []byte{},
    ))
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

// Generate bytes in range [0:256).
func GenerateRandomBytes(max int) []byte {
    var slice []byte = make([]byte, max)
    _, err := rand.Read(slice)
    if err != nil {
        return nil
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

    cipherText := make([]byte, blockSize+len(data))

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

    if len(data)%blockSize != 0 {
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
            Type:  "RSA PUBLIC KEY",
            Bytes: x509.MarshalPKCS1PublicKey(pub),
        },
    ))
}

// Translate private key as *rsa.PrivateKey to string.
func StringPrivate(priv *rsa.PrivateKey) string {
    return string(pem.EncodeToMemory(
        &pem.Block{
            Type:  "RSA PRIVATE KEY",
            Bytes: x509.MarshalPKCS1PrivateKey(priv),
        },
    ))
}

// Translate certificate bytes to string.
func StringCertificate(cert []byte) string {
    return string(pem.EncodeToMemory(
        &pem.Block{
            Type:  "CERTIFICATE",
            Bytes: cert,
        },
    ))
}

// base64.StdEncoding.EncodeToString
func Base64Encode(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// base64.StdEncoding.DecodeString
func Base64Decode(data string) []byte {
    result, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        return nil
    }
    return result
}

// Pack another types data to JSON.
func PackJSON(data interface{}) []byte {
    jsonData, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return nil
    }
    return jsonData
}

// Unpack JSON to another types data.
func UnpackJSON(jsonData []byte, data interface{}) interface{} {
    err := json.Unmarshal(jsonData, data)
    if err != nil {
        return nil
    }
    return data
}

// POW for check hash package by Nonce.
func ProofOfWork(blockHash []byte, difficulty uint8) uint64 {
    var (
        Target  = big.NewInt(1)
        intHash = big.NewInt(1)
        nonce   uint64
        hash    []byte
    )
    Target.Lsh(Target, 256-uint(difficulty))
    for nonce < math.MaxUint64 {
        hash = HashSum(bytes.Join([][]byte{
            ToBytes(nonce),
            blockHash,
        },
            []byte{},
        ))
        intHash.SetBytes(hash)
        if intHash.Cmp(Target) == -1 {
            break
        }
        nonce++
    }
    return nonce
}

// Return true if Nonce package equal POW(hash, DIFF).
func NonceIsValid(blockHash []byte, difficulty uint, nonce uint64) bool {
    var (
        Target  = big.NewInt(1)
        intHash = big.NewInt(1)
        hash    []byte
    )
    Target.Lsh(Target, 256-difficulty)
    hash = HashSum(bytes.Join([][]byte{
        ToBytes(nonce),
        blockHash,
    },
        []byte{},
    ))
    intHash.SetBytes(hash)
    if intHash.Cmp(Target) == -1 {
        return true
    }
    return false
}

// Translate uint64 to slice of bytes.
func ToBytes(num uint64) []byte {
    var data = new(bytes.Buffer)
    err := binary.Write(data, binary.BigEndian, num)
    if err != nil {
        return nil
    }
    return data.Bytes()
}
