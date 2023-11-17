package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	cWorkSizeKey = 1

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

func NewMessage(pSett ISettings, pPld payload.IPayload) IMessage {
	hash := getHash(pSett.GetNetworkKey(), pPld.ToBytes())
	proof := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits()).ProofBytes(hash)

	return &sMessage{
		fProof:   proof,
		fHash:    hash,
		fPayload: pPld,
	}
}

func LoadMessage(pSett ISettings, pData interface{}) IMessage {
	var msgBytes []byte

	switch x := pData.(type) {
	case []byte:
		msgBytes = x
	case string:
		return LoadMessage(pSett, encoding.HexDecode(x))
	default:
		return nil
	}

	if len(msgBytes) < encoding.CSizeUint64+hashing.CSHA256Size {
		return nil
	}

	proofBytes := msgBytes[:encoding.CSizeUint64]
	gotHash := msgBytes[encoding.CSizeUint64 : encoding.CSizeUint64+hashing.CSHA256Size]
	pldBytes := msgBytes[encoding.CSizeUint64+hashing.CSHA256Size:]

	proofArray := [encoding.CSizeUint64]byte{}
	copy(proofArray[:], proofBytes[:])

	proof := encoding.BytesToUint64(proofArray)
	puzzle := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits())
	if !puzzle.VerifyBytes(gotHash, proof) {
		return nil
	}

	if !bytes.Equal(gotHash, getHash(pSett.GetNetworkKey(), pldBytes)) {
		return nil
	}

	// check Head[u64]
	pld := payload.LoadPayload(pldBytes)
	if pld == nil {
		return nil
	}

	return &sMessage{
		fProof:   proof,
		fHash:    gotHash,
		fPayload: pld,
	}
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
	authKey := keybuilder.NewKeyBuilder(
		cWorkSizeKey,
		[]byte(cAuthSalt),
	).Build(networkKey)
	return hashing.NewHMACSHA256Hasher(authKey, pBytes).ToBytes()
}
