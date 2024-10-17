package conn

import (
	"bytes"
	"context"
	"errors"
	"math"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcHead = 12345
	tcBody = "hello, world!"
)

type tsConn struct {
	readDlError bool
	cancelBody  bool
	bodyPart    bool
	headSize    uint32
	bodySize    uint64
}

func (p *tsConn) Read(b []byte) (n int, err error) {
	if !p.bodyPart {
		headBytes := encoding.Uint32ToBytes(p.headSize)
		n = copy(b, headBytes[:])
		p.bodyPart = true
		return n, nil
	}
	if p.cancelBody {
		return 0, errors.New("some error1") // nolint: goerr113
	}
	bodyBytes := random.NewRandom().GetBytes(p.bodySize)
	n = copy(b, bodyBytes)
	return n, nil
}
func (p *tsConn) Write(_ []byte) (n int, err error) {
	return 0, errors.New("some error2") // nolint: goerr113
}
func (p *tsConn) Close() error                  { return nil }
func (p *tsConn) LocalAddr() net.Addr           { return &net.TCPAddr{} }
func (p *tsConn) RemoteAddr() net.Addr          { return &net.TCPAddr{} }
func (p *tsConn) SetDeadline(_ time.Time) error { return nil }
func (p *tsConn) SetReadDeadline(_ time.Time) error {
	if p.bodyPart && p.readDlError {
		return errors.New("some error3") // nolint: goerr113
	}
	return nil
}
func (p *tsConn) SetWriteDeadline(_ time.Time) error { return nil }

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SConnError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 6; i++ {
		testSettings(t, i)
	}
}

func testSettings(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = NewSettings(&SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{}),
			FWaitReadTimeout: time.Hour,
			FDialTimeout:     time.Minute,
			FReadTimeout:     time.Minute,
			FWriteTimeout:    time.Minute,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FMessageSettings:       message.NewSettings(&message.SSettings{}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		})
	case 2:
		_ = NewSettings(&SSettings{
			FMessageSettings:       message.NewSettings(&message.SSettings{}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		})
	case 3:
		_ = NewSettings(&SSettings{
			FMessageSettings:       message.NewSettings(&message.SSettings{}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
		})
	case 4:
		_ = NewSettings(&SSettings{
			FMessageSettings:       message.NewSettings(&message.SSettings{}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		})
	case 5:
		_ = NewSettings(&SSettings{
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		})
	}
}

func TestClosedConn(t *testing.T) {
	t.Parallel()

	listener := testNewService(t, testutils.TgAddrs[30], "")
	defer testFreeService(listener)

	conn, err := Connect(
		context.Background(),
		NewSettings(&SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: testutils.TCWorkSize,
			}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		testutils.TgAddrs[30],
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := conn.Close(); err != nil {
		t.Error(err)
		return
	}

	sett := message.NewConstructSettings(&message.SConstructSettings{
		FSettings: conn.GetSettings().GetMessageSettings(),
	})

	pld := payload.NewPayload32(1, []byte("aaa"))
	msg := message.NewMessage(sett, pld)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conn.WriteMessage(ctx, msg); err == nil {
		t.Error("success write payload to closed connection")
		return
	}

	readCh := make(chan struct{})
	go func() { <-readCh }()

	if _, err := conn.ReadMessage(ctx, readCh); err == nil {
		t.Error("success read payload from closed connection")
		return
	}

	sconn := conn.(*sConn)
	if err := sconn.sendBytes(ctx, []byte("hello, world!")); err == nil {
		t.Error("success send bytes to closed connection")
		return
	}

	if _, err := sconn.recvDataBytes(ctx, 128, time.Second); err == nil {
		t.Error("success recv data bytes from closed connection")
		return
	}
}

func TestInvalidConn(t *testing.T) {
	t.Parallel()

	_, err := Connect(
		context.Background(),
		NewSettings(&SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: testutils.TCWorkSize,
			}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		"INVALID_ADDRESS",
	)
	if err == nil {
		t.Error("success connect to invalid address")
		return
	}
}

func TestReadMessage(t *testing.T) {
	t.Parallel()

	rawConn := &tsConn{}
	conn := LoadConn(
		NewSettings(&SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: testutils.TCWorkSize,
			}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		rawConn,
	).(*sConn)

	wg := sync.WaitGroup{}
	wg.Add(1)

	ch := make(chan struct{})
	rawConn.bodyPart = false
	rawConn.headSize = message.CMessageHeadSize + 10
	rawConn.bodySize = message.CMessageHeadSize + 10
	go func() {
		defer wg.Done()
		ctx := context.Background()
		if _, err := conn.ReadMessage(ctx, ch); err == nil {
			t.Error("success read message with invalid conn 1")
			return
		}
	}()
	<-ch
	wg.Wait()

	wg.Add(1)
	rawConn.cancelBody = true
	rawConn.bodyPart = false
	rawConn.headSize = message.CMessageHeadSize + 10
	rawConn.bodySize = message.CMessageHeadSize + 10
	go func() {
		defer wg.Done()
		ctx := context.Background()
		if _, err := conn.ReadMessage(ctx, ch); err == nil {
			t.Error("success read message with invalid conn 2")
			return
		}
	}()
	<-ch
	wg.Wait()
}

func TestRecvDataBytes(t *testing.T) {
	t.Parallel()

	rawConn := &tsConn{}
	conn := LoadConn(
		NewSettings(&SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: testutils.TCWorkSize,
			}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		rawConn,
	).(*sConn)

	ch := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := conn.recvDataBytes(ctx, 1, 5*time.Second); err == nil {
		t.Error("success recv data bytes with invalid conn 1")
		return
	}

	rawConn.bodyPart = false
	rawConn.headSize = message.CMessageHeadSize + 10
	rawConn.bodySize = message.CMessageHeadSize + 10
	rawConn.readDlError = true
	go func() {
		ctx := context.Background()
		if _, err := conn.recvHeadBytes(ctx, ch, 5*time.Second); err == nil {
			t.Error("success recv data bytes with invalid conn 2")
			return
		}
	}()
	<-ch
}

func TestSendBytes(t *testing.T) {
	t.Parallel()

	rawConn := &tsConn{}
	conn := LoadConn(
		NewSettings(&SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: testutils.TCWorkSize,
			}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		rawConn,
	).(*sConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := conn.sendBytes(ctx, []byte{123}); err == nil {
		t.Error("success send bytes with invalid conn 1")
		return
	}

	if err := conn.sendBytes(context.Background(), []byte{123}); err == nil {
		t.Error("success send bytes with invalid conn 2")
		return
	}
}

func TestRecvHeadBytes(t *testing.T) {
	t.Parallel()

	rawConn := &tsConn{}
	conn := LoadConn(
		NewSettings(&SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: testutils.TCWorkSize,
			}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		rawConn,
	).(*sConn)

	ch := make(chan struct{})

	rawConn.bodyPart = false
	rawConn.headSize = 1
	go func() {
		ctx := context.Background()
		if _, err := conn.recvHeadBytes(ctx, ch, 5*time.Second); err == nil {
			t.Error("success recv head bytes with invalid conn 1")
			return
		}
	}()
	<-ch

	rawConn.bodyPart = false
	rawConn.headSize = math.MaxUint32
	go func() {
		ctx := context.Background()
		if _, err := conn.recvHeadBytes(ctx, ch, 5*time.Second); err == nil {
			t.Error("success recv head bytes with invalid conn 2")
			return
		}
	}()
	<-ch
}

func TestConnWithNetworkKey(t *testing.T) {
	t.Parallel()

	testConn(t, testutils.TgAddrs[17], "")
	// testConn(t, testutils.TgAddrs[17], "hello, world!")
}

func testConn(t *testing.T, pAddr, pNetworkKey string) {
	listener := testNewService(t, pAddr, pNetworkKey)
	defer testFreeService(listener)

	conn, err := Connect(
		context.Background(),
		NewSettings(&SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: testutils.TCWorkSize,
			}),
			FLimitMessageSizeBytes: testutils.TCMessageSize,
			FWaitReadTimeout:       time.Hour,
			FDialTimeout:           time.Minute,
			FReadTimeout:           time.Minute,
			FWriteTimeout:          time.Minute,
		}),
		pAddr,
	)
	if err != nil {
		t.Error(err)
		return
	}

	socket := conn.GetSocket()
	remoteAddr := strings.ReplaceAll(pAddr, "localhost", "127.0.0.1")
	if socket.RemoteAddr().String() != remoteAddr {
		t.Error("got incorrect remote address")
		return
	}

	msgSett := message.NewConstructSettings(&message.SConstructSettings{
		FSettings: conn.GetSettings().GetMessageSettings(),
	})

	pld := payload.NewPayload32(tcHead, []byte(tcBody))
	msg := message.NewMessage(msgSett, pld)
	ctx := context.Background()
	if err := conn.WriteMessage(ctx, msg); err != nil {
		t.Error(err)
		return
	}

	readCh := make(chan struct{})
	go func() { <-readCh }()

	msgRecv, err := conn.ReadMessage(ctx, readCh)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(msgRecv.GetPayload().GetBody(), []byte(tcBody)) {
		t.Error("load payload not equal new payload")
		return
	}
}

func testNewService(t *testing.T, pAddr, pNetworkKey string) net.Listener {
	listener, err := net.Listen("tcp", pAddr)
	if err != nil {
		t.Error(err)
		return nil
	}

	go func() {
		for {
			aconn, err := listener.Accept()
			if err != nil {
				break
			}

			conn := LoadConn(
				NewSettings(&SSettings{
					FMessageSettings: message.NewSettings(&message.SSettings{
						FWorkSizeBits: testutils.TCWorkSize,
						FNetworkKey:   pNetworkKey,
					}),
					FLimitMessageSizeBytes: testutils.TCMessageSize,
					FWaitReadTimeout:       time.Hour,
					FDialTimeout:           time.Minute,
					FReadTimeout:           time.Minute,
					FWriteTimeout:          time.Minute,
				}),
				aconn,
			)

			readCh := make(chan struct{})
			go func() { <-readCh }()

			ctx := context.Background()

			msg, err := conn.ReadMessage(ctx, readCh)
			if err != nil {
				break
			}

			ok := func() bool {
				defer conn.Close()
				return conn.WriteMessage(ctx, msg) == nil
			}()

			if !ok {
				break
			}
		}
	}()

	return listener
}

func testFreeService(listener net.Listener) {
	if listener == nil {
		return
	}
	listener.Close()
}
