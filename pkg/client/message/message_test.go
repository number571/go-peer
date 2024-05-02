package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcEnck = "4ac068c801a8b48637ae929440d829ed793e3a076477634542a8730ae7ad0e1987756316b6e5b41f8cdcf20c0b68616c1c6581424285eb984cf438bc8314bc02dc918b79f7f524ad7c9111922db7d516a2fb679044059d282b31e9d58a34abaee21ce9e9abc77c38ea5eca5cd53e29f298d3e081389b894bc2c2faee783d4c42"
	ecEncd = "04b775815e5dc0573cb949c1d681b733fba8855626997d97e11ca352cd47e0184ef35ccbe1ede94216eb874fbce503201ea24ceb66650a2cdbb8b2edc2f114cb6004aadfe56f9aa39fa391adf2adbf54da29dd124826132c498e569f64e069fd185632336b446be93f4996b5541098a345b4b95a29970aa71a9e4391d1cecbf586a1e6a7704e5ae2e6b7da9a96925bb564749654963cbfec56e53c63341912be9178b0136da9dc64c3f8c76143fbf5e03cb12b0f49b749fac08639a01c3202a65c22bbba8c985c1f4f3d2fed44f57c649375a8d5561f9735385a7aaba111a4ad8eb9ab14f8a4d54f0518f39332e97a79092573e7534a84b6897f326c5c24bbe81d26418249a951e02f591dcfb1556a8ae1fd7dfd5697eb2c2de35a2695623da14f596ef5dab55694b88754f7786d8ba505019b2fca21d49f6e0dbf140177577d5d9b3ce37960d8d2c32312cb1dcc0473b61c7342349c9733bc36b226be18c39a0bdbef3f4ca97b8fada4e3e1d948947bcea37b9de526a94638e36701f1538f49f91b16706f03011b2bb59cb388ca63b4fb6f16053472e2a5e6600d22e4caa1cc5de9844a06a88e05cc51b5c5a33f456534cba50cd7af2ac20acd27eca395eb0bf3e8ef778401c3328c12a8a9ab242e6c867eaa5e600e53ca313b66e4f7aff6a5ce3901a23efce3187f4065caa56b8725f8c41e2947d18289434ae4273b952de8181bb55696eddf5b1ce4846f38538b74afe8711f6feb510696aaa2881ff67447251de8fbe125837c8e6d7fa5bde64d8cae888fcf7758a4cf867a372f00e5659cc1b93c9684de86677541c076f6179b3534750a7c4fbb238a46caeaf5d2fe558970dca638957085fe3b845eef44aea0b7f19a70cf621781192685f66b66fa461e835d44980d9e19c064fb481428a911a9c0a6ff9e2c3982b818e9b559764b4ba4a9dd0dbd2f424b6f62f5ca83b6574c0553fe087382474f185a3be83b0a01c3e00b796722838f5e9328c60d538da732705fb8ab8ff5b43b58ecafbb5a448d7079c3fa8775ca3bdd4e56e60d06863c8d14750a99909eb5f13c83a0a4f79bef47024a581a782b0113d1d006fb2e3133094e9ed13783c4c062c8ae53bbfdc791c8164876ed4d96e68bd3d1a2d7b8e984f85de97fbd4bb0cbc65f18d82fea564ac6967af997684518083b81ac0cf0c3e8082e73806a4514e2a8b4a6085b833209df9c58715890b0f58c7bae6fa304b8afc9bcd4a7d54a08f9ef97b4f03780dfdb9696c630f52aae82a195bae213eaa8d79771f2eaae67780e049ba37aa0775da830df814b4a512cb2153b455abb775c5c52ad577fef4f4a4526f0467c3563aa45eda51b5bc84216e3f4e63ce41be2daece8e3a9406b6bcbd98a3b92d1ee66a6f9bf08e1b35a977588b93463e979d29515447ece633ed0b6f12301291e26fcc83bc0f4dce8d20f76757fd31157c22b1f9092b41e2372b0c9243fedea489ed76998d4ae297abe23e847c19b9a2081a9728c0297b0d5b01a79d2b9217fc7647809ff70b7a3460cf6cd727d000e0e7c5a4817f4e869afb0af27377482835a2303a73749414bcf17ae14ea2b6477018f11a068d9cf3b6af16ede29c55ca5b2fc88e5493853f37e782dcca9ff093c32b4274974b0d023e596166734d39bd8ace96b51413cb832d829b75838e73b7232514321496be18cf06522402f351237aad0da9d8cbc5705e2688a5368f3d52f633c773a82cd6e215802b601f1548f63bc9fcd51c00295e10bfdf58f69ed4056d4af5eb3dee517de240a1cad853274238821dc18aeb72ea936b80e7e80ed069232266a7fe657c3b7f4a1afdb574757544aa9476a4ecc71e346035bcdfd9b25c24e0b6cef9544bd5dc093085eeea2641a09b9d8e87f65d60433a6ea4aa6355024a1a25e6d131bf4bd0ac09e44141ead3b728965441a4a3f094029134129c48a0f88c4ce8bc0c2f74975f6afc36401ec76a80bb394be6eb358391ee0f8474947ba4c45b63b038c61c7d36d980e811821b5b325f53b247bb258fa144cba667cb8657434275b1eb4cb0dc6a0919b43b5d525174c4e7a47dd57d5999bae76abdb2ae2df8483745b039c479d4830c5096c637141d46b8c628e885f4bfad5a05f3d2714eabdacd79f6ce43d7280ab49480fe02e1ec0ec61bcf86cdf052c6d2087908be237fba09e20dd80972c857c59ff60b2e20dd9451efa999cdbd2fa697c74694966b02a52d20d993b91be76d01a7db727d90de582142b569a43311a1bbf73c8904ff6b0aedd64e6d9fadaa45d7d75e7d238a612f87768473a9529ea44886f941e6dab3ad856c15297636b89a9394773162e66c7ad852fb45c21addf96c702c3b2212a1e35306ba2de2f739e66164a9eeca60dcbae4900cadd6151cb58e1e1920a2f1d487fcd5bb276a17e1e40490cb62a1af932f961ddf8008ef914fec1dcf283720baa7e28070b7c982ee405d446540def8ec600486586dedab8267bee0cb7ad07f29bf8b34a91d53a98de1381ba7df0d1ad257d9f555de173cd976beadee3ddfec028c18f6c5ea67b01b39e82b6652481e9f183c0bd8e1184b04049c07f1b9dd569127a8f5dc8751194f1c47268d58e5202229684088da339b91aa0f3f6b71df7872aa192e3be048534ac1b38d98195d7cbd09b43b3f533bfbeee4368d5b227352c449a0f87edab40ae9a8e16d1f007"
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

	msg1 := NewMessage(msg.GetEnck(), msg.GetEncd())
	if !msg1.IsValid(params) {
		t.Error("new message is invalid")
		return
	}
}
