package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	// IV + Proof + Hash + PayloadSize + PayloadHead
	CMessageHeadSize = 0 +
		symmetric.CAESBlockSize +
		encoding.CSizeUint64 +
		hashing.CSHA256Size +
		encoding.CSizeUint64 +
		encoding.CSizeUint64
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
	fEncd    []byte           // E( K, P(HLMV) || HLMV || L(M) || M || V )
	fHash    []byte           // HLMV = H( K, L(M) || M || V )
	fVoid    []byte           // V
	fProof   uint64           // P(HLMV)
	fPayload payload.IPayload // M
}

func NewMessage(pSett ISettings, pPld payload.IPayload, pParallel, pLimitVoidSize uint64) IMessage {
	prng := random.NewStdPRNG()

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

	key := hashing.NewSHA256Hasher([]byte(pSett.GetNetworkKey())).ToBytes()
	hash := hashing.NewHMACSHA256Hasher(key, sizeXPayloadVoidBytes).ToBytes()
	cipher := symmetric.NewAESCipher(key)

	proof := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits()).ProofBytes(hash, pParallel)
	proofBytes := encoding.Uint64ToBytes(proof)

	return &sMessage{
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

	key := hashing.NewSHA256Hasher([]byte(pSett.GetNetworkKey())).ToBytes()
	dBytes := symmetric.NewAESCipher(key).DecryptBytes(msgBytes)

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

	newHash := hashing.NewHMACSHA256Hasher(key, dBytes[cHashIndex:]).ToBytes()
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
		fEncd:    msgBytes,
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
	return p.fEncd
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}
