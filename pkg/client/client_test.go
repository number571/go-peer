// nolint: goerr113
package client

import (
	"bytes"
	"errors"
	"testing"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

func TestPanicNewClient(t *testing.T) {
	t.Parallel()

	tcNewClientWithSmallMsgSize(t)
	tcNewClientWithInvalidPrivKey(t)
}

func tcNewClientWithSmallMsgSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewClient(asymmetric.NewPrivKey(), 8)
}

func tcNewClientWithInvalidPrivKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewClient(&tsPrivKey{}, (8 << 10))
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SClientError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestInvalidKeys(t *testing.T) {
	t.Parallel()

	_client := NewClient(asymmetric.NewPrivKey(), (8 << 10)).(*sClient)
	if _, err := _client.encryptWithParams(&tsKEMPubKey{}, []byte("hello"), 0); err == nil {
		t.Error("success encrypt with invalid pubkey")
		return
	}

	kemPubKey := _client.GetPrivKey().GetKEMPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := _client.EncryptMessage(kemPubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	_client.fPrivKey = &tsPrivKey{}
	if _, _, err := _client.DecryptMessage(enc); err == nil {
		t.Error("success decrypt with invalid privkey")
		return
	}
}

func TestInvalidClient(t *testing.T) {
	t.Parallel()

	msgsize := uint64(8 << 10)
	client := NewClient(asymmetric.NewPrivKey(), msgsize)
	kemPubKey := client.GetPrivKey().GetKEMPrivKey().GetPubKey()

	_client := client.(*sClient)
	msg1 := []byte("hello")
	pad1 := client.GetPayloadLimit() - uint64(len(msg1)) + 2*encoding.CSizeUint32

	enc1, err := _client.tcEncryptWithParamsInvalidMessageBytes(kemPubKey, msg1, pad1)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client.DecryptMessage(enc1); err == nil {
		t.Error("success decrypt message with invalid bytes structure (without joiner)")
		return
	}

	pad2 := client.GetPayloadLimit() - uint64(len(msg1)) + asymmetric.CDSAPubKeySize - 3
	enc2, err := _client.tcEncryptWithParamsInvalidDSAPublicKey(kemPubKey, msg1, pad2)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client.DecryptMessage(enc2); err == nil {
		t.Error("success decrypt message with invalid dsa public key")
		return
	}
}

func TestClient(t *testing.T) {
	t.Parallel()

	client := NewClient(asymmetric.NewPrivKey(), (8 << 10))

	kemPubKey := client.GetPrivKey().GetKEMPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := client.EncryptMessage(kemPubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	// _ = os.WriteFile("message/test_binary.msg", enc, 0600)
	// _ = os.WriteFile("message/test_string.msg", []byte(encoding.HexEncode(enc)), 0600)

	pubKey := client.GetPrivKey().GetPubKey()
	gotPubKey, dec, err := client.DecryptMessage(enc)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(pubKey.ToBytes(), gotPubKey.ToBytes()) {
		t.Error("invalid decrypt key")
		return
	}
	if !bytes.Equal(msg, dec) {
		t.Error("invalid decrypt message")
		return
	}

	// fmt.Println(client.GetPayloadLimit(), client.GetMessageSize())
	// fmt.Println(len(client.GetPrivKey().GetPubKey().ToString()))
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	client := NewClient(asymmetric.NewPrivKey(), (8 << 10))

	if _, _, err := client.DecryptMessage([]byte{123}); err == nil {
		t.Error("success decrypt with invalid ciphertext (1)")
		return
	}

	kemPubKey := client.GetPrivKey().GetKEMPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := client.EncryptMessage(kemPubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	enc[0] ^= 1
	if _, _, err := client.DecryptMessage(enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (2)")
		return
	}

	enc[0] ^= 1
	enc[len(enc)-1] ^= 1
	if _, _, err := client.DecryptMessage(enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (3)")
		return
	}

	enc[len(enc)-1] ^= 1
	enc[len(enc)-2000] ^= 1
	if _, _, err := client.DecryptMessage(enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (4)")
		return
	}
}

var (
	_ asymmetric.IPrivKey    = &tsPrivKey{}
	_ asymmetric.IKEMPubKey  = &tsKEMPubKey{}
	_ asymmetric.IDSAPubKey  = &tsDSAPubKey{}
	_ asymmetric.IKEMPrivKey = &tsKEMPrivKey{}
	_ asymmetric.IDSAPrivKey = &tsDSAPrivKey{}
)

type tsPrivKey struct{}
type tsKEMPubKey struct{}
type tsDSAPubKey struct{}
type tsKEMPrivKey struct{}
type tsDSAPrivKey struct{}

func (p *tsPrivKey) ToString() string                      { return "" }
func (p *tsPrivKey) ToBytes() []byte                       { return nil }
func (p *tsPrivKey) GetPubKey() asymmetric.IPubKey         { return nil }
func (p *tsPrivKey) GetKEMPrivKey() asymmetric.IKEMPrivKey { return &tsKEMPrivKey{} }
func (p *tsPrivKey) GetDSAPrivKey() asymmetric.IDSAPrivKey { return &tsDSAPrivKey{} }

func (p *tsKEMPubKey) ToBytes() []byte { return nil }
func (p *tsKEMPubKey) Encapsulate() ([]byte, []byte, error) {
	return nil, nil, errors.New("some error")
}

func (p *tsKEMPrivKey) ToBytes() []byte                    { return nil }
func (p *tsKEMPrivKey) GetPubKey() asymmetric.IKEMPubKey   { return &tsKEMPubKey{} }
func (p *tsKEMPrivKey) Decapsulate([]byte) ([]byte, error) { return nil, errors.New("some error") }

func (p *tsDSAPrivKey) ToBytes() []byte                  { return nil }
func (p *tsDSAPrivKey) GetPubKey() asymmetric.IDSAPubKey { return &tsDSAPubKey{} }
func (p *tsDSAPrivKey) SignBytes([]byte) []byte          { return nil }

func (p *tsDSAPubKey) ToBytes() []byte                 { return nil }
func (p *tsDSAPubKey) VerifyBytes([]byte, []byte) bool { return false }

func (p *sClient) tcEncryptWithParamsInvalidMessageBytes(
	pRecv asymmetric.IKEMPubKey,
	pMsg []byte,
	pPadd uint64,
) ([]byte, error) {
	var (
		rand = random.NewRandom()
		salt = rand.GetBytes(cSaltSize)
		sign = p.fPrivKey.GetDSAPrivKey()
	)

	data := bytes.Join([][]byte{pMsg, rand.GetBytes(pPadd)}, []byte{})
	hash := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{
			sign.GetPubKey().ToBytes(),
			pRecv.ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()

	ct, sk, err := pRecv.Encapsulate()
	if err != nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewCipher(sk)
	return message.NewMessage(
		ct,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			sign.GetPubKey().ToBytes(),
			salt,
			hash,
			sign.SignBytes(hash),
			data,
		})),
	).ToBytes(), nil
}

func (p *sClient) tcEncryptWithParamsInvalidDSAPublicKey(
	pRecv asymmetric.IKEMPubKey,
	pMsg []byte,
	pPadd uint64,
) ([]byte, error) {
	var (
		rand = random.NewRandom()
		salt = rand.GetBytes(cSaltSize)
		sign = p.fPrivKey.GetDSAPrivKey()
	)

	data := joiner.NewBytesJoiner32([][]byte{pMsg, rand.GetBytes(pPadd)})
	hash := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{
			sign.GetPubKey().ToBytes(),
			pRecv.ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()

	ct, sk, err := pRecv.Encapsulate()
	if err != nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewCipher(sk)
	return message.NewMessage(
		ct,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			[]byte("123"),
			salt,
			hash,
			sign.SignBytes(hash),
			data,
		})),
	).ToBytes(), nil
}
