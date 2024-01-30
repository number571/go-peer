# HLF

> Hidden Lake Filesharer

<img src="_images/hlf_logo.png" alt="hlf_logo.png"/>

The `Hidden Lake Filesharer` is a file sharing service based on the Anonymous Network Core (HLS) with theoretically provable anonymity. A feature of this file sharing service is the anonymity of the fact of transactions (file downloads), taking into account the existence of a global observer.

HLF is an application that implements a graphical user interface (GUI) on a browser-based HTML/CSS/JS display. Most of the code is based on the bootstrap library https://getbootstrap.com/. GUI is adapted to the size of the window, so it can be used both in a desktop and in a smartphone.

> More information about HLF in the [habr.com/ru/articles/789968](https://habr.com/ru/articles/789968/ "Habr HLF")

## Installation

```bash
$ go install github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/cmd/hlf@latest
```

## How it works

Most of the code is a call to API functions from the HLS kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service.

Unlike applications such as HLS, HLT, and HLM, the HLF application does not have a database. Instead, the storage is used, represented by the usual `hlf.stg` directory.

<p align="center"><img src="_images/hlf_download.gif" alt="hlf_download.gif"/></p>
<p align="center">Figure 1. Example of download file in HLF (x2 speed).</p>

File transfer is limited by the bandwidth of HLS itself. If we take into account that the packet generation period is `5 seconds`, then it will take about 10 seconds to complete the request-response cycle. HLS also limits the size of transmitted packets. If we assume that the limit is `8KiB`, taking into account the existing ~4KiB headers, then the transfer rate is defined as `4KiB/10s` or `410B/1s`.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ cd ./cmd/hidden_lake/applications/filesharer
$ make build # create hlf, hlf_[arch=amd64,arm64]_[os=linux,windows,darwin] and copy to ./bin
$ make run # run ./bin/hlf

> [INFO] 2023/06/03 15:30:31 HLF is running...
> ...
```

Open ports `9541` (HTTP, interface) and `9542` (HTTP, incoming).
Creates `./hlf.yml` or `./_mounted/hlf.yml` file (docker) and `./hlf.stg` or `./_mounted/hlf.stg` (docker) directory.
The directory `hlf.stg` stores all shared/loaded files. 

Default config `hlf.yml`

```yaml
settings:
  retry_num: 2
  page_offset: 10
logging:
- info
- warn
- erro
language: ENG
address:
  interface: 127.0.0.1:9541
  incoming: 127.0.0.1:9542
connection: 127.0.0.1:9572
```

Build and run with docker

```bash 
$ cd ./cmd/hidden_lake/applications/filesharer
$ make docker-build 
$ make docker-run

> [INFO] 2023/06/03 08:35:50 HLF is running...
> ...
```

```bash
$ ./hlf -path=/root
# path = path to config and storage
```

## Example

The example will involve (as well as in HLS) two nodes `node1_hlf and node2_hlf`. Both nodes are a combination of HLS and HLF, where HLF plays the role of an application and services (as shown in `Figure 3` of the HLS readme).

Build and run nodes
```bash
$ cd examples/anon_filesharing/default
$ make
```

Than open browser on `localhost:8080`. It is a `node1_hlf`. This node is a Alice.

<p align="center"><img src="_images/hlf_about.png" alt="hlf_about.png"/></p>
<p align="center">Figure 2. Home page of the HLF application.</p>

To see the another side of communication, you need to do all the same operations, but with `localhost:7070` as `node2_hlf`. This node will be Bob.

> More example images about HLF pages in the [github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/_images](https://github.com/number571/go-peer/tree/master/cmd/hidden_lake/applications/filesharer/_images "Path to HLF images")
