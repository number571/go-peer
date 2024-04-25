package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcEnck = "7c596e0796260f10a4c19f1545aa69deaaf367693e26249b7190d597233a9a3cc9e38402a94bace12d0c229f64093bb16205447598aef834a50cac1d4659f6f564b20c86132ac231cb16476004c8a6ec86d2913cd5abc08bd3d766a6870eef2b10cfd0b9b9584bc5dc7dfa34503d1b7ae29057f321ae79ba27634a9b82218ea6"
	ecEncd = "f2cc9d35729680f146bd4e1c0920dc51ebab433ddc9eca0d7aa5e44a94f9e2316a17ca6b729dd5f6e1744a9cb955179297e2758d0a4192fb4211220eff74ea98c4e38d317baa6fb904e2d4a76282fa8d29d30720a8c6cea10d602ca09120f65c8d7bc6330de9beebe2a3a3af6d4764d0b3ba8c2b80aa7ea908ba2193f6b8d65f1e60cffe1eaff02a32712995c83aa6caaac90f50307ed450f66724d810461049b02d84b64b4dc00686db155383faa4a45dca890437e9c85d1c6cbcdc99da0d19a1789ab164a974432c8ac1b7b98d0d6da8f6f151eb8fe0470fc9a59e576c432fe82d9136409736e2e368105de3ec94d6371f94c67b19fc0849c9b504bd1568fbf27bc60de758da7050c6081083a30aa8f3725ecfb8e1c67cef9097af29ebb66f0cf9c2dd23f16638b4555166bd37ee35bf1d4a819681e41c0a23493ea031ddc1222f8d8daec3b1fcfb551be6a33593fdd023c220e2426adfae87eadbc5224850a1b2ad0a8bc16157de456d2228878c5029aad7601d32df7b37bca69ef9621d359d1845e5d6b76c5a2166b9b2f677400ea047b19df5aa054685c930406175b2b6bd2170ff07bc5a045183f11a1505fb7a049e67a7db80cef53a035b2c6d574bbdaa2f39d29053e3dcacc14eb6f922b1e680201db1faebb3baa5d87871e0300d7884a73cae81d7dd9a658c4b8382dcb58bf149ffc77ea5410cc3656b7a36d4ead61087d3cca8ca3e0fbec2d0efc897932464714d7b103de8b92c03ed91f4126910457237e695fa68432b46ea61b31745080179c2361a4386b5eb932ea7fb6e54ab53c0e0cd7290674894e9351b7a6d3f6a0302fa77a89aa7bd0548e40920c272102e9e86a0095ab0344c59a408ccde2afd19183eaaa93c105c33b8cef188e73135c6ef4a0543e0a22500122c4ff5a8bd3fbd6e919ad8b1a584fe6b006bdf1da7a765c278688d516fcef854911905b8b44241dbae2245a04047f945a996215afdf6cf1b3d5c43fa17b7aa4fd0db2464ab58b40da903513edcadfa2e293221886ef8e30a007f3240394d810bec8b34c807905cc5db97d601fe4649f8180369d8ccd6669606cebf98a0ce37e2151753ee4b2edd79ddba2de59370b653358f4dc38cbf8b376d4a7ef536e4065a088222d6dba752e484171601e47f932825ea97521d4815aa205d4a8cd64c910ff65468d2b5a83de81ebacafe7abd2cefc87607c58026839c2a61f6f44d3cfff03eb6b349368935410b6699e667f4f7b52952ba33589f202450f392f871261e5c02664d7f582bb2653ce4c65125d0db14dfa0d93b81abdc33b6b31975b3afe9e47351081e0d1976f23fe3d9108b7e7b9c53eb1554103bdc0d79172bbdb1f0e77d40e667909f148550ec9a40b3d1696e08ac264474b823a0b27ce78a4eac40a04ef7d07bb281fd8d729a5b1fcca151d138b08eb6c7cbba44bd8ecd0872d50bafab8a8686a59071c69c388f3cba59ea832b1c6235636fbec636fa8b5fcf85faed6b6dd487992820034de72815d3bae64ea36b3b86fed5f535b26409d7d6dc0524c6a47bc4e7bbbfddb982329c5b5d1a6415df8eaec5ab0c6b363b4718758985287f5ffab07a18080ceb918f81c89e19e69e129a47fdd122c964f64affcba35e397bc789d8e15385b7156ad6be05dcbf8e94ba33a409afb41fbd88aa857a608674fe2dc238644d44fec918be214ff33b47d970d1120f44f20d51045c95ba14117064217bf2251850cd4bd1ef3b080210bc7a686c61f2db467155b34a1bc3526f27a27f9c966c5881853c8b60e85e5792f55b8f52eb88176fc477e7c33d2c2666c6c4ba0a4419dd1925987d1c4dde028d178345f9fa111f667b7019fb31aada13fb695d82455447e683b88fe4cec825f9fc1e5c08942718dab58fda90cf2f6dd6e78ccd6d0b3b5392c99a90cd14f75814a7c9348171a32752fa6cf3207b94657b90ef4d2f60db2a1fa557b035a008ae00ccdcb0eac50e0ce494344c373a91b9ceda876d36e6aa194e6ece69dd614cf458ae71a51dae96c2addb3ed2233f71181df8bcec13f8a49106b3d5ef73986bd1fafa6e5e81d60493d461b4ef4f07716d29d58a8c1d7803b7e6a3a3f7c47a30c99db7d72409e32ae5b8bc6345174bb13c62e1dbfdca58ba9566e70ceafe277bf27cc28132d499a734d93f40c3cf1518984d5f5e50848014c6e1ca32d0c22c5ff3a00181b1c6f54a7c2b6805087e68f8079eab897965702b307a39f89402e36ade0ef1ea1a399fd63c939147a5bf6217d45da4698f8247b837ea840052010e4da6093fb0d37e3c4a61c95db4c6cb14f9788f5d45e84788ec5e13c7f61bcb6c4138ecde34edbbe07e6126e8f9a6adf7bdbbff74a3c538d30e2916e8276c8a4303a0572a33bd486d472e8430152f70d3595cf915ce6b417b76e186bfebdaae0ba41cc61627c4715d7ee741127c27b0b34e53c3a010d6e795e21cee64707ff7c1fa3b79262ae751744010768bc2d3ceef461364ac301944609a7447b66ad56f10a0830da390ad92d5d90f577291db0262295489e40b87a33baab80e953c557b7c694fba6d35cde0cb2757771284f1a931ab65f52f111c982caf01476edcd6e2a078f494347b205f181c4424bd3ba685a1ea0ab0e7489b397d08da754a1d849289d01896b30c72c1935e7448281b9db9723c263aa2c25f570738f1ac"
)

func TestError(t *testing.T) {
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

	msg1 := NewMessage(msg.GetEnck(), msg.GetEncd())
	if !msg1.IsValid(params) {
		t.Error("new message is invalid")
		return
	}
}
