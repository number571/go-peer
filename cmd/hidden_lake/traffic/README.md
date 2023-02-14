# HLT

> Hidden Lake Traffic

<img src="../../../examples/images/hlt_logo.png" alt="hlt_logo.png"/>

The `Hidden Lake Traffic` is an application that saves traffic passed through HLS. The saved traffic can be used by other applications when they were offline. HLT provides an API for loading and unloading messages. Messages are stored in the database based on the "ring" structure. Thus, new messages will overwrite the old ones after some time.

## How it works

HLT emulates HLS to receive messages. In this scenario, HLT has only the functions of accepting messages, without the ability to generate or send them via HLS or independently.

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
