package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

const (
	// IV + Proof + Hash + (2 x Payload32Head) + Payload32Head
	CMessageHeadSize = 0 +
		1*symmetric.CAESBlockSize +
		1*encoding.CSizeUint64 +
		1*hashing.CSHA256Size +
		2*encoding.CSizeUint32 +
		1*encoding.CSizeUint32
)

const (
	cProofIndex = encoding.CSizeUint64
	cHashIndex  = cProofIndex + hashing.CSHA256Size
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fEncd    []byte             // E( K, P(HLMV) || HLMV || L(M) || M || V )
	fHash    []byte             // HLMV = H( K, L(M) || M || V )
	fVoid    []byte             // V
	fProof   uint64             // P(HLMV)
	fPayload payload.IPayload32 // M
}

func NewMessage(pSett IConstructSettings, pPld payload.IPayload32) IMessage {
	prng := random.NewCSPRNG()

	voidBytes := prng.GetBytes(prng.GetUint64() % (pSett.GetLimitVoidSizeBytes() + 1))
	bytesJoiner := joiner.NewBytesJoiner32([][]byte{pPld.ToBytes(), voidBytes})

	key := hashing.NewSHA256Hasher([]byte(pSett.GetNetworkKey())).ToBytes()
	hash := hashing.NewHMACSHA256Hasher(key, bytesJoiner).ToBytes()

	proof := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits()).ProofBytes(hash, pSett.GetParallel())
	proofBytes := encoding.Uint64ToBytes(proof)

	cipher := symmetric.NewAESCipher(key)
	return &sMessage{
		fEncd: cipher.EncryptBytes(bytes.Join(
			[][]byte{
				proofBytes[:],
				hash,
				bytesJoiner,
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

	key := hashing.NewSHA256Hasher([]byte(pSett.GetNetworkKey())).ToBytes()
	dBytes := symmetric.NewAESCipher(key).DecryptBytes(msgBytes)

	proofArr := [encoding.CSizeUint64]byte{}
	copy(proofArr[:], dBytes[:cProofIndex])
	proof := encoding.BytesToUint64(proofArr)

	hash := dBytes[cProofIndex:cHashIndex]
	puzzle := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits())
	if !puzzle.VerifyBytes(hash, proof) {
		return nil, ErrInvalidProofOfWork
	}

	bytesSlice, err := joiner.LoadBytesJoiner32(dBytes[cHashIndex:])
	if err != nil || len(bytesSlice) != 2 {
		return nil, ErrDecodeBytesJoiner
	}

	newHash := hashing.NewHMACSHA256Hasher(key, dBytes[cHashIndex:]).ToBytes()
	if !bytes.Equal(hash, newHash) {
		return nil, ErrInvalidAuthHash
	}

	payload := payload.LoadPayload32(bytesSlice[0])
	if payload == nil {
		return nil, ErrDecodePayload
	}

	return &sMessage{
		fEncd:    msgBytes,
		fHash:    hash,
		fVoid:    bytesSlice[1],
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

func (p *sMessage) GetPayload() payload.IPayload32 {
	return p.fPayload
}

func (p *sMessage) ToBytes() []byte {
	return p.fEncd
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}
