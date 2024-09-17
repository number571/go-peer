package message

import (
	"bytes"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

const (
	// IV + Proof + Hash + Timestamp + (2xPayload32Head)=[len(data)+len(void)] + Payload32Head=[head]
	// 16 + 8 + 32 + 8 + 2x4 + 4 = 76 additional bytes to origin message
	CMessageHeadSize = 0 +
		1*symmetric.CAESBlockSize +
		1*encoding.CSizeUint64 +
		1*hashing.CSHA256Size +
		1*encoding.CSizeUint64 +
		2*encoding.CSizeUint32 +
		1*encoding.CSizeUint32
)

const (
	cProofIndex     = encoding.CSizeUint64
	cHashIndex      = cProofIndex + hashing.CSHA256Size
	cTimestampIndex = cHashIndex + encoding.CSizeUint64
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fEncd    []byte             // E( K, P(HTLMR) || HTLMR || T || L(M) || M || L(R) || R )
	fHash    []byte             // HTLMR = H( K, T || L(M) || M || L(R) || R )
	fRand    []byte             // R
	fTime    uint64             // T
	fProof   uint64             // P(HTLMR)
	fPayload payload.IPayload32 // M
}

func NewMessage(pSett IConstructSettings, pPld payload.IPayload32) IMessage {
	prng := random.NewCSPRNG()

	timestamp := uint64(time.Now().UTC().Unix())
	randBytes := prng.GetBytes(prng.GetUint64() % (pSett.GetRandMessageSizeBytes() + 1))

	bytesJoiner := joiner.NewBytesJoiner32([][]byte{pPld.ToBytes(), randBytes})
	bytesJoiner = bytes.Join(
		[][]byte{
			func() []byte {
				bytes := encoding.Uint64ToBytes(timestamp)
				return bytes[:]
			}(),
			bytesJoiner,
		},
		[]byte{},
	)

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
		fRand:    randBytes,
		fTime:    timestamp,
		fProof:   proof,
		fPayload: pPld,
	}
}

func LoadMessage(pSett ISettings, pTSWindow time.Duration, pData interface{}) (IMessage, error) {
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

	timestamp := encoding.BytesToUint64(func() [encoding.CSizeUint64]byte {
		arr := [encoding.CSizeUint64]byte{}
		copy(arr[:], dBytes[cHashIndex:cTimestampIndex])
		return arr
	}())

	if pTSWindow != 0 {
		gotTimestamp := time.Unix(int64(timestamp), 0)
		switch nowTimestamp := time.Now().UTC(); {
		case gotTimestamp.After(nowTimestamp.Add(pTSWindow)):
			fallthrough
		case gotTimestamp.Before(nowTimestamp.Add(-pTSWindow)):
			return nil, ErrInvalidTimestamp
		}
	}

	bytesSlice, err := joiner.LoadBytesJoiner32(dBytes[cTimestampIndex:])
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
		fRand:    bytesSlice[1],
		fTime:    timestamp,
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

func (p *sMessage) GetTime() uint64 {
	return p.fTime
}

func (p *sMessage) GetRand() []byte {
	return p.fRand
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
