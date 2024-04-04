package stream

import (
	"context"
	"crypto/sha256"
	"hash"
	"io"

	internal_utils "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/utils"
	hlf_client "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/client"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
)

const (
	cRetryNum = 2
)

var (
	_ IReadSeeker = &sStream{}
)

type sStream struct {
	fContext   context.Context
	fHlfClient hlf_client.IClient

	fBuffer   []byte
	fPosition uint64

	fHasher   hash.Hash
	fFileHash string

	fAliasName string
	fFileName  string

	fFileSize  uint64
	fChunkSize uint64
}

func BuildStream(
	pCtx context.Context,
	pHlsClient hls_client.IClient,
	pAliasName string,
	pFileName string,
	pFileHash string,
	pFileSize uint64,
) IReadSeeker {
	chunkSize, err := internal_utils.GetMessageLimit(pCtx, pHlsClient)
	if err != nil {
		return nil
	}

	return &sStream{
		fContext: pCtx,

		fHlfClient: hlf_client.NewClient(
			hlf_client.NewBuilder(),
			hlf_client.NewRequester(pHlsClient),
		),

		fAliasName: pAliasName,
		fFileName:  pFileName,

		fHasher:   sha256.New(),
		fFileHash: pFileHash,

		fFileSize:  pFileSize,
		fChunkSize: chunkSize,
	}
}

func (p *sStream) Read(b []byte) (int, error) {
	if len(p.fBuffer) == 0 {
		chunk, err := p.loadFileChunk()
		if err != nil {
			return 0, utils.MergeErrors(ErrLoadFileChunk, err)
		}
		if _, err := p.fHasher.Write(chunk); err != nil {
			return 0, utils.MergeErrors(ErrWriteFileChunk, err)
		}
		p.fBuffer = chunk
	}

	n := copy(b, p.fBuffer)
	p.fBuffer = p.fBuffer[n:]
	p.fPosition += uint64(n)

	if p.fPosition < p.fFileSize {
		return n, nil
	}

	hashSum := encoding.HexEncode(p.fHasher.Sum(nil))
	if hashSum != p.fFileHash {
		return 0, ErrInvalidHash
	}

	return n, io.EOF
}

func (p *sStream) Seek(offset int64, whence int) (int64, error) {
	var pos int64
	switch whence {
	case io.SeekStart:
		pos = offset
	case io.SeekCurrent:
		pos = int64(p.fPosition) + offset
	case io.SeekEnd:
		pos = int64(p.fFileSize) + offset
	default:
		return 0, ErrInvalidWhence
	}
	if pos < 0 {
		return 0, ErrNegativePosition
	}
	p.fBuffer = p.fBuffer[:0]
	p.fPosition = uint64(pos)
	return pos, nil
}

func (p *sStream) loadFileChunk() ([]byte, error) {
	var lastErr error
	for i := 0; i <= cRetryNum; i++ {
		chunk, err := p.fHlfClient.LoadFileChunk(
			p.fContext,
			p.fAliasName,
			p.fFileName,
			p.fPosition/p.fChunkSize,
		)
		if err != nil {
			lastErr = err
			continue
		}
		return chunk, nil
	}
	return nil, utils.MergeErrors(ErrRetryFailed, lastErr)
}
