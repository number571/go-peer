package anon

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	testutils "github.com/number571/go-peer/test/utils"
)

var (
	_ conn.IConn = &tsConn{}
	_ net.Conn   = &tsNetConn{}
	_ net.Addr   = &tsAddr{}
)

type tsConn struct{}

func (p *tsConn) Close() error                                             { return nil }
func (p *tsConn) GetSettings() conn.ISettings                              { return nil }
func (p *tsConn) GetSocket() net.Conn                                      { return &tsNetConn{} }
func (p *tsConn) WriteMessage(context.Context, net_message.IMessage) error { return nil }
func (p *tsConn) ReadMessage(context.Context, chan<- struct{}) (net_message.IMessage, error) {
	return nil, nil
}

type tsNetConn struct{}

func (p *tsNetConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (p *tsNetConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (p *tsNetConn) Close() error                       { return nil }
func (p *tsNetConn) LocalAddr() net.Addr                { return &tsAddr{} }
func (p *tsNetConn) RemoteAddr() net.Addr               { return &tsAddr{} }
func (p *tsNetConn) SetDeadline(t time.Time) error      { return nil }
func (p *tsNetConn) SetReadDeadline(t time.Time) error  { return nil }
func (p *tsNetConn) SetWriteDeadline(t time.Time) error { return nil }

type tsAddr struct{}

func (p *tsAddr) Network() string { return "tcp" }
func (p *tsAddr) String() string  { return "192.168.0.1:2000" }

const (
	tcService = "TST"
	tcHash    = "96cb1f0968adba001ebc216708a02c8d2817b1a77fad1206012c22716a9b130b"
	tcFmtLog  = "service=TST type=ENQRQ hash=96CB1F09...6A9B130B addr=6245E00D...327047E9 proof=0000012345 size=1024B conn=192.168.0.1:2000"
)

func TestLoggerPanic(t *testing.T) {
	t.Parallel()

	logFunc := GetLogFunc()
	for i := 0; i < 3; i++ {
		testLoggerPanic(t, logFunc, i)
	}
}

func testLoggerPanic(t *testing.T, f logger.ILogFunc, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		f(struct{}{})
	case 1:
		logger := testNewAnonLogger()
		f(logger) // without type
	case 2:
		logger := testNewAnonLogger().WithType(255)
		f(logger) // with unknown type
	}
}

func TestLogger(t *testing.T) {
	t.Parallel()

	logger := testNewAnonLogger().
		WithType(anon_logger.CLogBaseEnqueueRequest)

	logFunc := GetLogFunc()
	if logFunc(logger) != tcFmtLog {
		t.Error("result fmtLog != tcFmtLog")
		return
	}
}

func testNewAnonLogger() anon_logger.ILogBuilder {
	return anon_logger.NewLogBuilder(tcService).
		WithHash(encoding.HexDecode(tcHash)).
		WithProof(12345).
		WithSize(1024).
		WithPubKey(asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0])).
		WithConn(&tsConn{})
}
