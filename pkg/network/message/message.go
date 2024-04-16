package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	// Salt = 128bit
	CSaltSize = 16

	// Salt(cipher) + Salt(auth) + IV + Hash + Proof + PayloadSize + PayloadHead
	CMessageHeadSize = 0 +
		CSaltSize +
		CSaltSize +
		symmetric.CAESBlockSize +
		hashing.CSHA256Size +
		encoding.CSizeUint64 +
		encoding.CSizeUint64 +
		encoding.CSizeUint64
)

const (
	cCipherSaltIndex = CSaltSize
	cAuthSaltIndex   = cCipherSaltIndex + CSaltSize
)

const (
	cProofIndex  = encoding.CSizeUint64
	cHashIndex   = cProofIndex + hashing.CSHA256Size
	cPldLenIndex = cHashIndex + encoding.CSizeUint64
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	// Result fields
	fSalt []byte // Sc || Sa
	fEncd []byte // E( KDF(K,Sc), P(HLMV) || HLMV || LM || M || V )

	// Only read fields
	fHash    []byte           // HLMV = H( KDF(K,Sa), LM || M || V )
	fVoid    []byte           // V
	fProof   uint64           // P(HLMV)
	fPayload payload.IPayload // M
}

func NewMessage(pSett ISettings, pPld payload.IPayload, pParallel, pLimitVoidSize uint64) IMessage {
	prng := random.NewStdPRNG()

	authSalt := prng.GetBytes(CSaltSize)
	cipherSalt := prng.GetBytes(CSaltSize)

	payloadBytes := pPld.ToBytes()
	payloadSize := encoding.Uint64ToBytes(uint64(len(payloadBytes)))

	voidBytes := prng.GetBytes(prng.GetUint64() % (pLimitVoidSize + 1))
	sizeXPayloadVoidBytes := bytes.Join(
		[][]byte{
			payloadSize[:],
			payloadBytes,
			voidBytes,
		},
		[]byte{},
	)

	hash := getAuthHash(pSett.GetNetworkKey(), authSalt, sizeXPayloadVoidBytes)
	proof := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits()).ProofBytes(hash, pParallel)
	proofBytes := encoding.Uint64ToBytes(proof)

	cipher := getCipher(pSett.GetNetworkKey(), cipherSalt)
	return &sMessage{
		fSalt: bytes.Join(
			[][]byte{cipherSalt, authSalt},
			[]byte{},
		),
		fEncd: cipher.EncryptBytes(bytes.Join(
			[][]byte{
				proofBytes[:],
				hash,
				sizeXPayloadVoidBytes,
			},
			[]byte{},
		)),
		fHash:    hash,
		fVoid:    voidBytes,
		fProof:   proof,
		fPayload: pPld,
	}
}

func LoadMessage(pSett ISettings, pData interface{}) (IMessage, error) {
	var msgBytes []byte

	switch x := pData.(type) {
	case []byte:
		msgBytes = x
	case string:
		msgBytes = encoding.HexDecode(x)
	default:
		return nil, ErrUnknownType
	}

	if len(msgBytes) < CMessageHeadSize {
		return nil, ErrInvalidHeaderSize
	}

	saltBytes := msgBytes[:cAuthSaltIndex]
	encdBytes := msgBytes[cAuthSaltIndex:]

	cipherSaltBytes := saltBytes[:cCipherSaltIndex]
	authSaltBytes := saltBytes[cCipherSaltIndex:]

	cipher := getCipher(pSett.GetNetworkKey(), cipherSaltBytes)
	dBytes := cipher.DecryptBytes(encdBytes)

	proofArr := [encoding.CSizeUint64]byte{}
	copy(proofArr[:], dBytes[:cProofIndex])
	proof := encoding.BytesToUint64(proofArr)

	payloadSizeArr := [encoding.CSizeUint64]byte{}
	copy(payloadSizeArr[:], dBytes[cHashIndex:cPldLenIndex])
	payloadLength := encoding.BytesToUint64(payloadSizeArr)

	hash := dBytes[cProofIndex:cHashIndex]
	puzzle := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits())
	if !puzzle.VerifyBytes(hash, proof) {
		return nil, ErrInvalidProofOfWork
	}

	payloadVoidBytes := dBytes[cPldLenIndex:]
	if payloadLength > uint64(len(payloadVoidBytes)) {
		return nil, ErrInvalidPayloadSize
	}

	newHash := getAuthHash(pSett.GetNetworkKey(), authSaltBytes, dBytes[cHashIndex:])
	if !bytes.Equal(hash, newHash) {
		return nil, ErrInvalidAuthHash
	}

	payloadBytes := payloadVoidBytes[:payloadLength]
	voidBytes := payloadVoidBytes[payloadLength:]

	payload := payload.LoadPayload(payloadBytes)
	if payload == nil {
		return nil, ErrDecodePayload
	}

	return &sMessage{
		fSalt:    saltBytes,
		fEncd:    encdBytes,
		fHash:    hash,
		fVoid:    voidBytes,
		fProof:   proof,
		fPayload: payload,
	}, nil
}

func (p *sMessage) GetProof() uint64 {
	return p.fProof
}

func (p *sMessage) GetHash() []byte {
	return p.fHash
}

func (p *sMessage) GetVoid() []byte {
	return p.fVoid
}

func (p *sMessage) GetPayload() payload.IPayload {
	return p.fPayload
}

func (p *sMessage) ToBytes() []byte {
	return bytes.Join(
		[][]byte{p.fSalt, p.fEncd},
		[]byte{},
	)
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func getAuthHash(networkKey string, pAuthSalt, pBytes []byte) []byte {
	authKey := keybuilder.NewKeyBuilder(1, pAuthSalt).Build(networkKey)
	return hashing.NewHMACSHA256Hasher(authKey, pBytes).ToBytes()
}

func getCipher(networkKey string, pCipherSalt []byte) symmetric.ICipher {
	cipherKey := keybuilder.NewKeyBuilder(1, pCipherSalt).Build(networkKey)
	return symmetric.NewAESCipher(cipherKey)
}
