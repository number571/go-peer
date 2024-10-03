package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcEnck = "287c039605bb4376a10942cb85080d472e43f883040ff6426799fe97f469bddb5a55fcfb87a73a716d0aebe89ef75e01e0fe57e206744d5efc1ee5ac3e087a7e05089eac529d3502774f14fc1e924f3a8aadbb9b6575f2dde58f42bdc6aa8878b300172a06a8bda40b453505904b30180f18996361e346d34e1cf6850444aa80"
	ecEncd = "799657192fc5eb22e606a948bd82aae6e06e344651109836c1883998c7afcbc274a511f8067cd235a601bee7fb6c52af5e789dfbdfdf5ef165c87b487bc79516d046c473298e8d72878d123a7fdb34ac23c5cd9a61d2f7d35e86af618606e57db7a2e3f7da0d2cc45414a4f2f7c05a9f2946029898a8184a4f2f2d957ab1534ae44d8422844b3c39ab60f1e847d6940e79dddcf2e8a5f04cd01be9a75cd7944234fef55bdbc56e6d90f42fcb134d4a1c543606a1cce46b373dc19b8fb1a64e7d40a0bac7268b3e19167ae31607804d706500d5a07bcc35b4bbd14b65728b1cb4621257c7c688a33b868f44c140c68c98a75e879cf960f9015262b739ae88d03c96fefa34393ee63aa9256c54253e54f3664f0d08a5b7c0daee3f38cbae96f65271edb70ed4852c4424200b2f63e5047af7439d51765c33d9b3e4e1ebcd7b85b16a0ddcf2e9028fc70268024b2ad0fec230eeab6f72ea6a58d01ca0423867fa8a8879b4c49a27790725107069e8de1a3735898e89b1e19fb5b16d2e10794512d5e24913b71797157587e6ff15ba386feaf9dfaa29c29990a1a0810ac97b643e886ab34da06ac444c639b0eeac7b780a3ed26fbca315c717a11156429e8b534c17ec9aef74a704383381dec0b9e7a72d8a71a79f84b30ee64a279657708daa246dcaf2114e993561d93627e978e6da4e9e8dd3ca94d7710c288f54e0d16308e059021f1203f8e057b4e0ccd747fce450f9305b938669d493272a56aef12894a67f86d5fc7c3ab84303ccf8e5787cac7afec735c747c8284ca58c9f4b3dcb7e62130fc25165a977fda33d73b6f403cfdc75825295632bb7ec608c96ce5d9ba5ce8e02aa825647d11b76eddc233ea43522d1ed4d305d6505ae9ee2c8a9334e761f242b1b4ce382027ee5d554a2dc0fffb1144447589d8901d15b9c53759533361dde2b646667d16fbf48fab4efd9f5a51a339f6f4a8fd2c82e158ba0776cab36866553e51ea8c0f4c6d431b5d2fdc9378aac9c86abe6c377726ea0c7d75f26aa47d7da16f61a9013fbc2b9580cce1fd23321b1fb356443cc0743325d28fd39f43406932f360db37d2743f04f9d1008374531cbb4b7af06363abc2d3399f44842d63a5020bbb7e977493de3491b20d7b7b9f3b1fca5188f28cd915e707a741956652be510ba2d8a2688c31764e9f62a244e4bc935e41ec4e2f13eb674c784333fc6d63664f992a94c7bea4a6ab235ebf6c28ec6f04d8e4a6cc442b373f3794ef82498c1481b07d9d84b762b24e4d70f6b42eae9afb08fc97adaf74ebe86661c7ba7ea588dca4cf7159238ea72f9dcf0b42886c61063ebce56d62485e5fbe18a68fad366dc20df624e12b354e775695c9937373ecc8cb3d7f740dbc1d535bc756e66c7b29bb2f763801fafd862e9fa58acbb72dea1a26ce01c18a02d8f193efb969d6a96b63193a4ddeb4da2c6fdc7642380eee9a3b11f67152253cadacc5fa25502dfd92a2969597f25c30927199be394f4871f24e0a38d31e4cca2a724cca5f238c4b0b5989b17c3a48fbdc4788fe957312522433becda133c90e91c61f2058cebf685fdd84bf471fceeaddcf4ab9edd4732ddb37c4371b0c301e10605c1ddda4af844273bbd694dfc210356e40dafe83617f89245759164b848d915e4eb6799cfd9a54e70046114d6c53facb96607955ffe87b769b264b30aa62a885c1b801983579a1cd7c3e7a71c25c8bd4e670a8cd17ce756739e6fd600a5a9728de0dbe7926fd59d5fd7de8b25bb525d9996cdffc5712cbe8124bf665763716b965eb502063ca72a4b59b33b940cbff194cac0b780e3b1d944989fcc1f94063f1c5e8887cd5c9e7bff9beae4da19cd2c78996c22ff5fa5355276dd22af4fa03180f170f29f2ede3e3313ab85fbedab1f7eff96e9790bcf558d2487d935d48f791c63bc7d9e1a428703267b6d0e28694832eb263b8a256a801b759316c28d71f3d912bcfc0a79bc65c079601041c767775c01a6939e6178e7b0c1983029d35dc3780e16e13f15f0e98b9fef8665ccd4012b122f5836a94172affa7758dfcbe069a921f9cf50831bafaf39625c8c4342c1dcf99952be371d4cb6620a80381f4a2fe5f910dd2d1ee02201db586a04d0205b28eddc182022a3e7c9e339997d57d6a2ae402dffede42f2505aabecd509277fcc4f51dc21c876742efcd8fa488fd898ff90c5526f9d9f3273d31fb5d97f09aa5978e134fd3a876dc865921fb66ce69a2083d892ecabc623efdb59ad01ff4c24216770e665121e587e035e3465757a339c3a125a69b42cb7086169a20ef208e4dab0c8c619ef821d48a02702b717cbd45c5d7c8dccfbaa6473273185a470e0a490e65052ca0a439737c7962b36b746dc83449bd509c32e9ff4b3374cd2aa18e2ddb63a225169434ea531019a2f2ecdaf46136c13acc8b6a94a9bf11875e1210ff921d163e38d67930f8d44a684a12dff25bc6b6dda20574c0c8b9f8c13bfbbae985ad8e55df0ef307a928ae6eb8e884c0b357d1c5a34e5cb2e7947124ddff481fef519e55ff7d99c6bf430706ae2ff8f19c602822b3fbbd87ca5876aa8ef1721c8d7274806b5b9e0ad0ce5666c9ba9f4743e7914f89e120f041dfcca9439c6f556fc855596c069bec16387fe2cbddfd7dded5b84e3319a3b2f6c7a3182d120a53c7a10365103f0e03f6c390550838d2c2c0011ca4d0de"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SMessageError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 2; i++ {
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
			FMessageSizeBytes: 1024,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FKeySizeBits: testutils.TcKeySize,
		})
	}
}

func TestPanicLoadMessage(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	params := NewSettings(&SSettings{
		FMessageSizeBytes: 64,
		FKeySizeBits:      testutils.TcKeySize,
	})

	_, _ = loadMessage(params, []byte{123})
}

func TestInvalidMessage(t *testing.T) {
	t.Parallel()

	params := NewSettings(&SSettings{
		FMessageSizeBytes: (2 << 10),
		FKeySizeBits:      testutils.TcKeySize,
	})

	if _, err := LoadMessage(params, struct{}{}); err == nil {
		t.Error("success load message with unknown type")
		return
	}

	if _, err := LoadMessage(params, []byte{123}); err == nil {
		t.Error("success load invalid message")
		return
	}

	msgBytes := joiner.NewBytesJoiner32([][]byte{[]byte("aaa"), []byte("bbb")})
	if _, err := LoadMessage(params, msgBytes); err == nil {
		t.Error("success load invalid message")
		return
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()

	params := NewSettings(&SSettings{
		FMessageSizeBytes: (2 << 10),
		FKeySizeBits:      testutils.TcKeySize,
	})

	msg1, err := LoadMessage(params, testutils.TGBinaryMessage)
	if err != nil {
		t.Error(err)
		return
	}
	testMessage(t, params, msg1)

	msg2, err := LoadMessage(params, testutils.TGStringMessage)
	if err != nil {
		t.Error(err)
		return
	}
	testMessage(t, params, msg2)
}

func testMessage(t *testing.T, params ISettings, msg IMessage) {
	if !bytes.Equal(msg.ToBytes(), testutils.TGBinaryMessage) {
		t.Error("invalid convert to bytes")
		return
	}

	if msg.ToString() != testutils.TGStringMessage {
		t.Error("invalid convert to string")
		return
	}

	if !bytes.Equal(msg.GetEnck(), encoding.HexDecode(tcEnck)) {
		t.Error("incorrect enck")
		return
	}

	if !bytes.Equal(msg.GetEncd(), encoding.HexDecode(ecEncd)) {
		t.Error("incorrect encd")
		return
	}

	msgBytes := bytes.Join([][]byte{msg.GetEnck(), msg.GetEncd()}, []byte{})
	if _, err := loadMessage(params, msgBytes); err != nil {
		t.Error("new message is invalid")
		return
	}
}
