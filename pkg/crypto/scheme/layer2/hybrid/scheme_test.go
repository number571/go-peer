// nolint: err113
package hybrid

import (
	"bytes"
	"errors"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

func TestPanicNewScheme(t *testing.T) {
	t.Parallel()

	tcNewSchemeWithSmallMsgSize(t)
	tcNewSchemeWithInvalidPrivKey(t)
}

func tcNewSchemeWithSmallMsgSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewScheme(asymmetric.NewPrivKey(), 8)
}

func tcNewSchemeWithInvalidPrivKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewScheme(&tsPrivKey{}, (8 << 10))
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SSchemeError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestInvalidKeys(t *testing.T) {
	t.Parallel()

	_schemePrivKey := asymmetric.NewPrivKey()
	_scheme := NewScheme(_schemePrivKey, (8 << 10)).(*sScheme)
	if _, err := _scheme.encryptWithPadding(&tsPubKey{}, []byte("hello"), 0); err == nil {
		t.Error("success encrypt with invalid pubkey")
		return
	}

	pubKey := _schemePrivKey.GetPubKey()
	msg := []byte("hello, world!")

	enc, err := _scheme.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	_scheme.fPrivKey = &tsPrivKey{}
	mapPubKeys := asymmetric.NewMapPubKeys()
	if _, _, err := _scheme.DecryptMessage(mapPubKeys, enc); err == nil {
		t.Error("success decrypt with invalid privkey")
		return
	}

	_key := make([]byte, symmetric.CCipherKeySize)
	if _, _, err := _scheme.DecryptMessage(symmetric.NewCipherCFB(_key), enc); err == nil {
		t.Error("success decrypt with another key type")
		return
	}
	if _, err := _scheme.EncryptMessage(symmetric.NewCipherCFB(_key), enc); err == nil {
		t.Error("success encrypt with another key type")
		return
	}
}

func TestInvalidScheme(t *testing.T) {
	t.Parallel()

	msgsize := uint64(8 << 10)

	schemePrivKey := asymmetric.NewPrivKey()
	scheme := NewScheme(schemePrivKey, msgsize)
	pubKey := schemePrivKey.GetPubKey()

	_scheme := scheme.(*sScheme)
	msg1 := []byte("hello")
	pad1 := scheme.GetPayloadLimit() - uint64(len(msg1)) + 2*encoding.CSizeUint32

	enc1, err := _scheme.encryptWithPadding(pubKey, msg1, pad1)
	if err != nil {
		t.Error(err)
		return
	}

	mapKeys := asymmetric.NewMapPubKeys(pubKey)
	if _, _, err := scheme.DecryptMessage(mapKeys, enc1); err == nil {
		t.Error("success decrypt message with invalid bytes structure (without joiner)")
		return
	}

	pad2 := scheme.GetPayloadLimit() - uint64(len(msg1)) + asymmetric.CDSAPubKeySize - 3
	enc2, err := tcEncryptWithParamsInvalidPKID(_scheme, pubKey, msg1, pad2)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := scheme.DecryptMessage(mapKeys, enc2); err == nil {
		t.Error("success decrypt message with invalid dsa public key")
		return
	}
}

func TestScheme(t *testing.T) {
	t.Parallel()

	schemePrivKey := asymmetric.NewPrivKey()
	scheme := NewScheme(schemePrivKey, (8 << 10))

	pubKey := schemePrivKey.GetPubKey()
	msg := []byte("hello, world!")

	enc, err := scheme.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	// _ = os.WriteFile("message/test_binary.msg", enc, 0600)
	// _ = os.WriteFile("message/test_string.msg", []byte(encoding.HexEncode(enc)), 0600)

	gotPubKey, dec, err := scheme.DecryptMessage(asymmetric.NewMapPubKeys(pubKey), enc)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(pubKey.ToBytes(), gotPubKey.(asymmetric.IPubKey).ToBytes()) {
		t.Error("invalid decrypt key")
		return
	}
	if !bytes.Equal(msg, dec) {
		t.Error("invalid decrypt message")
		return
	}

	// fmt.Println(scheme.GetPayloadLimit(), scheme.GetMessageSize())
	// fmt.Println(len(scheme.GetPrivKey().GetPubKey().ToString()))
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	schemePrivKey := asymmetric.NewPrivKey()
	scheme := NewScheme(schemePrivKey, (8 << 10))

	if _, _, err := scheme.DecryptMessage(asymmetric.NewMapPubKeys(), []byte{123}); err == nil {
		t.Error("success decrypt with invalid ciphertext (1)")
		return
	}

	pubKey := schemePrivKey.GetPubKey()
	msg := []byte("hello, world!")

	enc, err := scheme.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	mapKeys := asymmetric.NewMapPubKeys(pubKey)

	if _, _, err := scheme.DecryptMessage(mapKeys, enc); err != nil {
		t.Error(err)
		return
	}

	enc[0] ^= 1
	if _, _, err := scheme.DecryptMessage(mapKeys, enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (2)")
		return
	}

	enc[0] ^= 1
	enc[len(enc)-1] ^= 1
	if _, _, err := scheme.DecryptMessage(mapKeys, enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (3)")
		return
	}

	enc[len(enc)-1] ^= 1
	enc[asymmetric.CKEMCiphertextSize+symmetric.CCipherBlockSize+2*encoding.CSizeUint32+hashing.CHasherSize+1] ^= 1
	if _, _, err := scheme.DecryptMessage(mapKeys, enc); err == nil {
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
	p *sScheme,
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

	cipher := symmetric.NewCipherCFB(sk)
	return newMessage(
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
