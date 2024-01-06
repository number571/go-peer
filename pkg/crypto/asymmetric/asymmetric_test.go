package asymmetric

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
)

const (
	tcKeySize = 2048
	tcPrivKey = `PrivKey{308204A50201000282010100D7114DDEDB12344F4C7086819B91DCE5FE4DD914306BBAFDDA2EB1080A43F98EC0281DD07B5D1F9229DF38E2C986F8CA6B2AC42E4E8E9F021A38E82A03388B4F2314D58D738BB9D77243CC7D8641C574B4CDC00A1D2720329E1162F1EDC7FDA9BF2F9948E927CCF6C274321473E769AA68CB3669729B19BC1CE022E9BA1E0683F3970330ED7FC7F9E836CD1CDB67B571CD98B82FBEB624948E4BD3CF161250311F2661BB6D0885C3E81B66219CB3407619DF239427214328F1C10CC591C490417DB13CB6B90EC22FEC7B7E4BE8B9B3B75E9820D6446AE7F5690CE0ED048DA5ACE47105A1441F4C765CF345499EBFB55B1FFCFE9CAEF2795E3826075AECF2CED502030100010282010100B9E90F7371D44EBBADCC27B9AA0D70F2AFDE03A4DC2684422474F03B8F042B9A26A986FC4D67B67ED70B4B555FF7F8E0A1BB1A531D3D545EB0E4386CF8D3CC38E08E85FBFCC1F02839723A36D7F3CB0893B2B82B06006868D9131681239719C3BEAD1AC858243B9DA38266381FE90F026C0C1E4110FCDA462E7FE22E40E0EBA756B1BC47F9D517D9C92BAC70AFB22DC4B40F31FDDA1F6E854579239C46461D45C3F3E9CDD172BEB8171CE16F04EB20C10A86DCF364A7757FE948061776CD706F1C8B526637B7F53E83D8CDF7B66549A3FD9F6BDDA83A5C499B268EE08D016592BD5854C53D1CD0C85C1E50188817CF16F1EEAF3AD6A064B08B3894639F49C09102818100DFF1AA38A23CBEB84EE566C46CA055BD0382ACBF50EA55B3422A056CC2C07A7F8849252AC375358894D428C7C36BB54E3CDCBA561C49021C41B782B23455E365C48D07F7D7D12FD94BC5239A9B9C4A303CF1025B0207E07388787D6471C785F7E14DD42A9C7F4588ABF54BE793300856C19ACD231183729908F7706931DEE7F302818100F5DA5DEED020D9251B8FCE707442BFBEC2394CDDEDA79E5844D199F0FC2B02D8A4F930EAB2DD13EB7A83EA75500195B15D78846AD7FC72F9EE57CEDF571EB6B4F284F7BC11905EBD8B1D8E4B8DCD3F8D22C21004327662B71C88D8FD97D8221C47F8E7D5F9A8D23737AD49027F9A580D05A9ACF9FA252C8B173EC96D751A281702818100DCDF3FA247F15DB1EEAEB763383813182F642CF92CD752DB50809D851DB835999F537542EE30A6322587F308C3A771D4CE966D7A0CBFBD431D55DFA3DF966E87AB09E637FE3625D94DB00C63AAE2C5113AAA0246BC84044E2EE597D6FF99687A894EF7D9672CE7E9DAA03ED3120AA7CED978D2A6A9D959A7B27E49F296EB611D0281802941E9FD87A3DB8CE4A12F6DA3B507E485478464C1DB1D3186EAFDC07930E69B60A408D77A08ABAD1AB444864754DCC0150582834397B3DBC969A6E7C800F97C482E943C555E3AE7E80E9FB0822D6D7ACBD87143A30C46E89FBB3F5EDF3A800EEAED144ACE48CC6E43C3AABAE69B0A27B5499223A91CCFEACF8DD3D3B091212502818100C68D284CAC625B8AB21E2ABDC62B418C66F8E707402F763F4087C55BCAE98197C2D9829E3136FCE952A240230CDCB993BF6BDEF1178CEE8E7305CF538F42687B2CB57497187C885DCC38E4988D9FCDFAA5D505A85056FA248D6638E236FFBB68C34A7FC119ED3EC5270317DD39D3B9A6440C4094403466EE23FDE541F01DB2F8}`
	tcPubKey  = `PubKey{3082010A0282010100D7114DDEDB12344F4C7086819B91DCE5FE4DD914306BBAFDDA2EB1080A43F98EC0281DD07B5D1F9229DF38E2C986F8CA6B2AC42E4E8E9F021A38E82A03388B4F2314D58D738BB9D77243CC7D8641C574B4CDC00A1D2720329E1162F1EDC7FDA9BF2F9948E927CCF6C274321473E769AA68CB3669729B19BC1CE022E9BA1E0683F3970330ED7FC7F9E836CD1CDB67B571CD98B82FBEB624948E4BD3CF161250311F2661BB6D0885C3E81B66219CB3407619DF239427214328F1C10CC591C490417DB13CB6B90EC22FEC7B7E4BE8B9B3B75E9820D6446AE7F5690CE0ED048DA5ACE47105A1441F4C765CF345499EBFB55B1FFCFE9CAEF2795E3826075AECF2CED50203010001}`
	tcAddr    = "82bd6fdab3f4141bbd11fec786c7d3cc68b557fc1c96022791d3d34ae038ea1e"
)

func TestLoadRSAKeyUnknownType(t *testing.T) {
	t.Parallel()

	for i := 0; i < 2; i++ {
		testLoadKeyUnknownType(t, i)
	}
}

func testLoadKeyUnknownType(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = LoadRSAPrivKey(struct{}{})
	case 1:
		_ = LoadRSAPubKey(struct{}{})
	}
}

func TestLoadRSAKey(t *testing.T) {
	t.Parallel()

	priv := LoadRSAPrivKey(tcPrivKey)
	if priv == nil {
		t.Error("failed load private key")
		return
	}

	if priv := LoadRSAPrivKey([]byte{123}); priv != nil {
		t.Error("success load invalid private key (bytes)")
		return
	}

	str := string("123")
	if priv := LoadRSAPrivKey(str); priv != nil {
		t.Error("success load invalid private key (string)")
		return
	}

	prefix := cPrivKeyPrefix
	if priv := LoadRSAPrivKey(prefix + str); priv != nil {
		t.Error("success load invalid private key (string+prefix)")
		return
	}

	suffix := cKeySuffix
	if priv := LoadRSAPrivKey(prefix + str + suffix); priv != nil {
		t.Error("success load invalid private key (string+prefix+suffix)")
		return
	}

	pub := LoadRSAPubKey(tcPubKey)
	if pub == nil {
		t.Error("failed load public key")
		return
	}

	if pub := LoadRSAPubKey([]byte{123}); pub != nil {
		t.Error("success load invalid public key (bytes)")
		return
	}

	if pub := LoadRSAPubKey(str); pub != nil {
		t.Error("success load invalid public key (string)")
		return
	}

	prefixPub := cPubKeyPrefix
	if pub := LoadRSAPubKey(prefixPub + str); pub != nil {
		t.Error("success load invalid public key (string+prefix)")
		return
	}

	if pub := LoadRSAPubKey(prefixPub + str + suffix); pub != nil {
		t.Error("success load invalid public key (string+prefix+suffix)")
		return
	}

	if priv.GetPubKey().GetHasher().ToString() != pub.GetHasher().ToString() {
		t.Error("load public key have not relation with private key")
		return
	}

	if err := testRSAConverter(priv, pub); err != nil {
		t.Error(err)
		return
	}
}

func testRSAConverter(priv IPrivKey, pub IPubKey) error {
	if priv.GetSize() != tcKeySize {
		return fmt.Errorf("private key size != tcKeySize")
	}

	if pub.GetSize() != tcKeySize {
		return fmt.Errorf("public key size != tcKeySize")
	}

	if priv.ToString() != tcPrivKey {
		return fmt.Errorf("private key string != tcPrivKey")
	}

	if pub.ToString() != tcPubKey {
		return fmt.Errorf("public key string != tcPrivKey")
	}

	if pub.GetHasher().ToString() != tcAddr {
		return fmt.Errorf("address string != tcAddr")
	}

	return nil
}

func TestRSASign(t *testing.T) {
	t.Parallel()

	var (
		priv = NewRSAPrivKey(1024)
		msg  = []byte("hello, world!")
	)

	pub := priv.GetPubKey()
	sign := priv.SignBytes(msg)

	if !pub.VerifyBytes(msg, sign) {
		t.Error("signature is invalid")
		return
	}

	if pub.VerifyBytes(msg, []byte{123}) {
		t.Error("success verify with invalid signature")
		return
	}

	if pub.VerifyBytes([]byte{123}, msg) {
		t.Error("success verify with invalid message")
		return
	}
}

func TestRSAEncrypt(t *testing.T) {
	t.Parallel()

	var (
		priv = NewRSAPrivKey(1024)
		msg  = []byte("hello, world!")
	)

	pub := priv.GetPubKey()

	if enc := pub.EncryptBytes(random.NewStdPRNG().GetBytes(1 << 10)); enc != nil {
		t.Error("success encrypt message with size > key size")
		return
	}

	emsg := pub.EncryptBytes(msg)

	if !bytes.Equal(msg, priv.DecryptBytes(emsg)) {
		t.Error("decrypted message is invalid")
		return
	}

	if dec := priv.DecryptBytes([]byte{123}); dec != nil {
		t.Error("success decrypt invalid message")
		return
	}
}
