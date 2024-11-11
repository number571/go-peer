package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	// IV + Proof + Hash + Payload32Head=[head]
	// 16 + 8 + 48 + 4 = 76 additional bytes to origin message
	CMessageHeadSize = 0 +
		1*symmetric.CCipherBlockSize +
		1*encoding.CSizeUint64 +
		1*hashing.CHasherSize +
		1*encoding.CSizeUint32
)

const (
	cProofIndex = encoding.CSizeUint64
	cHashIndex  = cProofIndex + hashing.CHasherSize
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fEncd    []byte             // E( K, P(HM) || HM || M )
	fHash    []byte             // HM = H( K, M )
	fProof   uint64             // P(HM)
	fPayload payload.IPayload32 // M
}

func NewMessage(pSett IConstructSettings, pPld payload.IPayload32) IMessage {
	sett := pSett.GetSettings()
	pldBytes := pPld.ToBytes()

	keyBuilder := keybuilder.NewKeyBuilder(0, []byte{}) // the network_key must have good entropy
	key := keyBuilder.Build(sett.GetNetworkKey(), symmetric.CCipherKeySize)
	hash := hashing.NewHMACHasher(key, pldBytes).ToBytes()

	proof := puzzle.NewPoWPuzzle(sett.GetWorkSizeBits()).ProofBytes(hash, pSett.GetParallel())
	proofBytes := encoding.Uint64ToBytes(proof)

	cipher := symmetric.NewCipher(key)
	return &sMessage{
		fEncd: cipher.EncryptBytes(bytes.Join(
			[][]byte{
				proofBytes[:],
				hash,
				pldBytes,
			},
			[]byte{},
		)),
		fHash:    hash,
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

	keyBuilder := keybuilder.NewKeyBuilder(0, []byte{}) // the network_key must have good entropy
	key := keyBuilder.Build(pSett.GetNetworkKey(), symmetric.CCipherKeySize)
	dBytes := symmetric.NewCipher(key).DecryptBytes(msgBytes)

	proofArr := [encoding.CSizeUint64]byte{}
	copy(proofArr[:], dBytes[:cProofIndex])
	proof := encoding.BytesToUint64(proofArr)

	hash := dBytes[cProofIndex:cHashIndex]
	puzzle := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits())
	if !puzzle.VerifyBytes(hash, proof) {
		return nil, ErrInvalidProofOfWork
	}

	newHash := hashing.NewHMACHasher(key, dBytes[cHashIndex:]).ToBytes()
	if !bytes.Equal(hash, newHash) {
		return nil, ErrInvalidAuthHash
	}

	payload := payload.LoadPayload32(dBytes[cHashIndex:])
	if payload == nil {
		return nil, ErrDecodePayload
	}

	return &sMessage{
		fEncd:    msgBytes,
		fHash:    hash,
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

func (p *sMessage) GetPayload() payload.IPayload32 {
	return p.fPayload
}

func (p *sMessage) ToBytes() []byte {
	return p.fEncd
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}
