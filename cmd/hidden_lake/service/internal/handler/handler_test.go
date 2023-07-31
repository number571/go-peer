package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcServiceAddressInHLS = "hidden-echo-service"
	tcPathDB              = "database_test.db"
	tcPathConfig          = "config_test.cfg"
)

const (
	tcConfig = `{
	"settings": {
		"message_size_bytes": 8192,
		"work_size_bits": 20,
		"key_size_bits": 4096,
		"queue_period_ms": 1000,
		"limit_void_size_bytes": 4096
	},
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
		"test_name1": "PubKey(go-peer/rsa){3082020A0282020100C62F3CFA3D9809EE6DD77EBBFD38BC6796ABA76B795B3C76D3449F0AC808E01EDA8B2B08C58E508C306B2D842A2D317FF6B6D4A13EB76C7BBD5B157B663C3390B227476F4985EF649510D8CCA38FAB9FFCD67916FE73DB77595AB64FBE66D85892708A2DBCA94447A628F183FA6328136FCF158688CB6664EBA91F4C41621741786D50E3286AF9CAB81C101BDB19ACF42E10041CFDA5C6F30ACBBC4251E3D13C0E0781CBDC622E4ED490DD76BBA04D0A9C0012EBDAA77BD9F23183205A9D533C95A6C1FAAD8AB7C3B21FA4C76F7A3FB8EAEB231083ED925C1F71D23671E8C90E460C673A0DCD82ECFA956DF315200554571A99D79EB1E744681B9652389DBA6B9937CE476EBCAC34D02AEACF381DA40469B2F23E4F3DBFD5D8E04031708E46C31E3DC94342298E6F83CF7869C1209ACE2EA04FDB011D0FE265C8D51CF7D90C947160415B3415DFF9D1B16D5A9961F896109223B1408E740C421C6F413FA7B3D7094144DE4A0211DCAF043BC1A9FDE120251CBD654E705795D692A912F0543FF2F13EC733BD1E3AB83B915F95D3540EAA809C1E6E8C248A1EA1AE1D3B29C804F855167F64DA0AB06E5D89080D77D95A6E7199B079925922EA8735DF7654A01B350D67472F25B79DE5FF65B7E9156AEFC8818A1D9216BC4BE527DDC7D88F249B8745CF7DF1610A8237EB4BC1325C64FF47BD34B32CFE59720EC7FB52608D9009C70203010001}",
		"test_name2": "PubKey(go-peer/rsa){3082020A0282020100B752D35E81F4AEEC1A9C42EDED16E8924DD4D359663611DE2DCCE1A9611704A697B26254DD2AFA974A61A2CF94FAD016450FEF22F218CA970BFE41E6340CE3ABCBEE123E35A9DCDA6D23738DAC46AF8AC57902DDE7F41A03EB00A4818137E1BF4DFAE1EEDF8BB9E4363C15FD1C2278D86F2535BC3F395BE9A6CD690A5C852E6C35D6184BE7B9062AEE2AFC1A5AC81E7D21B7252A56C62BB5AC0BBAD36C7A4907C868704985E1754BAA3E8315E775A51B7BDC7ACB0D0675D29513D78CB05AB6119D3CA0A810A41F78150E3C5D9ACAFBE1533FC3533DECEC14387BF7478F6E229EB4CC312DC22436F4DB0D4CC308FB6EEA612F2F9E00239DE7902DE15889EE71370147C9696A5E7B022947ABB8AFBBC64F7840BED4CE69592CAF4085A1074475E365ED015048C89AE717BC259C42510F15F31DA3F9302EAD8F263B43D14886B2335A245C00871C041CBB683F1F047573F789673F9B11B6E6714C2A3360244757BB220C7952C6D3D9D65AA47511A63E2A59706B7A70846C930DCFB3D8CAFB3BD6F687CACF5A708692C26B363C80C460F54E59912D41D9BB359698051ABC049A0D0CFD7F23DC97DA940B1EDEAC6B84B194C8F8A56A46CE69EE7A0AEAA11C99508A368E64D27756AD0BA7146A6ADA3D5FA237B3B4EDDC84B71C27DE3A9F26A42197791C7DC09E2D7C4A7D8FCDC8F9A5D4983BB278FCE9513B1486D18F8560C3F31CC70203010001}",
		"test_recvr": "PubKey(go-peer/rsa){3082020A0282020100DE72151AC56C26AC83F1375D141FC9E00810EEF1C6BC5A329FFB2796A930A1824C2A6A5B3C46019BC8CACDA53EC2FF78C983388BAD4670BF1C9609F15942880530CB274A2F0BDC81EDE6F4C7A55FDD3B56E485DB4755C346B53A307F0FFE150822E0CA951902699A688A24DC3182CFC14E9CA9DF2D61A0671CBF46B21DE45A52A4D131C5C137CF43910F2BE89E0964330798D2DEEDA00410CAAED900EA97740B5C32B8B0A93EC8B3ED31310ECE82EA9B5F3E17C8B879A4901D6C90F9066FC8F7974B61E47B4BA6BF56549B6408F17C014965A5C5F2B50C77FF13DDC54CC2BA1032DB40E69C9C55F2104DD1A44B21C6663CCB335A6B6008787D2E3E1F3E834368E4CA98CBF142EBC3435E105C8F853FF0BBB3C6322B5B4061F8EDC800885883D73AF0517EFC672243C690F6F149C1E1E93925249BE9E7F2A7482D49FF2F36DB5C93F25B15C081D393B09D7D71391615068C8B47E9114A582E6F928FE529F3A14CE9A3EE156768501126625DB05030E2DC466DD048D57DFB2A9DF98EF659F9DEC3AFA179D50BA4FA8103398D21D8325FDD05572B4E2E76042520E3F7E34FD7B0E14767AC4EE18B480D808E3DD5F3B1C4DBDE4A27D9FD5AFAA45438A858082873E9184673FDBC1BA0CEAF9E4AD3EF0178EE84659995188659D7498E2047E85F1E22F5F522C061C93315C4D0A7F9A1C56B65379877979FA8F97F7B0B57BA922E47B30203010001}"
	},
	"services": {
		"test_service1": "test_address1",
		"test_service2": "test_address2",
		"test_service3": "test_address3"
	}
}`
)

func testStartServerHTTP(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testEchoPage)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testEchoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		FMessage string `json:"message"`
	}

	var resp struct {
		FEcho  string `json:"echo"`
		FError int    `json:"error"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.FError = 1
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.FEcho = req.FMessage
	json.NewEncoder(w).Encode(resp)
}

func testAllCreate(cfgPath, dbPath, srvAddr string) (config.IWrapper, anonymity.INode, *http.Server) {
	wcfg := testNewWrapper(cfgPath)
	node := testRunNewNode(dbPath, "")
	srvc := testRunService(wcfg, node, srvAddr)
	time.Sleep(200 * time.Millisecond)
	return wcfg, node, srvc
}

func testAllFree(node anonymity.INode, srv *http.Server) {
	defer func() {
		os.RemoveAll(tcPathDB)
		os.RemoveAll(tcPathConfig)
	}()
	types.StopAll([]types.ICommand{
		node,
		node.GetNetworkNode(),
	})
	types.CloseAll([]types.ICloser{
		srv,
		node.GetWrapperDB(),
	})
}

func testRunService(wcfg config.IWrapper, node anonymity.INode, addr string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, HandleConfigConnectsAPI(wcfg, node))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, HandleConfigFriendsAPI(wcfg, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, HandleNetworkOnlineAPI(node))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, HandleNetworkRequestAPI(wcfg, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkMessagePath, HandleNetworkMessageAPI(node))
	mux.HandleFunc(pkg_settings.CHandleNodeKeyPath, HandleNodeKeyAPI(wcfg, node))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testNewWrapper(cfgPath string) config.IWrapper {
	filesystem.OpenFile(cfgPath).Write([]byte(tcConfig))
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	return config.NewWrapper(cfg)
}

func testRunNewNode(dbPath, addr string) anonymity.INode {
	node := testNewNode(dbPath, addr).HandleFunc(pkg_settings.CServiceMask, nil)
	if err := node.Run(); err != nil {
		return nil
	}
	return node
}

func testNewNode(dbPath, addr string) anonymity.INode {
	db, err := database.NewKeyValueDB(
		storage.NewSettings(&storage.SSettings{
			FPath:      dbPath,
			FHashing:   true,
			FCipherKey: []byte("CIPHER"),
		}),
	)
	if err != nil {
		return nil
	}
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   "TEST",
			FRetryEnqueue:  0,
			FNetworkMask:   1,
			FFetchTimeWait: time.Minute,
		}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		anonymity.NewWrapperDB().Set(db),
		testNewNetworkNode(addr),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: testutils.TCQueueCapacity,
				FPoolCapacity: testutils.TCQueueCapacity,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FWorkSizeBits:     testutils.TCWorkSize,
					FMessageSizeBytes: testutils.TCMessageSize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
			),
		),
		asymmetric.NewListPubKeys(),
	)
	return node
}

func testNewNetworkNode(addr string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      addr,
			FCapacity:     testutils.TCCapacity,
			FMaxConnects:  testutils.TCMaxConnects,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
				FFetchTimeWait:    1, // not used
			}),
		}),
	)
}
