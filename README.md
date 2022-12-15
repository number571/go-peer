# go-peer

> Library for create secure and anonymity decentralized networks

<img src="examples/images/go-peer_logo.png" alt="go-peer_logo.png"/>

The `go-peer` library contains a large number of functions necessary to ensure the security of transmitted or stored information, as well as for the anonymity of nodes in the decentralized form. The library can be divided into several main modules:

1. The `crypto` module represents cryptographic primitives: 1) asymmetric encryption, decryption; 2) asymmetric signing and signature verification; 3) symmetric encryption and decryption; 4) hashing; 5) entropy enhancement; 6) computational problems (puzzles); 7) cryptographically stable pseudorandom number generator.
2. The `client` module for encrypting and decrypting information with the attached data integrity (hash), authentication (signature) and confirmation (work). It is a basic part of the `anonymity` module.
3. The `client/queue` module represents the generation, storage and issuance of encrypted messages every time the period specified by the application is reached. Uses the `client` module.
4. The `network` module is a decentralized communication between network nodes. It does not represent any protection of information and anonymity of participants.
5. The `network/anonymity` module to ensure anonymity based on the fifth stage. Presents the main functions for working with the network on top of the `network` and `queue` modules.
6. The `storage` module includes two types of storage: `in-memory` and"`crypto`. The second type of storage can be used for secure storage of passwords and private keys.
7. The `storage/database` module is a `key-value` database with the functions of value encryption and key hashing.

> Examples of works in the directory [https://github.com/number571/go-peer/examples/modules](https://github.com/number571/go-peer/tree/master/examples/modules "Modules");

## Library based applications

1. [Hidden Lake Service](#1-hidden-lake-service) 
2. [Hidden Lake Messenger](#2-hidden-lake-messenger) 
3. [Another applications](#3-another-applications) 

## 1. Hidden Lake Service

> [github.com/number571/go-peer/tree/master/cmd/hls](https://github.com/number571/go-peer/tree/master/cmd/hls "HLS")

<img src="examples/images/hls_logo.png" alt="hls_logo.png"/>

The `Hidden Lake Service` is the core of an anonymous network with theoretically provable anonymity. HLS is based on the `fifth^ stage` of anonymity and is an implementation of an `abstract` anonymous network based on `queues`. It is a `peer-to-peer` network communication with trusted `friend-to-friend` participants. All transmitted and received messages are in the form of `end-to-end` encryption.

A feature of HLS (compared to many other anonymous networks) is its easy adaptation to a hostile centralized environment. Anonymity can be restored literally from one node in the network, even if it is the only point of failure.

> More information about HLS in the [habr.com/ru/post/696504](https://habr.com/ru/post/696504/ "Habr HLS")

### How it works

Each network participant sets a message generation period for himself (the period can be a network constant for all system participants). When one cycle of the period ends and the next begins, each participant sends his encrypted message to all his connections (those in turn to all of their own, etc.). If there is no true message to send, then a pseudo message is generated (filled with random bytes) that looks like a normal encrypted one. The period property ensures the anonymity of the sender.

<p align="center"><img src="examples/images/hls_queue.jpg" alt="hls_queue.jpg"/></p>
<p align="center">Figure 1. Queue and message generation in HLS.</p>

Since the encrypted message does not disclose the recipient in any way, each network participant tries to decrypt the message with his private key. The true recipient is only the one who can decrypt the message. At the same time, the true recipient acts according to the protocol and further distributes the received packet, even knowing the meaninglessness of the subsequent dispatch. This property makes it impossible to determine the recipient.

> Simple example of the `client` module (encrypt/decrypt functions) in the directory [github.com/number571/go-peer/examples/modules/client](https://github.com/number571/go-peer/tree/master/examples/modules/client "Module client");

<p align="center"><img src="examples/images/hls_view.jpg" alt="hls_view.jpg"/></p>
<p align="center">Figure 2. Two participants are constantly generating messages for their periods on the network. It is impossible to determine their real activity.</p>

Data exchange between network participants is carried out using application services. HLS has a dual role: 1) packages traffic from pure to anonymizing and vice versa; 2) converts external traffic to internal and vice versa. The second property is the redirection of traffic from the network to the local service and back.

<p align="center"><img src="examples/images/hls_service.jpg" alt="hls_service.jpg"/></p>
<p align="center">Figure 3. Interaction of third-party services with the traffic anonymization service.</p>

As shown in the figure above, HLS acts as an anonymizer and handlers of incoming and outgoing traffic. The remaining parts in the form of applications and services depend on third-party components (as an example, `HLM`).

> More details in the work [Theory of the structure of hidden systems](https://github.com/number571/go-peer/blob/master/hidden_systems.pdf "TSHS")

### Example

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
```
PUSH_FORMAT="{
        \"receiver\":\"Pub(go-peer/rsa){3082020A0282020100B752D35E81F4AEEC1A9C42EDED16E8924DD4D359663611DE2DCCE1A9611704A697B26254DD2AFA974A61A2CF94FAD016450FEF22F218CA970BFE41E6340CE3ABCBEE123E35A9DCDA6D23738DAC46AF8AC57902DDE7F41A03EB00A4818137E1BF4DFAE1EEDF8BB9E4363C15FD1C2278D86F2535BC3F395BE9A6CD690A5C852E6C35D6184BE7B9062AEE2AFC1A5AC81E7D21B7252A56C62BB5AC0BBAD36C7A4907C868704985E1754BAA3E8315E775A51B7BDC7ACB0D0675D29513D78CB05AB6119D3CA0A810A41F78150E3C5D9ACAFBE1533FC3533DECEC14387BF7478F6E229EB4CC312DC22436F4DB0D4CC308FB6EEA612F2F9E00239DE7902DE15889EE71370147C9696A5E7B022947ABB8AFBBC64F7840BED4CE69592CAF4085A1074475E365ED015048C89AE717BC259C42510F15F31DA3F9302EAD8F263B43D14886B2335A245C00871C041CBB683F1F047573F789673F9B11B6E6714C2A3360244757BB220C7952C6D3D9D65AA47511A63E2A59706B7A70846C930DCFB3D8CAFB3BD6F687CACF5A708692C26B363C80C460F54E59912D41D9BB359698051ABC049A0D0CFD7F23DC97DA940B1EDEAC6B84B194C8F8A56A46CE69EE7A0AEAA11C99508A368E64D27756AD0BA7146A6ADA3D5FA237B3B4EDDC84B71C27DE3A9F26A42197791C7DC09E2D7C4A7D8FCDC8F9A5D4983BB278FCE9513B1486D18F8560C3F31CC70203010001}\",
        \"hex_data\":\"$(str2hex "$JSON_DATA")\"
}";
```

Build and run nodes
```bash
$ cd examples/cmd/echo_service
$ make
```

Logs from `middle_hls` node. When sending requests and receiving responses, `middle_hls` does not see the action. For him, all actions and moments of inaction are equivalent.

<p align="center"><img src="examples/images/hls_logger.png" alt="hls_logger.png"/></p>
<p align="center">Figure 4. Output of all actions and all received traffic from the middle_hls node.</p>

Send request
```bash
$ ./request.sh
```

Get response
```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 15 Dec 2022 07:42:49 GMT
Content-Length: 97

{"result":"7b226563686f223a2268656c6c6f2c20776f726c6421222c2272657475726e223a317d0a","return":1}
```

Decode response
```json
{"echo":"hello, world!","return":1}
```

> Simple examples of the `anonymity` module in the directory [github.com/number571/go-peer/examples/modules/network/anonymity](https://github.com/number571/go-peer/tree/master/examples/modules/network/anonymity "Module anonymity");

## 2. Hidden Lake Messenger

> [github.com/number571/go-peer/tree/master/cmd/hlm](https://github.com/number571/go-peer/tree/master/cmd/hlm "HLM");

<img src="examples/images/hlm_logo.png" alt="hlm_logo.png"/>

The `Hidden Lake Messenger` is a messenger based on the core of an anonymous network with theoretically provable anonymity of HLS. A feature of this messenger is the provision of anonymity of the fact of transactions (sending, receiving).

HLM is an application that implements a graphical user interface (GUI) on a browser-based HTML/CSS/JS display. Most of the code is based on the bootstrap library https://getbootstrap.com /. GUI is adapted to the size of the window, so it can be used both in a desktop and in a smartphone.

> More information about HLM in the [habr.com/ru/post/701488](https://habr.com/ru/post/701488/ "Habr HLM")

### How it works

Most of the code is a call to API functions from the HLS kernel. However, there are additional features aimed at the security of the HLM application itself.

Firstly, there is registration and authorization, which does not exist in the HLS core. Registration performs the role of creating / depositing a private key `PrivKey` in order to save it through encryption. 

The encryption of the private key is carried out on the basis of the entered `login (L) / password (P)`, where the login acts as a cryptographic salt. The concatenation of the login and password `L||P` is hashed `2^20` times `K = H(L||H(...L||(H(L||P)...))` to increase the password security by about `20 bits` of entropy and turn it into an encryption key `K`. The resulting `K` is additionally hashed by `H(K)` and stored together with the encrypted version of the private key `Q = E(K, PrivKey)`.

<p align="center"><img src="examples/images/hlm_auth.jpg" alt="hlm_auth.jpg"/></p>
<p align="center">Figure 5. Data encryption with different types of input parameters.</p>

Authorization is performed by entering a `login/password`, their subsequent conversion to `K' and H(K')`, subsequent comparison with the stored hash `H(K) = H(K')?` and subsequent decryption of the private key `D(K, Q) = D(K, E(K, PrivKey)) = PrivKey`.

Secondly, the received key K is also used to encrypt all incoming and outgoing messages `C = E(K, M)`. All personal encrypted messages `C` are stored in the local database of each individual network participant.

### Example

The example will involve (as well as in HLS) three nodes `middle_hls, node1_hlm and node2_hlm`. The first one is only needed for communication between `node1_hlm` and `node2_hlm` nodes. Each of the remaining ones is a combination of HLS and HLM, where HLM plays the role of an application and services, as it was depicted in `Figure 3`.

Build and run nodes
```bash
$ cd examples/cmd/anon_messenger
$ make
```

The output of the `middle_hls` node is similar to `Figure 4`.

Than open browser on `localhost:8080`
<p align="center"><img src="examples/images/hlm_about.png" alt="hlm_about.png"/></p>
<p align="center">Figure 6. Home page of the HLM application.</p>

Next, you need to register by going to the Sign up page. Enter your `login/password` and insert the private key `priv.key`. That key is located in `examples/cmd/anon_messenger/node2_hlm`.

After the registration procedure, re-enter your `login/password`. After that, you will have the functions of adding connections and friends, as well as the communication with friends itself. In the example, friend `Alice` will be added by default. 

To see the success of sending and receiving messages, you need to do all the same operations, but with `localhost:7070` and `node2_hlm`. This node will be Alice.

> More example images about HLM pages in the [github.com/number571/go-peer/cmd/hlm/examples/images](https://github.com/number571/go-peer/tree/master/cmd/hlm/examples/images "Path to HLM images")

## 3. Another applications

* Hidden Message Service: [github.com/number571/go-peer/tree/master/cmd/hms](https://github.com/number571/go-peer/tree/master/cmd/hms "HMS");
* Hidden Message Client: [github.com/number571/go-peer/tree/master/cmd/hmc](https://github.com/number571/go-peer/tree/master/cmd/hmc "HMC");
* Asymmetric Keys Generator: [github.com/number571/go-peer/tree/master/cmd/akg](https://github.com/number571/go-peer/tree/master/cmd/akg "AKG");
* Cryptographic Data Storage: [github.com/number571/go-peer/tree/master/cmd/cds](https://github.com/number571/go-peer/tree/master/cmd/cds "CDS");
* In developing... Union Block Chain: [github.com/number571/go-peer/tree/master/cmd/ubc](https://github.com/number571/go-peer/tree/master/cmd/ubc "UBC");

### Deprecated applications

* Hidden Lake (new release = cmd/hls+cmd/hlm+...): [github.com/number571/hidden-lake](https://github.com/number571/hidden-lake "HL");
* Hidden Email Service (new release = cmd/hms+cmd/hmc): [github.com/number571/hes](https://github.com/number571/hes "HES");
* Hidden Lake Service (new release = cmd/hls): [github.com/number571/hls](https://github.com/number571/hls "HLS");
