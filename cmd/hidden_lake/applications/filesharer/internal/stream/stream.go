package stream

import (
	"context"
	"crypto/sha256"
	"errors"
	"hash"
	"io"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/utils"
	hlf_client "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/client"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/encoding"
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
	chunkSize, err := utils.GetMessageLimit(pCtx, pHlsClient)
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
		chunk, err := p.loadFileChunk(p.fContext)
		if err != nil {
			return 0, err
		}
		if _, err := p.fHasher.Write(chunk); err != nil {
			return 0, err
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
		return 0, errors.New("invalid hash")
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
		return 0, errors.New("stream..Reader.Seek: invalid whence")
	}
	if pos < 0 {
		return 0, errors.New("stream..Reader.Seek: negative position")
	}
	p.fBuffer = p.fBuffer[:0]
	p.fPosition = uint64(pos)
	return pos, nil
}

func (p *sStream) loadFileChunk(pCtx context.Context) ([]byte, error) {
	var lastErr error
	for i := 0; i <= cRetryNum; i++ {
		chunk, err := p.fHlfClient.LoadFileChunk(
			pCtx,
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
	return nil, lastErr
}
