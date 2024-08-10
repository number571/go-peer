# HLC

> Hidden Lake Composite

<img src="_images/hlc_logo.png" alt="hlc_logo.png"/>

The `Hidden Lake Composite` combines several HL type's services into one application using startup config.

## Installation

```bash
$ go install github.com/number571/go-peer/cmd/hidden_lake/composite/cmd/hlc@latest
```

## How it works

The application HLC includes the download of all Hidden Lake services, and runs only the configurations selected by names in the file. The exact names of the services can be found in their `pkg/settings/settings.go` configuration files.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ cd ./cmd/hidden_lake/composite
$ make build # create hlc, hlc_[arch=amd64,arm64]_[os=linux,windows,darwin] and copy to ./bin
$ make run # run ./bin/hlc

> [INFO] 2023/12/03 02:12:51 HLC is running...
> ...
```

Creates `./hlc.yml` file.

Default config `hlc.yml`

```yaml
logging:
- info
- warn
- erro
services:
- hidden-lake-service
- hidden-lake-messenger
- hidden-lake-filesharer
```

## Running options

```bash
$ ./hlc -path=/root -parallel=1
# path     = path to config, database, key files
# parallel = num of parallel functions for PoW algorithm
```

## Config structure

```
"logging"  Enable loggins in/out actions in the network
"services" Names of Hidden Lake services 
```
