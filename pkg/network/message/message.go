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
	fProof   []byte
	fHash    []byte
	fPayload payload.IPayload
}

func NewMessage(pSett ISettings, pPld payload.IPayload) IMessage {
	hash := getHash(pSett.GetNetworkKey(), pPld.ToBytes())
	proof := puzzle.NewPoWPuzzle(pSett.GetWorkSizeBits()).ProofBytes(hash)
	proofBytes := encoding.Uint64ToBytes(proof)

	return &sMessage{
		fProof:   proofBytes[:],
		fHash:    hash,
		fPayload: pPld,
	}
}

func LoadMessage(pSett ISettings, pData []byte) IMessage {
	if len(pData) < encoding.CSizeUint64+hashing.CSHA256Size {
		return nil
	}

	proofBytes := pData[:encoding.CSizeUint64]
	gotHash := pData[encoding.CSizeUint64 : encoding.CSizeUint64+hashing.CSHA256Size]
	pldBytes := pData[encoding.CSizeUint64+hashing.CSHA256Size:]

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
		fProof:   proofArray[:],
		fHash:    gotHash,
		fPayload: pld,
	}
}

func (p *sMessage) GetProof() uint64 {
	proofArray := [encoding.CSizeUint64]byte{}
	copy(proofArray[:], p.fProof[:])
	return encoding.BytesToUint64(proofArray)
}

func (p *sMessage) GetHash() []byte {
	return p.fHash
}

func (p *sMessage) GetPayload() payload.IPayload {
	return p.fPayload
}

func (p *sMessage) ToBytes() []byte {
	return bytes.Join(
		[][]byte{
			p.fProof,
			p.fHash,
			p.fPayload.ToBytes(),
		},
		[]byte{},
	)
}

func getHash(networkKey string, pBytes []byte) []byte {
	authKey := keybuilder.NewKeyBuilder(
		cWorkSizeKey,
		[]byte(cAuthSalt),
	).Build(networkKey)
	return hashing.NewHMACSHA256Hasher(authKey, pBytes).ToBytes()
}
