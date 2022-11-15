package config

import (
	"os"
	"testing"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/filesystem"
)

const (
	tcConfigFile = "config_test.txt"
)

const (
	tcNetwork     = "test_network_key"
	tcAddressTCP  = "test_address_tcp"
	tcAddressHTTP = "test_address_http"
)

const (
	tcConfig = `{
	"network": "test_network_key",
	"address": {
		"tcp": "test_address_tcp",
		"http": "test_address_http"
	},
	"connections": [
		"test_connect1",
		"test_connect2",
		"test_connect3"
	],
	"friends": {
		"test_name": "Pub(go-peer/rsa){3082020A0282020100C62F3CFA3D9809EE6DD77EBBFD38BC6796ABA76B795B3C76D3449F0AC808E01EDA8B2B08C58E508C306B2D842A2D317FF6B6D4A13EB76C7BBD5B157B663C3390B227476F4985EF649510D8CCA38FAB9FFCD67916FE73DB77595AB64FBE66D85892708A2DBCA94447A628F183FA6328136FCF158688CB6664EBA91F4C41621741786D50E3286AF9CAB81C101BDB19ACF42E10041CFDA5C6F30ACBBC4251E3D13C0E0781CBDC622E4ED490DD76BBA04D0A9C0012EBDAA77BD9F23183205A9D533C95A6C1FAAD8AB7C3B21FA4C76F7A3FB8EAEB231083ED925C1F71D23671E8C90E460C673A0DCD82ECFA956DF315200554571A99D79EB1E744681B9652389DBA6B9937CE476EBCAC34D02AEACF381DA40469B2F23E4F3DBFD5D8E04031708E46C31E3DC94342298E6F83CF7869C1209ACE2EA04FDB011D0FE265C8D51CF7D90C947160415B3415DFF9D1B16D5A9961F896109223B1408E740C421C6F413FA7B3D7094144DE4A0211DCAF043BC1A9FDE120251CBD654E705795D692A912F0543FF2F13EC733BD1E3AB83B915F95D3540EAA809C1E6E8C248A1EA1AE1D3B29C804F855167F64DA0AB06E5D89080D77D95A6E7199B079925922EA8735DF7654A01B350D67472F25B79DE5FF65B7E9156AEFC8818A1D9216BC4BE527DDC7D88F249B8745CF7DF1610A8237EB4BC1325C64FF47BD34B32CFE59720EC7FB52608D9009C70203010001}"
	},
	"services": {
		"test_service1": "test_address1",
		"test_service2": "test_address2",
		"test_service3": "test_address3"
	}
}`
)

var (
	tgConnects = []string{
		"test_connect1",
		"test_connect2",
		"test_connect3",
	}
	tgPubKeys = map[string]string{
		"test_name": `Pub(go-peer/rsa){3082020A0282020100C62F3CFA3D9809EE6DD77EBBFD38BC6796ABA76B795B3C76D3449F0AC808E01EDA8B2B08C58E508C306B2D842A2D317FF6B6D4A13EB76C7BBD5B157B663C3390B227476F4985EF649510D8CCA38FAB9FFCD67916FE73DB77595AB64FBE66D85892708A2DBCA94447A628F183FA6328136FCF158688CB6664EBA91F4C41621741786D50E3286AF9CAB81C101BDB19ACF42E10041CFDA5C6F30ACBBC4251E3D13C0E0781CBDC622E4ED490DD76BBA04D0A9C0012EBDAA77BD9F23183205A9D533C95A6C1FAAD8AB7C3B21FA4C76F7A3FB8EAEB231083ED925C1F71D23671E8C90E460C673A0DCD82ECFA956DF315200554571A99D79EB1E744681B9652389DBA6B9937CE476EBCAC34D02AEACF381DA40469B2F23E4F3DBFD5D8E04031708E46C31E3DC94342298E6F83CF7869C1209ACE2EA04FDB011D0FE265C8D51CF7D90C947160415B3415DFF9D1B16D5A9961F896109223B1408E740C421C6F413FA7B3D7094144DE4A0211DCAF043BC1A9FDE120251CBD654E705795D692A912F0543FF2F13EC733BD1E3AB83B915F95D3540EAA809C1E6E8C248A1EA1AE1D3B29C804F855167F64DA0AB06E5D89080D77D95A6E7199B079925922EA8735DF7654A01B350D67472F25B79DE5FF65B7E9156AEFC8818A1D9216BC4BE527DDC7D88F249B8745CF7DF1610A8237EB4BC1325C64FF47BD34B32CFE59720EC7FB52608D9009C70203010001}`,
	}
	tgServices = map[string]string{
		"test_service1": "test_address1",
		"test_service2": "test_address2",
		"test_service3": "test_address3",
	}
)

func testConfigDefaultInit(configPath string) {
	filesystem.OpenFile(configPath).Write([]byte(tcConfig))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
	}

	if cfg.Network() != tcNetwork {
		t.Errorf("network is invalid")
	}

	if cfg.Address().TCP() != tcAddressTCP {
		t.Errorf("address_tcp is invalid")
	}

	if cfg.Address().HTTP() != tcAddressHTTP {
		t.Errorf("address_http is invalid")
	}

	for i, v := range cfg.Connections() {
		if v != tgConnects[i] {
			t.Errorf("connection '%d' is invalid", i)
		}
	}

	for k, v := range tgServices {
		v1, ok := cfg.Service(k)
		if !ok {
			t.Errorf("service undefined '%s'", k)
		}
		if v != v1 {
			t.Errorf("service address is invalid '%s'", v1)
		}
	}

	for name, pubStr := range tgPubKeys {
		v1 := cfg.Friends()[name]
		pubKey := asymmetric.LoadRSAPubKey(pubStr)
		if pubKey.Address().String() != v1.Address().String() {
			t.Errorf("public key is invalid '%s'", v1)
		}
	}
}

func TestWrapper(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
	}

	if len(cfg.Friends()) == 0 {
		t.Errorf("list of friends should be is not nil for tests")
		return
	}

	wrapper := NewWrapper(cfg)
	wrapper.Editor().UpdateFriends(nil)

	if len(cfg.Friends()) != 0 {
		t.Errorf("friends is not nil for current config")
		return
	}

	cfg, err = LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
	}

	if len(cfg.Friends()) != 0 {
		t.Errorf("friends is not nil for loaded config")
		return
	}
}
