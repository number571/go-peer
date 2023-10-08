# HLS

> Hidden Lake Service

<img src="_images/hls_logo.png" alt="hls_logo.png"/>

The `Hidden Lake Service` is the core of an anonymous network with theoretically provable anonymity. HLS is based on the `fifth^ stage` of anonymity and is an implementation of an `abstract` anonymous network based on `queues`. It is a `peer-to-peer` network communication with trusted `friend-to-friend` participants. All transmitted and received messages are in the form of `end-to-end` encryption.

Features / Anonymity networks |  Queue-networks (5^stage)               |  Entropy-networks (6stage)              |  DC-networks (1^stage)
:-----------------------------:|:-----------------------------:|:------------------------------:|:------------------------------:
Theoretical provability  |  +  |  +  |  +
Ease of software implementation  |  +  |  -  |  -
Polymorphism of information  |  -  |  +  |  +
Static communication delay  |  +  |  -  |  +
Network scales easily  |  -  |  -  |  -

A feature of HLS (compared to many other anonymous networks) is its easy adaptation to a hostile centralized environment. Anonymity can be restored literally from one node in the network, even if it is the only point of failure.

> More information about HLS in the [habr.com/ru/post/696504](https://habr.com/ru/post/696504/ "Habr HLS")

## How it works

Each network participant sets a message generation period for himself (the period can be a network constant for all system participants). When one cycle of the period ends and the next begins, each participant sends his encrypted message to all his connections (those in turn to all of their own, etc.). If there is no true message to send, then a pseudo message is generated (filled with random bytes) that looks like a normal encrypted one. The period property ensures the anonymity of the sender.

<p align="center"><img src="_images/hls_queue.jpg" alt="hls_queue.jpg"/></p>
<p align="center">Figure 1. Queue and message generation in HLS.</p>

Since the encrypted message does not disclose the recipient in any way, each network participant tries to decrypt the message with his private key. The true recipient is only the one who can decrypt the message. At the same time, the true recipient acts according to the protocol and further distributes the received packet, even knowing the meaninglessness of the subsequent dispatch. This property makes it impossible to determine the recipient.

> Simple example of the `client` package (encrypt/decrypt functions) in the directory [github.com/number571/go-peer/pkg/client/_examples](https://github.com/number571/go-peer/tree/master/pkg/client/_examples "Package client");

<p align="center"><img src="_images/hls_view.jpg" alt="hls_view.jpg"/></p>
<p align="center">Figure 2. Two participants are constantly generating messages for their periods on the network. It is impossible to determine their real activity.</p>

Data exchange between network participants is carried out using application services. HLS has a dual role: 1) packages traffic from pure to anonymizing and vice versa; 2) converts external traffic to internal and vice versa. The second property is the redirection of traffic from the network to the local service and back.

<p align="center"><img src="_images/hls_service.jpg" alt="hls_service.jpg"/></p>
<p align="center">Figure 3. Interaction of third-party services with the traffic anonymization service.</p>

As shown in the figure above, HLS acts as an anonymizer and handlers of incoming and outgoing traffic. The remaining parts in the form of applications and services depend on third-party components (as an example, `HLM`).

###  More details in the works 

1. [Theory of the structure of hidden systems](https://github.com/number571/go-peer/blob/master/docs/theory_of_the_structure_of_hidden_systems.pdf "TotSoHS")
2. [Monolithic cryptographic protocol](https://github.com/number571/go-peer/blob/master/docs/monolithic_cryptographic_protocol.pdf "MCP")
3. [Abstract anonymous networks](https://github.com/number571/go-peer/blob/master/docs/abstract_anonymous_networks.pdf "AAN")
4. [Decentralized key exchange protocol](https://github.com/number571/go-peer/blob/master/docs/decentralized_key_exchange_protocol.pdf "DKEP")

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Minimum system requirements

1. Processor: `1x2.2GHz` (more than two cores per processor are recommended)
2. Memory: `0.5GB RAM` (~250MB of memory is consumed)
3. Storage: `5Gib available space` (the size of hashes per year from one node)

## Build and run

Default build and run

```bash 
$ cd ./cmd/hidden_lake/service
$ make build # create hls, hls_[arch=amd64,arm64]_[os=linux,windows,darwin] and copy to ./bin
$ make run # run ./bin/hls

> [INFO] 2023/06/03 14:32:40 HLS is running...
> [INFO] 2023/06/03 14:32:42 service=HLS type=BRDCS hash=43A5E9C5...BA73DF43 addr=211494E4...EEA12BBC proof=0000000002256145 conn=127.0.0.1:
> [INFO] 2023/06/03 14:32:47 service=HLS type=BRDCS hash=EFDDC1D4...C47588AD addr=211494E4...EEA12BBC proof=0000000000090086 conn=127.0.0.1:
> [INFO] 2023/06/03 14:32:52 service=HLS type=BRDCS hash=8549E257...EDEB2748 addr=211494E4...EEA12BBC proof=0000000000634328 conn=127.0.0.1:
> ...
```

Service was running with random private key. Open ports `9571` (TCP, traffic) and `9572` (HTTP, interface).
Creates `./hls.cfg` or `./_mounted/hls.cfg` (docker) and `./hls.db` or `./_mounted/hls.db` (docker) files. 
The file `hls.db` stores hashes of sent/received messages.

Default config `hls.cfg`

```json
{
	"settings": {
		"message_size_bytes": 8192,
		"work_size_bits": 20,
		"key_size_bits": 4096,
		"queue_period_ms": 5000,
		"limit_void_size_bytes": 4096
	},
	"logging": [
		"info",
		"warn",
		"erro"
	],
	"address": {
		"tcp": "127.0.0.1:9571",
		"http": "127.0.0.1:9572"
	},
	"services": {
		"go-peer/hidden-lake-messenger": "127.0.0.1:9592"
	}
}
```

If service works not in docker's environment than need rewrite connection host in `hls.cfg` file from `messenger`to IP address (example: `127.0.0.1:9592` for local network).

Build and run with docker

```bash 
$ cd ./cmd/hidden_lake/service
$ make docker-build 
$ make docker-run

> [INFO] 2023/06/03 07:36:49 HLS is running...
> [INFO] 2023/06/03 07:36:51 service=HLS type=BRDCS hash=AF90439F...9F29A036 addr=BB58A8A2...B64D62C2 proof=0000000000479155 conn=127.0.0.1:
> [INFO] 2023/06/03 07:36:56 service=HLS type=BRDCS hash=2C4CE60A...E55BF9C4 addr=BB58A8A2...B64D62C2 proof=0000000000521434 conn=127.0.0.1:
> [INFO] 2023/06/03 07:37:01 service=HLS type=BRDCS hash=A9285F98...F96DB93D addr=BB58A8A2...B64D62C2 proof=0000000001256786 conn=127.0.0.1:
> ...
```

## Example

There are three nodes in the network `send_hls`, `recv_hls` and `middle_hls`. The `send_his` and `recv_hls` nodes connects to `middle_hls`. As a result, a link of the form `send_his <-> middle_hls <-> recv_hls` is created. Due to the specifics of HLS, the centralized `middle_hls` node does not violate the security and anonymity of the `send_hls` and `recv_hls` subjects in any way. All nodes, including the `middle_hls` node, set periods and adhere to the protocol of constant message generation.

The `recv_hls` node contains its `echo_service`, which performs the role of redirecting the request body back to the client as a response. Access to this service is carried out by its alias `hidden-echo-service`, put forward by the recv_hls node.

```go
...
// handle: "/echo"
// return format: {"echo":string,"return":int}
func echoPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response(w, 2, "failed: incorrect method")
		return
	}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		response(w, 3, "failed: read body")
		return
	}
	response(w, 1, string(res))
}
...
```

Identification between `recv_hls` and `send_hls` nodes is performed using public keys. This is the main method of identification and routing in the HLS network. IP addresses are only needed to connect to such a network and no more. Requests and responses structure are HEX encoded.

Structure of request. The body `hello, world!` is encoded base64.
```bash
JSON_DATA='{
        "method":"POST",
        "host":"hidden-echo-service",
        "path":"/echo",
        "head":{
            "Accept": "application/json"
        },
        "body":"aGVsbG8sIHdvcmxkIQ=="
}';
```

Request format
```bash
PUSH_FORMAT='{
        "receiver":"Alice",
        "hex_data":"'$(str2hex "$JSON_DATA")'"
}';
```

Build and run nodes
```bash
$ cd examples/echo_service
$ make
```

Logs from `middle_hls` node. When sending requests and receiving responses, `middle_hls` does not see the action. For him, all actions and moments of inaction are equivalent.

<p align="center"><img src="_images/hls_logger.gif" alt="hls_logger.gif"/></p>
<p align="center">Figure 4. Output of all actions and all received traffic from the middle_hls node.</p>

Send request
```bash
$ ./request.sh
```

Get response
```bash
HTTP/1.1 200 OK
Date: Mon, 22 May 2023 18:18:34 GMT
Content-Length: 113
Content-Type: text/plain; charset=utf-8

{"code":200,"head":{"Content-Type":"application/json"},"body":"eyJlY2hvIjoiaGVsbG8sIHdvcmxkISIsInJldHVybiI6MX0K"}
Request took 8 seconds
```

Return code 200 is HTTP code = StatusOK. Decode base64 response body
```bash
echo "eyJlY2hvIjoiaGVsbG8sIHdvcmxkISIsInJldHVybiI6MX0K" | base64 -d
> {"echo":"hello, world!","return":1}
```

<p align="center"><img src="_images/hls_request.gif" alt="hls_request.gif"/></p>
<p align="center">Figure 5. Example of running HLS with internal service.</p>

Also you can run example with docker-compose. In this example, all nodes have logging enabled
```bash
$ cd examples/echo_service/_docker/default
$ make
```

> Simple examples of the `anonymity` package in the directory [github.com/number571/go-peer/pkg/network/anonymity/_examples](https://github.com/number571/go-peer/tree/master/pkg/network/anonymity/_examples "Package anonymity");

## Cryptographic algorithms and functions

1. AES-256-CFB (Data encryption)
2. RSA-4096-OAEP (Key encryption)
3. RSA-4096-PSS (Hash signing)
4. SHA-256 (Data hashing)
5. HMAC-SHA-256 (Network hashing)
6. PoW-20 (Hash proof)

## Config structure

```
"logging"      Enable loggins in/out actions in the network
"address.tcp"  Connection address for anonymity network, may be void
"address.http" Connection address for API functions
"services"     Map with redirects requests from network to services
"network_key"  A network key created to encapsulate connections
"connections"  Connection addresses of the another nodes in network
"friends"      Friend addresses for send or receive messages over network
```

```json
{
	"settings": {
		"message_size_bytes": 8192,
		"work_size_bits": 20,
		"key_size_bits": 4096,
		"queue_period_ms": 5000,
		"limit_void_size_bytes": 4096
	},
	"logging": ["info", "warn", "erro"],
	"address": {
		"tcp": "localhost:9571",
		"http": "localhost:9572"
	},
	"services": {
		"hidden-default-service": "localhost:8080"
	},
	"network_key": "hls-network-key",
	"connections": [
		"localhost:8571"
	],
	"friends": {
		"alias-name": "PubKey(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}"
	}
}
```

## Request structure in HLS for internal services

```
Need encode this json to hex format
and put result to "hex_data" HLS API

"body" is base64 string
```

```json
{
	"method":"GET",
	"host":"hidden-default-service",
	"path":"/",
	"head":{
		"Accept":"application/json"
	},
	"body":"aGVsbG8sIHdvcmxkIQ=="
}
```

## Response structure from HLS API

```
"result" is string
"return" is int; 1 = success
```

```json
{
	"result":"go-peer/hidden-lake-service",
	"return":1
}
```

## HLS API

```
1. GET/POST/DELETE /api/config/connects
2. GET/POST/DELETE /api/config/friends
3. GET/DELETE      /api/network/online
4. POST/PUT        /api/network/request
5. GET/POST        /api/network/key
6. GET/POST/DELETE /api/node/key
```

### 1. /api/config/connects

#### 1.1. GET Request

```bash
curl -i -X GET -H 'Accept: application/json' http://localhost:9572/api/config/connects
```

#### 1.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 07 Aug 2023 00:21:25 GMT
Content-Length: 35
```

```json
["localhost:9581","localhost:8888"]
```

#### 1.2. POST Request

```bash
curl -i -X POST -H 'Accept: application/json' http://localhost:9572/api/config/connects --data 'localhost:8888'
```

#### 1.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Mon, 07 Aug 2023 00:21:17 GMT
Content-Length: 27

success: update connections
```

#### 1.3. DELETE Request

```bash
curl -i -X DELETE -H 'Accept: application/json' http://localhost:9572/api/config/connects --data 'localhost:8888'
```

#### 1.3. DELETE Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Mon, 07 Aug 2023 00:21:46 GMT
Content-Length: 26

success: delete connection
```

### 2. /api/config/friends

#### 2.1. GET Request

```bash
curl -i -X GET -H 'Accept: application/json' http://localhost:9572/api/config/friends
```

#### 2.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 07 Aug 2023 00:22:38 GMT
Transfer-Encoding: chunked
```

```json
[{"alias_name":"Bob","public_key":"PubKey(go-peer/rsa){3082020A0282020100B752D35E81F4AEEC1A9C42EDED16E8924DD4D359663611DE2DCCE1A9611704A697B26254DD2AFA974A61A2CF94FAD016450FEF22F218CA970BFE41E6340CE3ABCBEE123E35A9DCDA6D23738DAC46AF8AC57902DDE7F41A03EB00A4818137E1BF4DFAE1EEDF8BB9E4363C15FD1C2278D86F2535BC3F395BE9A6CD690A5C852E6C35D6184BE7B9062AEE2AFC1A5AC81E7D21B7252A56C62BB5AC0BBAD36C7A4907C868704985E1754BAA3E8315E775A51B7BDC7ACB0D0675D29513D78CB05AB6119D3CA0A810A41F78150E3C5D9ACAFBE1533FC3533DECEC14387BF7478F6E229EB4CC312DC22436F4DB0D4CC308FB6EEA612F2F9E00239DE7902DE15889EE71370147C9696A5E7B022947ABB8AFBBC64F7840BED4CE69592CAF4085A1074475E365ED015048C89AE717BC259C42510F15F31DA3F9302EAD8F263B43D14886B2335A245C00871C041CBB683F1F047573F789673F9B11B6E6714C2A3360244757BB220C7952C6D3D9D65AA47511A63E2A59706B7A70846C930DCFB3D8CAFB3BD6F687CACF5A708692C26B363C80C460F54E59912D41D9BB359698051ABC049A0D0CFD7F23DC97DA940B1EDEAC6B84B194C8F8A56A46CE69EE7A0AEAA11C99508A368E64D27756AD0BA7146A6ADA3D5FA237B3B4EDDC84B71C27DE3A9F26A42197791C7DC09E2D7C4A7D8FCDC8F9A5D4983BB278FCE9513B1486D18F8560C3F31CC70203010001}"},{"alias_name":"Eve","public_key":"PubKey(go-peer/rsa){3082020A0282020100C971936CE7F037D60E72613552F9E5DDAA7367CF414E1EC97B1C622E9DAC66DEE488048CECB51ACA082EF1EA1F6DF05FEB595B8A075C012634EEFDA62905717D2FC4DCBEDB824F0E92015E6124C81FD9D6B5E0EE13C685F6E226CDA5646DC2BE32D2FDDE486B0B15F4B7455CD1311F604A822C321B304ECBD599D2D7B4A8FB380F38AAEBCC2D1176E1D2BA85F38E25B7879DD61A8C290F55BAB4502221F23DDF6F75B5B3CA631D63B736FD7B7E6F8F9A82F55DE5B673862F0F324F4F911502810477E820946057F951B57E44EC79525BD10B472D05F57A7CAAB835AE55E71129CF9B1CC54175989E1BE86697F9A4C560D09C179CD332E05550F169DEC6318D4E8172F009DC82837A418E454A75E4CAE5A098161099BEA499FFD56E98433ABECAE1961A864388D355EF29C02DC1DEE315C03D16DA6687B6AD67D544A20E541ADB450D1CC57869EB21D3B53368CB716DCBAF18E625A4A68081651C2E5AEA28549F141DBBBB1F500EE970303DE1BC82098B130D202234322D7C1E8A71D71F016E10ACF523303DB48BD1DA9B1D2E623012557CFB81176F2195872F244E6149FE03395951AEE6F90B4808A88796875264A4FF177504D5139EE4729D9603FC3B0F448E0F3E95865CD5234A169DEF8EE07A067DED78E782A534F12DA6313597522E0592C69D381E60A2CA66364F429CB182BB32CCE3727974484B4A23F61E99AC494C710203010001}"}]
```

#### 2.2. POST Request

```bash
curl -i -X POST -H 'Accept: application/json' http://localhost:9572/api/config/friends --data '{"alias_name": "Eve", "public_key":"PubKey(go-peer/rsa){3082020A0282020100C971936CE7F037D60E72613552F9E5DDAA7367CF414E1EC97B1C622E9DAC66DEE488048CECB51ACA082EF1EA1F6DF05FEB595B8A075C012634EEFDA62905717D2FC4DCBEDB824F0E92015E6124C81FD9D6B5E0EE13C685F6E226CDA5646DC2BE32D2FDDE486B0B15F4B7455CD1311F604A822C321B304ECBD599D2D7B4A8FB380F38AAEBCC2D1176E1D2BA85F38E25B7879DD61A8C290F55BAB4502221F23DDF6F75B5B3CA631D63B736FD7B7E6F8F9A82F55DE5B673862F0F324F4F911502810477E820946057F951B57E44EC79525BD10B472D05F57A7CAAB835AE55E71129CF9B1CC54175989E1BE86697F9A4C560D09C179CD332E05550F169DEC6318D4E8172F009DC82837A418E454A75E4CAE5A098161099BEA499FFD56E98433ABECAE1961A864388D355EF29C02DC1DEE315C03D16DA6687B6AD67D544A20E541ADB450D1CC57869EB21D3B53368CB716DCBAF18E625A4A68081651C2E5AEA28549F141DBBBB1F500EE970303DE1BC82098B130D202234322D7C1E8A71D71F016E10ACF523303DB48BD1DA9B1D2E623012557CFB81176F2195872F244E6149FE03395951AEE6F90B4808A88796875264A4FF177504D5139EE4729D9603FC3B0F448E0F3E95865CD5234A169DEF8EE07A067DED78E782A534F12DA6313597522E0592C69D381E60A2CA66364F429CB182BB32CCE3727974484B4A23F61E99AC494C710203010001}"}'
```

#### 2.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Mon, 07 Aug 2023 00:22:19 GMT
Content-Length: 23

success: update friends
```

#### 2.3. DELETE Request

```bash
curl -i -X DELETE -H 'Accept: application/json' http://localhost:9572/api/config/friends --data '{"alias_name": "Eve"}'
```

#### 2.3. DELETE Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Mon, 07 Aug 2023 00:23:12 GMT
Content-Length: 22

success: delete friend
```

### 3. /api/network/online

#### 3.1. GET Request

```bash
curl -i -X GET -H 'Accept: application/json' http://localhost:9572/api/network/online
```

#### 3.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 07 Aug 2023 00:23:25 GMT
Content-Length: 18
```

```json
["localhost:9581"]
```

#### 3.2. DELETE Request

```bash
curl -i -X DELETE -H 'Accept: application/json' http://localhost:9572/api/network/online --data 'localhost:9581'
```

#### 3.2. DELETE Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Mon, 07 Aug 2023 00:24:14 GMT
Content-Length: 33

success: delete online connection
```

### 4. /api/network/request

#### Prefix script

```bash
#!/bin/bash

str2hex() {
    local str=${1:-""}
    local fmt="%02X"
    local chr
    local -i i
    for i in `seq 0 $((${#str}-1))`; do
        chr=${str:i:1}
        printf "${fmt}" "'${chr}"
    done
}

JSON_DATA='{
        "method":"POST",
        "host":"hidden-echo-service",
        "path":"/echo",
        "head":{
            "Accept": "application/json"
        },
        "body":"aGVsbG8sIHdvcmxkIQ=="
}';

# POST = request with response from service
# PUT  = broadcast without response from service
PUSH_FORMAT='{
        "receiver":"Alice",
        "hex_data":"'$(str2hex "$JSON_DATA")'"
}';
```

#### 4.1. POST Request

```bash
curl -i -X POST -H 'Accept: application/json' http://localhost:9572/api/network/request --data "${PUSH_FORMAT}"
```

#### 4.1. POST Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 07 Aug 2023 00:31:27 GMT
Content-Length: 113
```

```json
{"code":200,"head":{"Content-Type":"application/json"},"body":"eyJlY2hvIjoiaGVsbG8sIHdvcmxkISIsInJldHVybiI6MX0K"}
```

#### 4.2. PUT Request

```bash
curl -i -X PUT -H 'Accept: application/json' http://localhost:9572/api/network/request --data "${PUSH_FORMAT}"
```

#### 4.2. PUT Response

```
HTTP/1.1 200 OK
Date: Sun, 06 Aug 2023 23:21:20 GMT
Content-Length: 18
Content-Type: text/plain; charset=utf-8

success: broadcast
```

### 5. /api/network/key

#### 5.1. GET Request

```bash
curl -i -X GET -H 'Accept: application/json' http://localhost:9572/api/network/key
```

#### 5.1. GET Response

```
HTTP/1.1 200 OK
Date: Sun, 06 Aug 2023 23:23:01 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8

used_network_key
```

#### 5.2. POST Request

```bash
curl -i -X POST -H 'Accept: application/json' http://localhost:9572/api/network/key --data "used_network_key"'
```

#### 5.2. POST Response

```
HTTP/1.1 200 OK
Date: Sun, 06 Aug 2023 23:22:49 GMT
Content-Length: 24
Content-Type: text/plain; charset=utf-8

success: set network key
```

### 6. /api/node/key

#### 6.1. GET Request

```bash
curl -i -X GET -H 'Accept: application/json' http://localhost:9572/api/node/key
```

#### 6.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 07 Aug 2023 00:35:45 GMT
Transfer-Encoding: chunked
```

```json
["PubKey(go-peer/rsa){3082020A0282020100C17B6FA53983050B0339A0AB60D20A8A5FF5F8210564464C45CD2FAC2F266E8DDBA3B36C6F356AE57D1A71EED7B612C4CBC808557E4FCBAF6EDCFCECE37494144F09D65C7533109CE2F9B9B31D754453CA636A4463594F2C38303AE1B7BFFE738AC57805C782193B4854FF3F3FACA2C6BF9F75428DF6C583FBC29614C0B3329DF50F7B6399E1CC1F12BED77F29F885D7137ADFADE74A43451BB97A32F2301BE8EA866AFF34D6C7ED7FF1FAEA11FFB5B1034602B67E7918E42CA3D20E3E68AA700BE1B55A78C73A1D60D0A3DED3A6E5778C0BA68BAB9C345462131B9DC554D1A189066D649D7E167621815AB5B93905582BF19C28BCA6018E0CD205702968885E92A3B1E3DB37A25AC26FA4D2A47FF024ECD401F79FA353FEF2E4C2183C44D1D44B44938D32D8DBEDDAF5C87D042E4E9DAD671BE9C10DD8B3FE0A7C29AFE20843FE268C6A8F14949A04FF25A3EEE1EBE0027A99CE1C4DC561697297EA9FD9E23CF2E190B58CA385B66A235290A23CBB3856108EFFDD775601B3DE92C06C9EA2695C2D25D7897FD9D43C1AE10016E51C46C67F19AC84CD25F47DE2962A48030BCD8A0F14FFE4135A2893F62AC3E15CC61EC2E4ACADE0736C9A8DBC17D439248C42C5C0C6E08612414170FBE5AA6B52AE64E4CCDAE6FD3066BED5C200E07DBB0167D74A9FAD263AF253DFA870F44407F8EF3D9F12B8D910C4D803AD82ABA136F93F0203010001}","PubKey(go-peer/rsa){3082020A0282020100A4FDC84A2BE3212F32C75F54BB259CFB8F8B701D7B3E3B1FD6D0AA4CD4B9E2DDE44005362C1483895065902F0C68C7B4C5EDA89BE8A5C08019C1D3337E3D3660620AF3793224816866140A4E9E47747F53EDC32A31ADFAF493DF222A5BF484EF91AFB6AC8AE426F8E3F981180689F666000EBC9285DC74C12CA6B287CB93AD1F87097E9104851565DE127CBE675A3E6486E702EFE027A53A98A19C5E0D12458E620DD9443FF99D2ACFD0F1C813262E0A85C7E442F9ABE4BD316C794B3F95676E59B6E5220AB9744419BD3CA1B082C5540CF537B15737E5B3019EE50AB9501CEBF2D38CA4D3F90719D575BB8B6ACC0D8EDEB2BDF69A5293945F7F0C0B5DB9AAC15EEB13CE107B8A7C72E670EC82A53783DEFDD7A836698CC5FE0F4C95B72C6A64ED114E4E4F69B22F46D5255AA26E357999BDBBF003917C295A50682F4E39F5E89BF8A40360EA77DA0FEDB29C902B11D403B2DBE851F8CCAD0B888EF7161D4A4DEB70CA43515FE5CA3E19AA0587A9CE162BAA0401E81425E76CDCEB851075F32A97C958E0C8617684563C39DF8B89B6C348D734F3F5996F026B149AFB200C465B28D86D0D06652D8D2B66F5E5E1820AC42AE6183A3C0AFA7DAE39DDF880E1EB1DB76F937788CC24A60EE001CD2A5EBA9C919F6FEF54A4531F1B53DD8507FEE30605CA88E98499ED2809BA3A1CE1EF4EB60953B638B3395722CFAA6564896E7DCF0203010001}"]
```

#### 6.2. POST Request

```bash
curl -i -X POST -H 'Accept: application/json' http://localhost:9572/api/node/key --data '{"priv_key":"PrivKey(go-peer/rsa){3082092802010002820201009D5159F646E8E6F55B0FCD2445BDD5320F75AB6A5A5CEDFBB60B5B52FFC54E08543A8F5A329B1CCE1BE815C6AD7A30BEC3BE8EDCE0DB212C503F75D208AE2D1472C891DF7D5ECBC088358B33462579C95E43BD2E7C1577E9FB95A5B1105A55CD2459D499660659A843CC801F9260FA5846C67211BFE456DDD1121A8E0CBC7849A5301FEB18496B8DA3B6A7669229744D20C360CFEBB1F6245326F1D1EEE13A902A7A8E14A63457053DB6A285ABA6654960A310B7B54F43F191CCFFFEB35375C23FEEAF3D039FF57650CCF52FD76FA7EC9FD18627D8DE95486FA5EB8B50AF79A20CDF2DCBC4BC4B6931C57842D4E20AF531CCCB89D70F83232938AE0290963A379152C912989D8C85898123428C2A14D5046475D6EE9B6658B7B99DB4BCD863268792CB1185C0DDE874978F7DFA8DCAA1E97D092CA852C061F2FEBA6ECDC7F9CF08791BDADC8A2488405EDE35D59D407AD39C5FC2A167F2AC7F631CE64F7A099FFB41BBEB7C6036058E97CD9DF43C287B9A3A1D4E8D8273CBE65DE103F05CA87C0B3940131B3388AF97B383DB7D0C6A1BCFC93971C5D8D4CE6C3A61E1DA431C43B32D393422FDD6716BE27D9A5928402634940ADEE8A15C18B1F036636312843FD1738B98A4F6A3677455FB61D8310FEE2C1A0B705E5DBCECB4806F1603018F8E9CA18B0EA715F2D9A5498CF78BAE6B2CB4ADA658B4F8E074D78E5045B174D62302030100010282020021B59BDC4CC77D2DD7EC63DDC0DFF37DFD980E3A04D0E2E1CBD955214CD31F6C637804DDA3F85ECCBF6814BA74D3B8FC377F6EA75FBB34B9851C8407947A96084AAC35ADB8F4861E64516CD978CF70F03835B5A4EF4BBE5D31DE98197FD28B8E209AEB164FA94EAEE2904068037AAA4A1E2849AB09FE48AAD130DAE5D34ED34B9C8CDA5A0AE3389BAA17EA78ED1ADAE3E800558F5806D322677AF1D83522A7E4DA65566A904EA8D2E3AD6DD7CCB723FEFC2914DCF889DA9A39CEBE8FFA270915AD935C936B626C3B8506D6070157D898B88A31FFF9D580117C73062CDD062CBF0F9906FC21D4E327D0556AF68F1D3C91DBB0F17040D7FF169AAF9D81C92F979B9888618AEA88917281D14078891B944942ACB5073A67DFF2F889E470BD5FF29EEE2F04131629324724A7CD065200C34C4B8E44C64F2A3D25E9712A7CC83FA7F66929CFBFEE12137C0BF32FA109F630D3872910B20A3DD7DC12CFD0EB7E5D13883CA01E9EA032A6993304E9BFC2DA30A7D02FB3522661F1B6842EF097FC6E3D64381775DB516063F8762E036DAB436016757F393F0D7E1D6147EAF5BB8D7D3BFD764CC7F64EEB2C1163E4EC7E6D72AD50CC273E0D9B53D39788020C2101F7F1AE0223FF74E6ADEB659F21E2C2143AA1352E39C16447873F3EB060BC016C5F1A73B89DB0115AEDA6E2FE3588611064D6C9A288E2904BBA6D4985F007E6AC548B010282010100C239EF2FC77DE8A6022A34C9A0DAD65427841CC3BC067E12AE2A25FC282541B2FBE857773252058BEE309EA18DDB9EAA6038AF8DA0E36E901A57B369A257A814DC00D8A2E24E95CD491BAC1C258306EFEA51FA3CCEDDBC97A72B06EA3D9A42E07AF1EF83282092BEAFD117648DEA9F3D887B05136213C471BC6A01ED4C5586F661B61C8D5F130D80BF92C4C641CA49EB953DFD7E7B035816CF0EBB797E8A9018952710D4B67415C3269CD437225A640E213643AF24446595BE0E29AB5C9D1C7BBB8BCFBCB6113592DF68C693D7961A6223D9DDBA6DC422BE7AAE1C86B38F37F1DD10B5132101AD4B18ACC26B8F413D84D0A5A091BA4D3A453C06562CDADA60810282010100CF5A4B60FBC10BE3FD78FABA01C8E24839488DF13EF530B3FD7A02FC0D6C32EAEAD81CC12B4DE0B9016F7F164DAF8A8AABC3FD139BF5EC13AE69AADEB109514530F7ED49D88B441879F72178821620800744877B98C0F8B416546B754AD4E526470BB010A114F28F85CB7E027B154266D57839FA21A42B52FBB6B6FC01B4CF0AE04C6555F354F0A3C901BEAB60898B1BD652995F663AEB6FF7389023CAD11374EC94A2EDF6DFECBFFD964D8A30B6A8597505B0339FF352BAAB84C8AF3C9A0DB19A996875C0CBE7F102ADB227255C9B058B2D80046B961B3632FA4B152EBBA102240A6113139D58285A089BFA2BADFFFB52832047ACFB00BBE193F046CF3764A3028201004B77E9699E515D29CC238C39604848099105338C16AE4B24850A19925E2303E14122A981C64ABA9F01A160B21385E7A3FA196C955293ACAD4F9F0E36987F08EF7A00D62C8C54CEBE628EDF145CBB09E205216F635B5A2B629DF10911D177F44E775734A2B8DFD74542D9B3063E6291177EC596564EC0B18F240FE8C9C4E462B9AF83EC9A3DFC103E1BB232C57A60D8D2323E5116694406616E7921FD765EEED8AE73EC854A93D6B4EA76FBBAA49D8CCD34B87A1A3EB458E8935DBB713B5E4CE8031AB350774A3E8FE0413D0FCB3026F64549ED6EF821C3025276FEDC943EAD154CB9A632559BEA3308D670010D3BE3648D121E4F219DCA6B86844BCDC8081C810282010100CD24F164CF4ECBFBD1C00A9752C2B3956F0F2857A0C926593D13A4B648755EDEEA5FCBFB1563E44C456E5116F8DF0EBB697AEAFCA695A4EE47E58546F3725B7490210A23C058F09322BFECDE741D7E240C8CB15A07E40B6AE898B7040178260A3BCA05743E5A222CFADB3C1D2A36FB4E102EF5755229412FC5979CEC30A7F91B329482C1898FE4D0B642C2A87D473758E02F324C9F30F5D3FD8C7996DFC7006FF2CC8F71CD88F78B6F516FFFA3786390B5E55DD185934FAB1D9CAE8C28F1E5506CBB100D4824B4A1CEDB98618066617D1798798A6602C98352E62CB89556CED1F6644A6C7C407482DFA89AE0E4AC2E9130AE4896813E881859F26A8B33E202EF028201001AD162216179FAA9CA9012CCBE241AABCD2217011DB1FD049FC434491EA0A8C73FCFBDA0B92A2E34C7BC21D532C3A5FD36020BCAECD6C3C14E2BD642FA4FD0DF079047C2D061F148AED765A36EBAC1FC02E978AB8D397582457F3D11B5CE7A269868326F49E7087DE9A1E75DCD5BA249CBA60A8405EF53E86FBA1936A76B3542DD78031406BDEC421DEB969F5478AA9C53E3B34FD23D62986399C6978777B57DA715C72AB193210C48B26133DEA5BDDEE087E2E9809E55030FBDBB71F87728D5D923B15BEC9E1A38E85D33BBA27AC2BCBCC9AD7F759B76105120E758A58EAA359E793D05E23493A9A2925032E1234782D7728055EEC0FF28B30D3F613EE3A356}"}'
```

#### 6.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Mon, 07 Aug 2023 00:41:09 GMT
Content-Length: 27

success: update private key
```

#### 6.3. DELETE Request

```bash
curl -i -X DELETE -H 'Accept: application/json' http://localhost:9572/api/node/key
```

#### 6.3. DELETE Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 08 Oct 2023 07:10:24 GMT
Content-Length: 26

success: reset private key
```
