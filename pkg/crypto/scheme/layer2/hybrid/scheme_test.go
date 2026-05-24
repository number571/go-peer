// nolint: err113
package hybrid

import (
	"bytes"
	"errors"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SSchemeError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestInvalidKeys(t *testing.T) {
	t.Parallel()

	if _, err := NewScheme(asymmetric.NewPrivKey(), 8); err == nil {
		t.Fatal("success init scheme with struct size >= message size")
	}
	if _, err := NewScheme(&tsPrivKey{}, (8 << 10)); err == nil {
		t.Fatal("success init scheme with invalid private key")
	}

	_schemePrivKey := asymmetric.NewPrivKey()
	_scheme1, err := NewScheme(_schemePrivKey, (8 << 10))
	if err != nil {
		t.Fatal(err)
	}
	_scheme := _scheme1.(*sScheme)
	if _, err := _scheme.encryptWithPadding(&tsPubKey{}, []byte("hello"), 0); err == nil {
		t.Fatal("success encrypt with invalid pubkey")
	}

	pubKey := _schemePrivKey.GetPubKey()
	msg := []byte("hello, world!")

	enc, err := _scheme.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Fatal(err)
	}

	_scheme.fPrivKey = &tsPrivKey{}
	keysContainer := layer2.NewKeysContainer()
	if _, _, err := _scheme.DecryptMessage(keysContainer, enc); err == nil {
		t.Fatal("success decrypt with invalid privkey")
	}

	_key := make([]byte, symmetric.CCipherKeySize)
	_cipher := symmetric.NewCipherCFB(_key)

	_keysContainer := layer2.NewKeysContainer()
	_keysContainer.Add(_cipher)
	if _, _, err := _scheme.DecryptMessage(_keysContainer, enc); err == nil {
		t.Fatal("success decrypt with another key type")
	}
	if _, err := _scheme.EncryptMessage(_cipher, enc); err == nil {
		t.Fatal("success encrypt with another key type")
	}
}

func TestInvalidScheme(t *testing.T) {
	t.Parallel()

	msgsize := uint64(8 << 10)

	schemePrivKey := asymmetric.NewPrivKey()
	scheme, err := NewScheme(schemePrivKey, msgsize)
	if err != nil {
		t.Fatal(err)
	}

	pubKey := schemePrivKey.GetPubKey()

	_scheme := scheme.(*sScheme)
	msg1 := []byte("hello")
	pad1 := scheme.GetPayloadLimit() - uint64(len(msg1)) + 2*encoding.CSizeUint32

	enc1, err := _scheme.encryptWithPadding(pubKey, msg1, pad1)
	if err != nil {
		t.Fatal(err)
	}

	keysContainer := layer2.NewKeysContainer()
	if _, _, err := scheme.DecryptMessage(keysContainer, enc1); err == nil {
		t.Fatal("success decrypt message with invalid bytes structure (without joiner)")
	}

	pad2 := scheme.GetPayloadLimit() - uint64(len(msg1)) + asymmetric.CDSAPubKeySize - 3
	enc2, err := tcEncryptWithParamsInvalidPKID(_scheme, pubKey, msg1, pad2)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := scheme.DecryptMessage(keysContainer, enc2); err == nil {
		t.Fatal("success decrypt message with invalid dsa public key")
	}
}

func TestScheme(t *testing.T) {
	t.Parallel()

	schemePrivKey := asymmetric.NewPrivKey()
	scheme, err := NewScheme(schemePrivKey, (8 << 10))
	if err != nil {
		t.Fatal(err)
	}

	pubKey := schemePrivKey.GetPubKey()
	msg := []byte("hello, world!")

	enc, err := scheme.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Fatal(err)
	}

	// _ = os.WriteFile("message/test_binary.msg", enc, 0600)
	// _ = os.WriteFile("message/test_string.msg", []byte(encoding.HexEncode(enc)), 0600)

	keysContainer := layer2.NewKeysContainer(pubKey)
	gotPubKey, dec, err := scheme.DecryptMessage(keysContainer, enc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pubKey.ToBytes(), gotPubKey.(asymmetric.IPubKey).ToBytes()) {
		t.Fatal("invalid decrypt key")
	}
	if !bytes.Equal(msg, dec) {
		t.Fatal("invalid decrypt message")
	}

	// fmt.Println(scheme.GetPayloadLimit(), scheme.GetMessageSize())
	// fmt.Println(len(scheme.GetPrivKey().GetPubKey().ToString()))
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	schemePrivKey := asymmetric.NewPrivKey()
	scheme, err := NewScheme(schemePrivKey, (8 << 10))
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := scheme.DecryptMessage(layer2.NewKeysContainer(), []byte{123}); err == nil {
		t.Fatal("success decrypt with invalid ciphertext (1)")
	}

	pubKey := schemePrivKey.GetPubKey()
	msg := []byte("hello, world!")

	enc, err := scheme.EncryptMessage(pubKey, msg)
	if err != nil {
		t.Fatal(err)
	}

	mapKeys := layer2.NewKeysContainer(pubKey)

	if _, _, err := scheme.DecryptMessage(mapKeys, enc); err != nil {
		t.Fatal(err)
	}

	enc[0] ^= 1
	if _, _, err := scheme.DecryptMessage(mapKeys, enc); err == nil {
		t.Fatal("success decrypt with invalid ciphertext (2)")
	}

	enc[0] ^= 1
	enc[len(enc)-1] ^= 1
	if _, _, err := scheme.DecryptMessage(mapKeys, enc); err == nil {
		t.Fatal("success decrypt with invalid ciphertext (3)")
	}

	enc[len(enc)-1] ^= 1
	enc[asymmetric.CKEMCiphertextSize+symmetric.CCipherBlockSize+2*encoding.CSizeUint32+hashing.CHasherSize+1] ^= 1
	if _, _, err := scheme.DecryptMessage(mapKeys, enc); err == nil {
		t.Fatal("success decrypt with invalid ciphertext (4)")
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
