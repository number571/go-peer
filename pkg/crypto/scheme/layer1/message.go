package layer1

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	// IV + Proof + HMAC
	// 16 + 8 + 48 = 72 additional bytes to origin message
	CMessageHeadSize = 0 +
		1*symmetric.CCipherBlockSize +
		1*encoding.CSizeUint64 +
		1*hashing.CHasherSize
)

const (
	cProofIndex = encoding.CSizeUint64
	cHashIndex  = cProofIndex + hashing.CHasherSize
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fProof uint64 // P(HM)
	fEncd  []byte // E( K, P(HM) || HM || M )
	fHash  []byte // H( M )
	fHmac  []byte // HM = HMAC( K, M )
	fBody  []byte // M
}

func NewMessage(pSett IConstructSettings, pBody []byte) IMessage {
	sett := pSett.GetSettings()
	hash := hashing.NewHasher(pBody).ToBytes()

	keyBuilder := keybuilder.NewKeyBuilder(0, []byte{}) // the network_key must have good entropy
	key := keyBuilder.Build(sett.GetNetworkKey(), symmetric.CCipherKeySize)
	hmac := hashing.NewHMACHasher(key, pBody).ToBytes()

	proof := puzzle.NewPoWPuzzle(sett.GetWorkSizeBits()).ProofBytes(hmac, pSett.GetParallel())
	proofBytes := encoding.Uint64ToBytes(proof)

	cipher := symmetric.NewCipherCFB(key)
	return &sMessage{
		fEncd: cipher.EncryptBytes(bytes.Join(
			[][]byte{
				proofBytes[:],
				hmac,
				pBody,
			},
			[]byte{},
		)),
		fHash:  hash,
		fHmac:  hmac,
		fProof: proof,
		fBody:  pBody,
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

	keyBuilder := keybuilder.NewKeyBuilder(0, []byte{}) // the network_key must have good entropy
	key := keyBuilder.Build(pSett.GetNetworkKey(), symmetric.CCipherKeySize)
	dBytes := symmetric.NewCipherCFB(key).DecryptBytes(msgBytes)

	proofArr := [encoding.CSizeUint64]byte{}
	copy(proofArr[:], dBytes[:cProofIndex])
	proof := encoding.BytesToUint64(proofArr)

	hmac := dBytes[cProofIndex:cHashIndex]
	puzzle := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits())
	if !puzzle.VerifyBytes(hmac, proof) {
		return nil, ErrInvalidProofOfWork
	}

	body := dBytes[cHashIndex:]
	newHmac := hashing.NewHMACHasher(key, body).ToBytes()
	if !bytes.Equal(hmac, newHmac) {
		return nil, ErrInvalidAuthHash
	}

	hash := hashing.NewHasher(body).ToBytes()
	return &sMessage{
		fEncd:  msgBytes,
		fHash:  hash,
		fHmac:  hmac,
		fProof: proof,
		fBody:  body,
	}, nil
}

func (p *sMessage) GetProof() uint64 {
	return p.fProof
}

func (p *sMessage) GetHash() []byte {
	return p.fHash
}

func (p *sMessage) GetHmac() []byte {
	return p.fHmac
}

func (p *sMessage) GetBody() []byte {
	return p.fBody
}

func (p *sMessage) ToBytes() []byte {
	return p.fEncd
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}
