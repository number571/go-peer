package initapp

import (
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	tcTestdataPath           = "./testdata/"
	tcTestdataDirPath        = "./testdata/directory"
	tcPrivKey1024Path        = tcTestdataPath + "priv1024.key"
	tcPasswordPath           = tcTestdataPath + "password.key"
	tcTmpPrivKey512Path      = tcTestdataPath + "tmp_priv512.key"
	tcTmpPasswordPath        = tcTestdataPath + "tmp_password.key"
	tcInvalidPrivKey1024Path = tcTestdataPath + "invalid_priv1024.key"

	tcPassword    = "ee12f9c090cb1903708dea269bdcefa352171611da1c635a0a0694244bf9c049" // nolint: gosec
	tcPrivKey1024 = `PrivKey{3082025E02010002818100D9858A9B2E81BA4543A2548E443936DDC9F2134892315AF511F295A50B869B7A3F6D5F6B774AD2A2139113BAC1CBBDCDE19C3C6152A625A5926C6FA2F0F536A21424520F54D98855A57415B6176170DE9768A988129AC98923D7CB12ACA5548D2D12758B485606FE882C0F3E5169905C9F894120752DADBE858CFA42B400C943020301000102818037F15F676FBB8F8376D48DF894D53E262664EACEB4429B490217A8A2ECE6EE9FAF265AEF119C1DB5EF605579A793D5B9D8774D141EA47A742DC753A2CD63D36BBC6D60B5CD3D64742D794615BB3BDF43B381FE49BF21BD48D79B190F67C6DC0C7C5676BE691137DEB5543B642D36BB359AF3E0590661AEDA0BCAA23679E2B5E1024100E9A4AC8D71A6E72CFE89A17008DA7C3A23776F904578A16F418DE45CC712E839A1589F4169AD853BB53607EB8B43DA9990A95F40800D9ADA58A1664B74B31E11024100EE55F3CCFAC465A07D47145FCBD36BE87083AD0236D2694FA70649A3D3386B6505915BDB64754ECE1F067B5B35A5838381EB4CA68CB1C998AA0FE4B0B9E6AE1302410093F5EC2C7ADFE6A090E559EE183D3CD498A74768870638BDBB36FF7A5DBBB482E291BBF0F1DAA878426EE01F2387AA04FC1EB6AAA32D7A7672106C36B6C5C3F10241009D965D86AA4493C1C333ED67CDF8B43FD3AD6D06AAC30378F442370CC88B648F3E5837795FFA24AA2B5F78CEFD30BC3D86F8D30CC8B881489D21B71F973BCCDB024100A47F526668703511E592A2F450B15286AF1A88DFF3F9688E749E4389A73F2B8100370727E06D9260B7EE71B87F6D16EC389D7309E4CDAB7A1F9F59DED0C2688B}`
)

func testDeleteFile(f string) {
	os.RemoveAll(f)
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestPassword(t *testing.T) {
	t.Parallel()

	testDeleteFile(tcTmpPasswordPath)
	defer testDeleteFile(tcTmpPasswordPath)

	password, err := GetPassword(tcPasswordPath)
	if err != nil {
		t.Error(err)
		return
	}

	if password != tcPassword {
		t.Error("diff passwords")
		return
	}

	if _, err := GetPassword("./random/not_exist/path/57199u140291724y121291d1/priv.key"); err == nil {
		t.Error("success get private key with not exist directory")
		return
	}

	tmpPassword, err := GetPassword(tcTmpPasswordPath)
	if err != nil {
		t.Error(err)
		return
	}
	tmpPasswordX, err := GetPassword(tcTmpPasswordPath)
	if err != nil {
		t.Error(err)
		return
	}
	if tmpPassword != tmpPasswordX {
		t.Error("diff tmp passwords")
		return
	}
}

func TestGetPrivKey(t *testing.T) {
	t.Parallel()

	testDeleteFile(tcTmpPrivKey512Path)
	defer testDeleteFile(tcTmpPrivKey512Path)

	privKey, err := GetPrivKey(tcPrivKey1024Path, 1024)
	if err != nil {
		t.Error(err)
		return
	}
	if privKey.ToString() != asymmetric.LoadRSAPrivKey(tcPrivKey1024).ToString() {
		t.Error("diff private keys")
		return
	}

	if _, err := GetPrivKey(tcInvalidPrivKey1024Path, 1024); err == nil {
		t.Error("success get invalid private key")
		return
	}
	if _, err := GetPrivKey(tcPrivKey1024Path, 2048); err == nil {
		t.Error("success get private key with diff size")
		return
	}
	if _, err := GetPrivKey("./random/not_exist/path/57199u140291724y121291d1/priv.key", 512); err == nil {
		t.Error("success get private key with not exist directory")
		return
	}
	if _, err := GetPrivKey(tcTestdataDirPath, 512); err == nil {
		t.Error("success get private key as directory")
		return
	}

	tmpPrivKey, err := GetPrivKey(tcTmpPrivKey512Path, 512)
	if err != nil {
		t.Error(err)
		return
	}
	tmpPrivKeyX, err := GetPrivKey(tcTmpPrivKey512Path, 512)
	if err != nil {
		t.Error(err)
		return
	}
	if tmpPrivKey.ToString() != tmpPrivKeyX.ToString() {
		t.Error("diff tmp private keys")
		return
	}
}
