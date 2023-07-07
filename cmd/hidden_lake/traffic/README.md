# HLT

> Hidden Lake Traffic

<img src="../../../images/hlt_logo.png" alt="hlt_logo.png"/>

The `Hidden Lake Traffic` is an application that saves traffic passed through HLS. The saved traffic can be used by other applications when they were offline. HLT provides an API for loading and unloading messages. Messages are stored in the database based on the "ring" structure. Thus, new messages will overwrite the old ones after some time.

> More information about HLT in the [habr.com/ru/post/717184](https://habr.com/ru/post/717184/ "Habr HLT")

## How it works

HLT emulates HLS to receive messages. In this scenario, HLT has only the functions of accepting messages, without the ability to generate or send them via HLS or independently.

<p align="center"><img src="../../../examples/images/hlt_client.gif" alt="hlt_client.gif"/></p>
<p align="center">Figure 1. Example of running HLT client.</p>

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ cd ./cmd/hidden_lake/traffic
$ make build # create hlt, hlt_linux, hlt_windows.exe and copy to ./bin
$ make run # run ./bin/hlt

> [INFO] 2023/06/03 15:39:13 HLT is running...
> ...
```

Open ports `9581` (HTTP, interface).
Creates `./hlt.cfg` or `./_mounted/hlt.cfg` (docker), `./hlt.db` or `./_mounted/hlt.db` (docker) files.
The file `hlm.db` stores all sent/received messages as structure `ring` from network HL. 

Default config `hlt.cfg`

```json
{
	"logging": [
		"info",
		"warn",
		"erro"
	],
	"address": ":9581",
	"connection": "service:9571"
}
```

If traffic works not in docker's enviroment than need rewrite connection host in `hlt.cfg` file from `service` to IP address (example: `127.0.0.1:9571` for local network).

Build and run with docker

```bash 
$ cd ./cmd/hidden_lake/traffic
$ make docker-build 
$ make docker-run

> [INFO] 2023/06/03 08:44:14 HLT is running...
> ...
```

## Example 

Build and run service
```bash
$ cd examples/_cmd/traffic_keeper
$ make
```

Run client
```bash
$ cd client
$ go run ./main.go w 'hello, world!'
$ go run ./main.go h
$ go run ./main.go r cb3c6558fe0cb64d0d2bad42dffc0f0d9b0f144bc24bb8f2ba06313af9297be4 # hash get by 'h' option
```

## Config structure

```
"network"      A network key created to encapsulate connections
"address"      API address for HLT functions
"connection"   Connection to HLS as HLS
```

```json
{
    "network": "network-key",
	"address": "localhost:9573",
	"connection": "localhost:9571"
}
```

## Response structure from HLT API

```
"result" is string
"return" is int; 1 = success
```

```json
{
	"result":"go-peer/hidden-lake-traffic",
	"return":1
}
```

## HLT API

```
1. GET      /api/hashes
2. GET/POST /api/message
```

### 1. /api/hashes

#### 1.1. GET Request

```bash
curl -i -X GET -H 'Accept: application/json' http://localhost:9573/api/hashes
```

#### 1.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 14 Feb 2023 12:51:23 GMT
Content-Length: 154

{"result":"6d815acf176c71ab9e55b65a38dbfc266bcda6a12ac1a5e660720a077ea4bd23,31f3f211c1ccbe2ab367e743e109f7f9702521447e5f6348ef0d8ab7a1ccd756","return":1}
```

### 2. /api/message

#### 2.1. GET Request

```bash
curl -i -X GET -H 'Accept: application/json' 'http://localhost:9573/api/message?hash=31f3f211c1ccbe2ab367e743e109f7f9702521447e5f6348ef0d8ab7a1ccd756'
```

#### 2.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 14 Feb 2023 12:50:19 GMT
Transfer-Encoding: chunked

{"result":"7b2268656164223a7...3030303966383334227d7d","return":1}
```

#### 2.2. POST Request

```bash
curl -i -X POST -H 'Accept: application/json' http://localhost:9573/api/message -d @README_message.txt
```

#### 2.2. POST Response

```
HTTP/1.1 100 Continue

HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 14 Feb 2023 12:45:56 GMT
Content-Length: 32

{"result":"success","return":1}
```
