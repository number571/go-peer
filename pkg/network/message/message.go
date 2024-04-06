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
	// Salt(cipher) + Salt(auth) + IV + Hash + Proof + PayloadSize + PayloadHead
	cSaltSize        = 16
	CMessageHeadSize = 2*cSaltSize +
		symmetric.CAESBlockSize +
		hashing.CSHA256Size +
		3*encoding.CSizeUint64
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fSalt    [2][]byte // cipher_salt+auth_salt
	fEncPPHP []byte    // enc(proof_psize_hash_payload)

	// not used in LoadMessage(), ToBytes(), ToString()
	// used only in GetVoid(), GetHash(), in GetProof(), GetPayload()
	fHash    []byte
	fVoid    []byte
	fProof   uint64
	fPayload payload.IPayload
}

func NewMessage(pSett ISettings, pPld payload.IPayload, pParallel, pLimitVoidSize uint64) IMessage {
	payloadBytes := pPld.ToBytes()
	payloadSize := encoding.Uint64ToBytes(uint64(len(payloadBytes)))

	prng := random.NewStdPRNG()

	randVoidSize := prng.GetUint64() % (pLimitVoidSize + 1)
	voidBytes := prng.GetBytes(randVoidSize)

	payloadVoidBytes := bytes.Join(
		[][]byte{
			payloadSize[:],
			payloadBytes,
			voidBytes,
		},
		[]byte{},
	)

	authSalt := prng.GetBytes(cSaltSize)
	hash := getAuthHash(pSett.GetNetworkKey(), authSalt, payloadVoidBytes)
	proof := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits()).ProofBytes(hash, pParallel)

	cipherSalt := prng.GetBytes(cSaltSize)
	cipher := getCipher(pSett.GetNetworkKey(), cipherSalt)
	proofBytes := encoding.Uint64ToBytes(proof)

	return &sMessage{
		fSalt: [2][]byte{
			cipherSalt,
			authSalt,
		},
		fEncPPHP: cipher.EncryptBytes(bytes.Join(
			[][]byte{
				proofBytes[:],
				hash,
				payloadVoidBytes,
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

	const (
		cipherSaltIndex = cSaltSize
		authSaltIndex   = cipherSaltIndex + cSaltSize
	)

	const (
		proofIndex  = encoding.CSizeUint64
		hashIndex   = proofIndex + hashing.CSHA256Size
		pldLenIndex = hashIndex + encoding.CSizeUint64
	)

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

	cipherSaltBytes := msgBytes[:cipherSaltIndex]
	authSaltBytes := msgBytes[cipherSaltIndex:authSaltIndex]
	encPPHPBytes := msgBytes[authSaltIndex:]

	cipher := getCipher(pSett.GetNetworkKey(), cipherSaltBytes)
	pphpBytes := cipher.DecryptBytes(encPPHPBytes)

	proofArr := [encoding.CSizeUint64]byte{}
	copy(proofArr[:], pphpBytes[:proofIndex])
	proof := encoding.BytesToUint64(proofArr)

	payloadSizeArr := [encoding.CSizeUint64]byte{}
	copy(payloadSizeArr[:], pphpBytes[hashIndex:pldLenIndex])
	payloadLength := encoding.BytesToUint64(payloadSizeArr)

	hash := pphpBytes[proofIndex:hashIndex]
	puzzle := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits())
	if !puzzle.VerifyBytes(hash, proof) {
		return nil, ErrInvalidProofOfWork
	}

	payloadVoidBytes := pphpBytes[hashIndex:]
	if (payloadLength + encoding.CSizeUint64) > uint64(len(payloadVoidBytes)) {
		return nil, ErrInvalidPayloadSize
	}

	newHash := getAuthHash(pSett.GetNetworkKey(), authSaltBytes, payloadVoidBytes)
	if !bytes.Equal(hash, newHash) {
		return nil, ErrInvalidAuthHash
	}

	payloadBytes := payloadVoidBytes[encoding.CSizeUint64:][:payloadLength]
	payload := payload.LoadPayload(payloadBytes)
	if payload == nil {
		return nil, ErrDecodePayload
	}

	return &sMessage{
		fSalt: [2][]byte{
			cipherSaltBytes,
			authSaltBytes,
		},
		fEncPPHP: encPPHPBytes,
		fHash:    hash,
		fVoid:    payloadVoidBytes[payloadLength:],
		fProof:   proof,
		fPayload: payload,
	}, nil
}

func (p *sMessage) GetProof() uint64 {
	return p.fProof
}

func (p *sMessage) GetSalt() [2][]byte {
	return p.fSalt
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
		[][]byte{
			p.fSalt[0],
			p.fSalt[1],
			p.fEncPPHP,
		},
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
