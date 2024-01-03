package message

import (
	"bytes"
	"errors"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	// third digits of PI
	cAuthSalt = "8214808651_3282306647_0938446095_5058223172_5359408128"
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fProof   uint64
	fHash    []byte
	fPayload payload.IPayload
}

func NewMessage(pSett ISettings, pPld payload.IPayload, pParallel uint64) IMessage {
	hash := getHash(pSett.GetNetworkKey(), pPld.ToBytes())
	proof := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits()).ProofBytes(hash, pParallel)

	return &sMessage{
		fProof:   proof,
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
		return LoadMessage(pSett, encoding.HexDecode(x))
	default:
		return nil, errors.New("unknown type of message")
	}

	if len(msgBytes) < encoding.CSizeUint64+hashing.CSHA256Size {
		return nil, errors.New("length of message bytes < size of header")
	}

	proofBytes := msgBytes[:encoding.CSizeUint64]
	gotHash := msgBytes[encoding.CSizeUint64 : encoding.CSizeUint64+hashing.CSHA256Size]
	pldBytes := msgBytes[encoding.CSizeUint64+hashing.CSHA256Size:]

	proofArray := [encoding.CSizeUint64]byte{}
	copy(proofArray[:], proofBytes[:])

	proof := encoding.BytesToUint64(proofArray)
	puzzle := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits())
	if !puzzle.VerifyBytes(gotHash, proof) {
		return nil, errors.New("got invalid proof of work")
	}

	if !bytes.Equal(gotHash, getHash(pSett.GetNetworkKey(), pldBytes)) {
		return nil, errors.New("got invalid auth hash")
	}

	// check Head[u64]
	pld := payload.LoadPayload(pldBytes)
	if pld == nil {
		return nil, errors.New("failed to load payload")
	}

	return &sMessage{
		fProof:   proof,
		fHash:    gotHash,
		fPayload: pld,
	}, nil
}

func (p *sMessage) GetProof() uint64 {
	return p.fProof
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
			p.fHash,
			p.fPayload.ToBytes(),
		},
		[]byte{},
	)
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func getHash(networkKey string, pBytes []byte) []byte {
	authKey := keybuilder.NewKeyBuilder(1, []byte(cAuthSalt)).Build(networkKey)
	return hashing.NewHMACSHA256Hasher(authKey, pBytes).ToBytes()
}
