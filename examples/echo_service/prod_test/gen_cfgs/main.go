package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	cMessageSizeBytes = (8 << 10)
)

type sNodeHLS struct {
	fPath     string
	fConfig   string
	fMakefile string
}

var (
	gNumOtherNodes  int
	gNetworkKey     string
	gListOfConnects []string
)

func initRun() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("len args != 3")
	}

	switch os.Args[1] {
	case "prod_1":
		gNetworkKey = cNetworkKey_1
		gListOfConnects = gListOfConnects_1
	case "prod_2":
		gNetworkKey = cNetworkKey_2
		gListOfConnects = gListOfConnects_2
	default:
		return fmt.Errorf("unknown param")
	}

	numNodes, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	gNumOtherNodes = numNodes
	return nil
}

func main() {
	if err := initRun(); err != nil {
		panic(err)
	}

	strConnectsBuilder := strings.Builder{}
	strConnectsBuilder.Grow(len(gListOfConnects))
	for _, c := range gListOfConnects {
		strConnectsBuilder.WriteString("\n  - " + c)
	}
	strConnects := strConnectsBuilder.String()

	firstConnect := gListOfConnects[0]
	secondConnect := gListOfConnects[len(gListOfConnects)-1]

	nodes := make([]*sNodeHLS, 0, 2+gNumOtherNodes) // 2 = recv+send
	nodes = append(nodes, initRecvNode(firstConnect), initSendNode(secondConnect))

	os.Mkdir("other_nodes", 0744)
	for i := 1; i <= gNumOtherNodes; i++ {
		node := initOtherNode(i, strConnects)
		nodes = append(nodes, node)
		os.Mkdir(node.fPath, 0744)
	}

	for _, node := range nodes {
		os.WriteFile(fmt.Sprintf("%s/hls.yml", node.fPath), []byte(node.fConfig), 0644)
		if node.fMakefile != "" {
			os.WriteFile(fmt.Sprintf("%s/Makefile", node.fPath), []byte(node.fMakefile), 0644)
		}
	}
}

func initRecvNode(pConnects string) *sNodeHLS {
	return &sNodeHLS{
		fPath: "recv_hls",
		fConfig: fmt.Sprintf(`settings:
  message_size_bytes: %d
  work_size_bits: 22
  key_size_bits: 4096
  queue_period_ms: 5000
  limit_void_size_bytes: 4096
  network_key: %s
logging: 
  - info
  - warn
  - erro
services:
  hidden-echo-service: 
    host: localhost:8080
connections: 
  - %s
friends:
  Alice: PubKey{3082020A0282020100C17B6FA53983050B0339A0AB60D20A8A5FF5F8210564464C45CD2FAC2F266E8DDBA3B36C6F356AE57D1A71EED7B612C4CBC808557E4FCBAF6EDCFCECE37494144F09D65C7533109CE2F9B9B31D754453CA636A4463594F2C38303AE1B7BFFE738AC57805C782193B4854FF3F3FACA2C6BF9F75428DF6C583FBC29614C0B3329DF50F7B6399E1CC1F12BED77F29F885D7137ADFADE74A43451BB97A32F2301BE8EA866AFF34D6C7ED7FF1FAEA11FFB5B1034602B67E7918E42CA3D20E3E68AA700BE1B55A78C73A1D60D0A3DED3A6E5778C0BA68BAB9C345462131B9DC554D1A189066D649D7E167621815AB5B93905582BF19C28BCA6018E0CD205702968885E92A3B1E3DB37A25AC26FA4D2A47FF024ECD401F79FA353FEF2E4C2183C44D1D44B44938D32D8DBEDDAF5C87D042E4E9DAD671BE9C10DD8B3FE0A7C29AFE20843FE268C6A8F14949A04FF25A3EEE1EBE0027A99CE1C4DC561697297EA9FD9E23CF2E190B58CA385B66A235290A23CBB3856108EFFDD775601B3DE92C06C9EA2695C2D25D7897FD9D43C1AE10016E51C46C67F19AC84CD25F47DE2962A48030BCD8A0F14FFE4135A2893F62AC3E15CC61EC2E4ACADE0736C9A8DBC17D439248C42C5C0C6E08612414170FBE5AA6B52AE64E4CCDAE6FD3066BED5C200E07DBB0167D74A9FAD263AF253DFA870F44407F8EF3D9F12B8D910C4D803AD82ABA136F93F0203010001}
`,
			cMessageSizeBytes,
			gNetworkKey,
			pConnects,
		),
	}
}

func initSendNode(pConnects string) *sNodeHLS {
	return &sNodeHLS{
		fPath: "send_hls",
		fConfig: fmt.Sprintf(`settings:
  message_size_bytes: %d
  work_size_bits: 22
  key_size_bits: 4096
  queue_period_ms: 5000
  limit_void_size_bytes: 4096
  network_key: "%s"
logging: 
  - info
  - warn
  - erro
address:
  http: localhost:7572
connections:
  - %s
friends:
  Bob: PubKey{3082020A0282020100B752D35E81F4AEEC1A9C42EDED16E8924DD4D359663611DE2DCCE1A9611704A697B26254DD2AFA974A61A2CF94FAD016450FEF22F218CA970BFE41E6340CE3ABCBEE123E35A9DCDA6D23738DAC46AF8AC57902DDE7F41A03EB00A4818137E1BF4DFAE1EEDF8BB9E4363C15FD1C2278D86F2535BC3F395BE9A6CD690A5C852E6C35D6184BE7B9062AEE2AFC1A5AC81E7D21B7252A56C62BB5AC0BBAD36C7A4907C868704985E1754BAA3E8315E775A51B7BDC7ACB0D0675D29513D78CB05AB6119D3CA0A810A41F78150E3C5D9ACAFBE1533FC3533DECEC14387BF7478F6E229EB4CC312DC22436F4DB0D4CC308FB6EEA612F2F9E00239DE7902DE15889EE71370147C9696A5E7B022947ABB8AFBBC64F7840BED4CE69592CAF4085A1074475E365ED015048C89AE717BC259C42510F15F31DA3F9302EAD8F263B43D14886B2335A245C00871C041CBB683F1F047573F789673F9B11B6E6714C2A3360244757BB220C7952C6D3D9D65AA47511A63E2A59706B7A70846C930DCFB3D8CAFB3BD6F687CACF5A708692C26B363C80C460F54E59912D41D9BB359698051ABC049A0D0CFD7F23DC97DA940B1EDEAC6B84B194C8F8A56A46CE69EE7A0AEAA11C99508A368E64D27756AD0BA7146A6ADA3D5FA237B3B4EDDC84B71C27DE3A9F26A42197791C7DC09E2D7C4A7D8FCDC8F9A5D4983BB278FCE9513B1486D18F8560C3F31CC70203010001}
`,
			cMessageSizeBytes,
			gNetworkKey,
			pConnects,
		),
	}
}

func initOtherNode(pI int, pConnects string) *sNodeHLS {
	return &sNodeHLS{
		fPath: fmt.Sprintf("other_nodes/other_hls%d", pI),
		fMakefile: fmt.Sprintf(`
GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_other_hls%[1]d &
clean:
	pkill -15 prog_other_hls%[1]d || true
	rm -rf prog_other_hls%[1]d hls.db
`,
			pI),
		fConfig: fmt.Sprintf(`settings:
  message_size_bytes: %d
  work_size_bits: 22
  key_size_bits: 4096
  queue_period_ms: 5000
  limit_void_size_bytes: 4096
  network_key: %s
connections: %s
`,
			cMessageSizeBytes,
			gNetworkKey,
			pConnects,
		),
	}
}
