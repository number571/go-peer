package config

import (
	"fmt"
	"os"
	"testing"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/internal/slices"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcNewNetworkKey = "abc_network_key"
)

var (
	tgNewConnections = []string{"a", "b", "c", "b"}
	tgNewFriends     = map[string]asymmetric.IPubKey{
		"a": asymmetric.LoadRSAPubKey(testutils.TgPubKeys[2]),
		"b": asymmetric.LoadRSAPubKey(testutils.TgPubKeys[3]),
	}

	// diff size of keys: 1024, 4096 bits
	tgNewIncorrect1Friends = map[string]asymmetric.IPubKey{
		"a": asymmetric.LoadRSAPubKey(testutils.TgPubKeys[2]),
		"b": asymmetric.LoadRSAPubKey(`PubKey{3082020A0282020100C62F3CFA3D9809EE6DD77EBBFD38BC6796ABA76B795B3C76D3449F0AC808E01EDA8B2B08C58E508C306B2D842A2D317FF6B6D4A13EB76C7BBD5B157B663C3390B227476F4985EF649510D8CCA38FAB9FFCD67916FE73DB77595AB64FBE66D85892708A2DBCA94447A628F183FA6328136FCF158688CB6664EBA91F4C41621741786D50E3286AF9CAB81C101BDB19ACF42E10041CFDA5C6F30ACBBC4251E3D13C0E0781CBDC622E4ED490DD76BBA04D0A9C0012EBDAA77BD9F23183205A9D533C95A6C1FAAD8AB7C3B21FA4C76F7A3FB8EAEB231083ED925C1F71D23671E8C90E460C673A0DCD82ECFA956DF315200554571A99D79EB1E744681B9652389DBA6B9937CE476EBCAC34D02AEACF381DA40469B2F23E4F3DBFD5D8E04031708E46C31E3DC94342298E6F83CF7869C1209ACE2EA04FDB011D0FE265C8D51CF7D90C947160415B3415DFF9D1B16D5A9961F896109223B1408E740C421C6F413FA7B3D7094144DE4A0211DCAF043BC1A9FDE120251CBD654E705795D692A912F0543FF2F13EC733BD1E3AB83B915F95D3540EAA809C1E6E8C248A1EA1AE1D3B29C804F855167F64DA0AB06E5D89080D77D95A6E7199B079925922EA8735DF7654A01B350D67472F25B79DE5FF65B7E9156AEFC8818A1D9216BC4BE527DDC7D88F249B8745CF7DF1610A8237EB4BC1325C64FF47BD34B32CFE59720EC7FB52608D9009C70203010001}`),
	}

	// duplicated public keys
	tgNewIncorrect2Friends = map[string]asymmetric.IPubKey{
		"a": asymmetric.LoadRSAPubKey(testutils.TgPubKeys[2]),
		"b": asymmetric.LoadRSAPubKey(testutils.TgPubKeys[2]),
	}
)

type tsConfig struct{}

var (
	_ IConfig = &tsConfig{}
)

func (p *tsConfig) GetSettings() IConfigSettings              { return nil }
func (p *tsConfig) GetLogging() logger.ILogging               { return nil }
func (p *tsConfig) GetShare() bool                            { return false }
func (p *tsConfig) GetAddress() IAddress                      { return nil }
func (p *tsConfig) GetNetworkKey() string                     { return "" }
func (p *tsConfig) GetConnections() []string                  { return nil }
func (p *tsConfig) GetFriends() map[string]asymmetric.IPubKey { return nil }
func (p *tsConfig) GetService(_ string) (IService, bool)      { return nil, false }
func (p *tsConfig) GetF2FDisabled() bool                      { return false }

func TestPanicEditor(t *testing.T) {
	t.Parallel()

	for i := 0; i < 2; i++ {
		testPanicEditor(t, i)
	}
}

func testPanicEditor(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = newEditor(nil)
	case 1:
		_ = newEditor(&tsConfig{})
	}
}

func TestEditor(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 4)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	wrapper := NewWrapper(cfg)

	config := wrapper.GetConfig()
	editor := wrapper.GetEditor()

	beforeNetworkKey := config.GetSettings().GetNetworkKey()
	beforeConnections := config.GetConnections()
	beforeFriends := config.GetFriends()

	if err := editor.UpdateNetworkKey(tcNewNetworkKey); err != nil {
		t.Error(err)
		return
	}
	afterNetworkKey := config.GetSettings().GetNetworkKey()
	if beforeNetworkKey == afterNetworkKey {
		t.Error("beforeNetworkKey == afterNetworkKey")
		return
	}
	if afterNetworkKey != tcNewNetworkKey {
		t.Error("afterNetworkKey != tcNewNetworkKey")
		return
	}

	if err := editor.UpdateConnections(tgNewConnections); err != nil {
		t.Error(err)
		return
	}
	afterConnections := config.GetConnections()
	if len(afterConnections) != 3 {
		t.Error("failed deduplicate strings (connections)")
		return
	}
	hasNewConn := false
	for _, ac := range afterConnections {
		if !slices.HasInSlice(beforeConnections, ac) {
			hasNewConn = true
			break
		}
	}
	if !hasNewConn {
		t.Error("beforeConnections == afterConnections")
		return
	}
	for _, nc := range tgNewConnections {
		if !slices.HasInSlice(afterConnections, nc) {
			t.Error("afterConnections != tgNewConnections")
			return
		}
	}

	if err := editor.UpdateFriends(tgNewFriends); err != nil {
		t.Error(err)
		return
	}
	afterFriends := config.GetFriends()
	if len(afterFriends) != 2 {
		t.Error("failed deduplicate public keys (friends)")
		return
	}
	for af := range afterFriends {
		if _, ok := beforeFriends[af]; ok {
			t.Error("beforeFriends == afterFriends")
			return
		}
	}
	for nf := range tgNewFriends {
		if _, ok := afterFriends[nf]; !ok {
			t.Error("afterFriends != tgNewFriends")
			return
		}
	}

	if err := editor.UpdateFriends(tgNewIncorrect2Friends); err == nil {
		t.Error("success update friends with duplicates")
		return
	}
}

func TestIncorrectFilepathEditor(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 5)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	wrapper := NewWrapper(cfg)

	config := wrapper.GetConfig().(*SConfig)
	editor := wrapper.GetEditor()

	config.fFilepath = random.NewStdPRNG().GetString(32)

	if err := editor.UpdateNetworkKey(tcNewNetworkKey); err == nil {
		t.Error("success update network key with incorrect filepath")
		return
	}

	if err := editor.UpdateConnections(tgNewConnections); err == nil {
		t.Error("success update connections with incorrect filepath")
		return
	}

	if err := editor.UpdateFriends(tgNewIncorrect1Friends); err == nil {
		t.Error("success update friends with incorrect key sizes")
		return
	}

	if err := editor.UpdateFriends(tgNewFriends); err == nil {
		t.Error("success update friends with incorrect filepath")
		return
	}
}
