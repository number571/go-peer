# HLR

> Hidden Lake Remoter

<img src="_images/hlr_logo.png" alt="hlr_logo.png"/>

The `Hidden Lake Remoter` this is a service that provides the ability to make remote calls on the anonymous network core (HLS) with theoretically provable anonymity.

> [!CAUTION]
> This application can be extremely dangerous. Use HLR with caution.

> More information about HLR in the [habr.com/ru/articles/830130](https://habr.com/ru/articles/830130/ "Habr HLR")

## Installation

```bash
$ go install github.com/number571/go-peer/cmd/hidden_lake/applications/remoter/cmd/hlr@latest
```

## How it works

Most of the code is a call to API functions from the HLS kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service.

The server providing the remote access service is waiting for a request in the form of a command. The command does not depend on the operating system and therefore should have a small additional syntax separating the launch of the main command and its arguments.

As an example, to create a file with the contents of "hello, world!" and then reading from the same file, you will need to run the following command:

```bash
bash[@remoter-separator]-c[@remoter-separator]echo 'hello, world' > file.txt && cat file.txt
```

The `[@remoter-separator]` label means that the arguments are separated for the main command.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ cd ./cmd/hidden_lake/applications/remoter
$ make build # create hlr, hlr_[arch=amd64,arm64]_[os=linux,windows,darwin] and copy to ./bin
$ make run # run ./bin/hlr

> [INFO] 2023/06/03 15:30:31 HLR is running...
> ...
```

Open port `9532` (HTTP, incoming).
Creates `./hlr.yml` file.

Default config `hlr.yml`

```yaml
settings:
  exec_timeout_ms: 5000
  password: 4otg9sohTw8Lv8PheDZ7fOD5j5v5sU
logging:
- info
- warn
- erro
address:
  incoming: 127.0.0.1:9532
```

## Running options

```bash
$ ./hlr -path=/root
# path = path to config
```

## Example

The example will involve three nodes `recv_hlc, send_hls` and three repeaters `middle_hlt_1, middle_hlt_2, middle_hlt3_`. The three remaining nodes are used only for the successful connection of the two main nodes. In other words, HLT nodes are traffic relay nodes.

Build and run nodes
```bash
$ cd examples/anonymity/remoter/routing
$ make
```

Than run command
```bash
$ cd examples/anonymity/remoter
$ make request # go run ./_request/main.go
```

Got response
```json
{"code":200,"head":{"Content-Type":"application/octet-stream"},"body":"aGVsbG8sIHdvcmxkCg=="}
```
