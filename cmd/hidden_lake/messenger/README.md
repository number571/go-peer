# HLM

> Hidden Lake Messenger

<img src="_images/hlm_logo.png" alt="hlm_logo.png"/>

The `Hidden Lake Messenger` is a messenger based on the core of an anonymous network with theoretically provable anonymity of HLS. A feature of this messenger is the provision of anonymity of the fact of transactions (sending, receiving).

HLM is an application that implements a graphical user interface (GUI) on a browser-based HTML/CSS/JS display. Most of the code is based on the bootstrap library https://getbootstrap.com /. GUI is adapted to the size of the window, so it can be used both in a desktop and in a smartphone.

> More information about HLM in the [habr.com/ru/post/701488](https://habr.com/ru/post/701488/ "Habr HLM")

## How it works

Most of the code is a call to API functions from the HLS kernel. However, there are additional features aimed at the security of the HLM application itself.

<p align="center"><img src="_images/hlm_chat.gif" alt="hlm_chat.gif"/></p>
<p align="center">Figure 1. Example of chat room in HLM.</p>

Firstly, there is registration and authorization, which does not exist in the HLS core. Registration performs the role of creating / depositing a private key `PrivKey` in order to save it through encryption. 

The encryption of the private key is carried out on the basis of the entered `login (L) / password (P)`, where the login acts as a cryptographic salt. The concatenation of the login and password `L||P` is hashed `2^20` times `K = H(L||H(...L||(H(L||P)...))` to increase the password security by about `20 bits` of entropy and turn it into an encryption key `K`. The resulting `K` is additionally hashed by `H(K)` and stored together with the encrypted version of the private key `Q = E(K, PrivKey)`.

<p align="center"><img src="_images/hlm_auth.jpg" alt="hlm_auth.jpg"/></p>
<p align="center">Figure 2. Data encryption with different types of input parameters.</p>

Authorization is performed by entering a `login/password`, their subsequent conversion to `K' and H(K')`, subsequent comparison with the stored hash `H(K) = H(K')?` and subsequent decryption of the private key `D(K, Q) = D(K, E(K, PrivKey)) = PrivKey`.

Secondly, the received key K is also used to encrypt all incoming and outgoing messages `C = E(K, M)`. All personal encrypted messages `C` are stored in the local database of each individual network participant.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ cd ./cmd/hidden_lake/messenger
$ make build # create hlm, hlm_[arch=amd64,arm64]_[os=linux,windows,darwin] and copy to ./bin
$ make run # run ./bin/hlm

> [INFO] 2023/06/03 15:30:31 HLM is running...
> ...
```

Open ports `9591` (HTTP, interface) and `9592` (HTTP, incoming).
Creates `./hlm.cfg` or `./_mounted/hlm.cfg` (docker), `./hlm.db` or `./_mounted/hlm.db` (docker) files and `./hlm.stg` or `./_mounted/hlm.stg` (docker).
The file `hlm.db` stores all sent/received messages in encrypted view. The file `hlm.stg` stores all auth information (logins, passwords, private keys) in encrypted view.

Default config `hlm.cfg`

```json
{
	"settings": {
		"message_size_bytes": 1048576,
		"work_size_bits": 20,
		"key_size_bits": 4096
	},
	"logging": [
		"info",
		"warn",
		"erro"
	],
	"address": {
		"interface": ":9591",
		"incoming": ":9592"
	},
	"connection": {
		"service": "127.0.0.1:9572",
		"traffic": "127.0.0.1:9582"
	}
}
```

If messenger works not in docker's enviroment than need rewrite connection hosts in `hlm.cfg` file from `service` and `traffic` to IP addresses (example: `127.0.0.1:9572` and also `127.0.0.1:9581` for local network).

Build and run with docker

```bash 
$ cd ./cmd/hidden_lake/messenger
$ make docker-build 
$ make docker-run

> [INFO] 2023/06/03 08:35:50 HLM is running...
> ...
```

## Example

The example will involve (as well as in HLS) three nodes `middle_hls, node1_hlm and node2_hlm`. The first one is only needed for communication between `node1_hlm` and `node2_hlm` nodes. Each of the remaining ones is a combination of HLS and HLM, where HLM plays the role of an application and services, as it was depicted in `Figure 3` (HLS).

Build and run nodes
```bash
$ cd examples/anon_messenger
$ make
```

The output of the `middle_hls` node is similar to `Figure 4` (HLS).
Than open browser on `localhost:8080`. It is a `node1_hlm`. This node is a Bob.

<p align="center"><img src="_images/hlm_about.png" alt="hlm_about.png"/></p>
<p align="center">Figure 3. Home page of the HLM application.</p>

Next, you need to login by going to the Sign in page. Enter your `login=user` and `password=password` than the private key, friends and connections will be automatically loaded from the storage.

To see the success of sending and receiving messages, you need to do all the same operations, but with `localhost:7070` as `node2_hlm`. This node will be Alice.

Also you can run example with docker-compose. In this example, all nodes have logging enabled
```bash
$ cd examples/anon_messenger/_docker/default
$ make
```

<p align="center"><img src="_images/hlm_logger.png" alt="hlm_logger.png"/></p>
<p align="center">Figure 4. Log of the three nodes with request/response actions.</p>

> More example images about HLM pages in the [github.com/number571/go-peer/cmd/hidden_lake/messenger/_images](https://github.com/number571/go-peer/tree/master/cmd/hlm/_images "Path to HLM images")

## Pages

### Signup page

Sign up login/password and additional private key. If field with private key is null than private key generated.

<img src="images/v2/signup.png" alt="signup.png"/>

### Signin page

Sign in with login/password. Authorized client identified by a private key and can push messages into database.

<img src="images/v2/signin.png" alt="signin.png"/>

### About page

Base information about projects HLM and HLS with links to source.

<img src="images/v2/about.png" alt="about.png"/>

### Settings page

Information about public key and connections. Connections can be appended and deleted.

<img src="images/v2/settings.png" alt="settings.png"/>

### Friends page

Information about friends. Friends can be appended and deleted.

<img src="images/v2/friends.png" alt="friends.png"/>

### Chat page

Chat with friend. The chat is based on web sockets, so it can update messages in real time. Messages can be sent.

<img src="images/v2/chat.png" alt="chat.png"/>
