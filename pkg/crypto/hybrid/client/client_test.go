// nolint: err113
package client

import (
	"bytes"
	"errors"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/message/layer2"
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
	if _, err := _client.encryptWithPadding(&tsPubKey{}, []byte("hello"), 0); err == nil {
		t.Error("success encrypt with invalid pubkey")
		return
	}

	pubKey := _client.GetPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := _client.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	_client.fPrivKey = &tsPrivKey{}
	mapPubKeys := asymmetric.NewMapPubKeys()
	if _, _, err := _client.DecryptMessage(mapPubKeys, enc); err == nil {
		t.Error("success decrypt with invalid privkey")
		return
	}
}

func TestInvalidClient(t *testing.T) {
	t.Parallel()

	msgsize := uint64(8 << 10)
	client := NewClient(asymmetric.NewPrivKey(), msgsize)
	pubKey := client.GetPrivKey().GetPubKey()

	_client := client.(*sClient)
	msg1 := []byte("hello")
	pad1 := client.GetPayloadSize() - uint64(len(msg1)) + 2*encoding.CSizeUint32

	enc1, err := _client.encryptWithPadding(pubKey, msg1, pad1)
	if err != nil {
		t.Error(err)
		return
	}

	mapKeys := asymmetric.NewMapPubKeys(pubKey)
	if _, _, err := client.DecryptMessage(mapKeys, enc1); err == nil {
		t.Error("success decrypt message with invalid bytes structure (without joiner)")
		return
	}

	pad2 := client.GetPayloadSize() - uint64(len(msg1)) + asymmetric.CDSAPubKeySize - 3
	enc2, err := tcEncryptWithParamsInvalidPKID(_client, pubKey, msg1, pad2)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client.DecryptMessage(mapKeys, enc2); err == nil {
		t.Error("success decrypt message with invalid dsa public key")
		return
	}
}

func TestClient(t *testing.T) {
	t.Parallel()

	client := NewClient(asymmetric.NewPrivKey(), (8 << 10))

	pubKey := client.GetPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := client.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	// _ = os.WriteFile("message/test_binary.msg", enc, 0600)
	// _ = os.WriteFile("message/test_string.msg", []byte(encoding.HexEncode(enc)), 0600)

	gotPubKey, dec, err := client.DecryptMessage(asymmetric.NewMapPubKeys(pubKey), enc)
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

	if _, _, err := client.DecryptMessage(asymmetric.NewMapPubKeys(), []byte{123}); err == nil {
		t.Error("success decrypt with invalid ciphertext (1)")
		return
	}

	pubKey := client.GetPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := client.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	mapKeys := asymmetric.NewMapPubKeys(pubKey)

	if _, _, err := client.DecryptMessage(mapKeys, enc); err != nil {
		t.Error(err)
		return
	}

	enc[0] ^= 1
	if _, _, err := client.DecryptMessage(mapKeys, enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (2)")
		return
	}

	enc[0] ^= 1
	enc[len(enc)-1] ^= 1
	if _, _, err := client.DecryptMessage(mapKeys, enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (3)")
		return
	}

	enc[len(enc)-1] ^= 1
	enc[asymmetric.CKEMCiphertextSize+symmetric.CCipherBlockSize+2*encoding.CSizeUint32+hashing.CHasherSize+1] ^= 1
	if _, _, err := client.DecryptMessage(mapKeys, enc); err == nil {
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
type tsPubKey struct{}
type tsKEMPubKey struct{}
type tsDSAPubKey struct{}
type tsKEMPrivKey struct{}
type tsDSAPrivKey struct{}

func (p *tsPubKey) ToString() string                    { return "" }
func (p *tsPubKey) ToBytes() []byte                     { return nil }
func (p *tsPubKey) GetHasher() hashing.IHasher          { return hashing.NewHasher([]byte{}) }
func (p *tsPubKey) GetKEMPubKey() asymmetric.IKEMPubKey { return &tsKEMPubKey{} }
func (p *tsPubKey) GetDSAPubKey() asymmetric.IDSAPubKey { return &tsDSAPubKey{} }

func (p *tsPrivKey) ToString() string                      { return "" }
func (p *tsPrivKey) ToBytes() []byte                       { return nil }
func (p *tsPrivKey) GetPubKey() asymmetric.IPubKey         { return &tsPubKey{} }
func (p *tsPrivKey) GetKEMPrivKey() asymmetric.IKEMPrivKey { return &tsKEMPrivKey{} }
func (p *tsPrivKey) GetDSAPrivKey() asymmetric.IDSAPrivKey { return &tsDSAPrivKey{} }

func (p *tsKEMPubKey) ToBytes() []byte { return nil }
func (p *tsKEMPubKey) Encapsulate() ([]byte, []byte, error) {
	return nil, nil, errors.New("some error") //nolint:err113
}

func (p *tsKEMPrivKey) ToBytes() []byte                    { return nil }
func (p *tsKEMPrivKey) GetPubKey() asymmetric.IKEMPubKey   { return &tsKEMPubKey{} }
func (p *tsKEMPrivKey) Decapsulate([]byte) ([]byte, error) { return nil, errors.New("some error") } //nolint:err113

func (p *tsDSAPrivKey) ToBytes() []byte                  { return nil }
func (p *tsDSAPrivKey) GetPubKey() asymmetric.IDSAPubKey { return &tsDSAPubKey{} }
func (p *tsDSAPrivKey) SignBytes([]byte) []byte          { return nil }

func (p *tsDSAPubKey) ToBytes() []byte                 { return nil }
func (p *tsDSAPubKey) VerifyBytes([]byte, []byte) bool { return false }

func tcEncryptWithParamsInvalidPKID(
	p *sClient,
	pRecv asymmetric.IPubKey,
	pMsg []byte,
	pPadd uint64,
) ([]byte, error) {
	var (
		rand = random.NewRandom()
		salt = rand.GetBytes(cSaltSize)
		pkid = p.fPrivKey.GetPubKey().GetHasher().ToBytes()
	)

	data := joiner.NewBytesJoiner32([][]byte{pMsg, rand.GetBytes(pPadd)})
	hash := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{
			pkid,
			pRecv.ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()

	ct, sk, err := pRecv.GetKEMPubKey().Encapsulate()
	if err != nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewCipher(sk)
	return layer2.NewMessage(
		ct,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			[]byte("123"),
			salt,
			hash,
			p.fPrivKey.GetDSAPrivKey().SignBytes(hash),
			data,
		})),
	).ToBytes(), nil
}
