package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	CSaltSize = 32
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fProof   uint64
	fSalt    []byte
	fHash    []byte
	fPayload payload.IPayload
}

func NewMessage(pSett ISettings, pPld payload.IPayload, pParallel uint64) IMessage {
	salt := random.NewStdPRNG().GetBytes(CSaltSize)
	hash := getAuthHash(pSett.GetNetworkKey(), salt, pPld.ToBytes())
	proof := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits()).ProofBytes(hash, pParallel)

	return &sMessage{
		fProof:   proof,
		fSalt:    salt,
		fHash:    hash,
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

	if len(msgBytes) < encoding.CSizeUint64+CSaltSize+hashing.CSHA256Size {
		return nil, ErrInvalidHeaderSize
	}

	proofBytes := msgBytes[:encoding.CSizeUint64]
	gotSalt := msgBytes[encoding.CSizeUint64 : encoding.CSizeUint64+CSaltSize]
	gotHash := msgBytes[encoding.CSizeUint64+CSaltSize : encoding.CSizeUint64+CSaltSize+hashing.CSHA256Size]
	pldBytes := msgBytes[encoding.CSizeUint64+CSaltSize+hashing.CSHA256Size:]

	proofArray := [encoding.CSizeUint64]byte{}
	copy(proofArray[:], proofBytes)

	proof := encoding.BytesToUint64(proofArray)
	puzzle := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits())
	if !puzzle.VerifyBytes(gotHash, proof) {
		return nil, ErrInvalidProofOfWork
	}

	if !bytes.Equal(gotHash, getAuthHash(pSett.GetNetworkKey(), gotSalt, pldBytes)) {
		return nil, ErrInvalidAuthHash
	}

	// check Head[u64]
	pld := payload.LoadPayload(pldBytes)
	if pld == nil {
		return nil, ErrDecodePayload
	}

	return &sMessage{
		fProof:   proof,
		fSalt:    gotSalt,
		fHash:    gotHash,
		fPayload: pld,
	}, nil
}

func (p *sMessage) GetProof() uint64 {
	return p.fProof
}

func (p *sMessage) GetSalt() []byte {
	return p.fSalt
}

func (p *sMessage) GetHash() []byte {
	return p.fHash
}

func (p *sMessage) GetPayload() payload.IPayload {
	return p.fPayload
}

func (p *sMessage) ToBytes() []byte {
	proofBytes := encoding.Uint64ToBytes(p.fProof)
	return bytes.Join(
		[][]byte{
			proofBytes[:],
			p.fSalt,
			p.fHash,
			p.fPayload.ToBytes(),
		},
		[]byte{},
	)
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func getAuthHash(networkKey string, pAuthSalt, pBytes []byte) []byte {
	authKey := keybuilder.NewKeyBuilder(1, []byte(pAuthSalt)).Build(networkKey)
	return hashing.NewHMACSHA256Hasher(authKey, pBytes).ToBytes()
}
