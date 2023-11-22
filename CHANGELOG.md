# CHANGELOG

<!-- ... -->

## v1.5.22

*??? ??, ????*

### IMPROVEMENTS

- Update `cmd/hidden_lake`: append HLL service (Hidden Lake Loader)

### CHANGES

- Update `pkg/client/message`: delete proof of work
- Delete `pkg/errors`: delete package errors -> replace standard "errors" package
- Delete `pkg/file_system`: delete package file_system -> replace standard "os" package
- Delete `cmd/hidden_lake`: move network_key param from config main block to settings block
- Delete `pkg/network/message`: LoadMessage not return error

<!-- ... -->

## v1.5.21

*November 18, 2023*

### IMPROVEMENTS

- Update `pkg/storage/database`: append tryRecover function to NewKeyValueDB
- Update `pkg/network/message`: append Proof field to IMessage interface
- Update `pkg/network`: BroadcastPayload -> BroadcastMessage, (WritePayload, ReadPayload) -> (WriteMessage, ReadMessage)
- Update `pkg`: append doc.go files to all packages 

### CHANGES

- Update `pkg/stringtools`: rename -> slices, replace from pkg/ to internal/
- Update `pkg/client/queue`: delete UpdateClient method
- Update `cmd/tools/pmanager`: EOF -> EOL
- Update `pkg/client/message`: replace Payload field from SBody to SMessage 
- Update `pkg/client`: replace init check of message from custom to msg.IsValid
- Update `pkg/client/message`: move convert functions to internal/msgconv
- Update `pkg`: _examples/ -> examples/
- Update `pkg/types`: move CloseAll, StopAll functions to internal/interrupt
- Update `cmd/hidden_lake/service`: delete HandleMessage API
- Update `cmd/hidden_lake/traffic`: database (GetHashes, Load): string arg -> []byte arg
- Update `Dockerfile*`: append modifier '--platform linux/amd64' to section FROM 
- Update `cmd/hidden_lake/traffic`: interfaces with message (pkg/client/message) -> message (pkg/network/message)

### BUG FIXES

- Update `pkg/network`: rewrite inMapWithSet -> inQueueWithSet
- Update `cmd/hidden_lake/traffic`: fix HandleMessage API
- Update `cmd/hidden_lake/service`: delete field 'messages_capacity' from config

<!-- ... -->

## v1.5.20

*November 02, 2023*

### IMPROVEMENTS

- Update `cmd/hidden_lake/service`: update _daemon/install_hls.sh script, append key generation
- Update `cmd/hidden_lake/service`: append client.ResetPrivKey method
- Update `pkg/storage/database`: append check on correct input auth/encryption key
- Update `pkg`: update tests (coverage > 90%)
- Update `*_test`: change all tests to parallel actions
- Update `README.md`: append new badge - coverage

### CHANGES

- Update `cmd/hidden_lake/service`: CNetworkMaxConns 64 -> 256
- Update `examples/echo_service/prod_test`: append switch prod_1/prod_2 in Makefile with PROD param
- Update `cmd/hidden_lake`: delete jino, timeweb.cloud providers
- Update `cmd/hidden_lake/messenger`: delete auth
- Update `cmd/hidden_lake/service`: delete SetPrivKey/ResetPrivKey
- Update `cmd/hidden_lake/service`: generates priv key file
- Update `cmd/hidden_lake/*`: workSize, storageKey are can be null value
- Update `cmd/hidden_lake/service`: move backup_connections from HLM to HLS

### BUG FIXES

- Update `cmd/hidden_lake/service`: append check on size of input private key
- Update `cmd/hidden_lake/service`: update _math directory
- Update `pkg/client`: append payload check decode
- Update `pkg/client/message`: append check on unknown type
- Update `cmd/hidden_lake/service`: fix check duplicate public key in config actions

## v1.5.19

*October 05, 2023*

### IMPROVEMENTS

- Update `cmd/hidden_lake/messenger`: append _daemon scripts
- Update `pkg/network`: append return error for IHandlerF
- Update `cmd/hidden_lake/traffic`: now HLT redirect messages from producers to network/consumers
- Update `pkg/network/anonymity`: replace logbuilder's string format into internal/logger/anon
- Update `pkg/network/conn`: append ReadTimeout param for function ReadPayload

### CHANGES

- Update `pkg/network/anonymity`: delete hash field from IHandlerF
- Update `cmd/hidden_lake/messenger`: refactoring sMessage fields (timestamp, blockUID)
- Update `cmd/hidden_lake/service`: rename service headerd
- Update `cmd/hidden_lake`: refactoring running pprof service 
- Update `pkg/storage`: delete getHashing setting, check FPassword setting
- Update `pkg/storage/database`: turn on opt DisableBlockCache
- Update `pkg/network/anonymity`: update store hash of message
- Delete `cmd/hidden_lake/composite/mobile`: delete android support, delete fyne's dependency 
- Update `go.mod`: go1.17 -> go1.16 
- Update `cmd/hidden_lake/[service|traffic]`: append GetConnDeadline 

### BUG FIXES

- Update `cmd/hidden_lake/service`: fix return JSON format for '/api/config/settings'

## v1.5.18

*September 19, 2023*

### IMPROVEMENTS

- Update `cmd/hidden_lake/composite`: append builds for android/arm64 and android/amd64
- Update `cmd/hidden_lake`: append new provider - eternalhost.net
- Update `cmd/hidden_lake`: switch provider jino.ru from HLTr to HLTs
- Update `cmd/hidden_lake`: append provider timeweb.cloud as HLTr

### CHANGES

- Update `cmd/hidden_lake/messenger`: now not deleted connections from HLS config after logout
- Update `pkg/client`: move GetMessageLimit from func to method on *sClient
- Update `cmd/hidden_lake/messenger`: change sizes of buttons, card blocks in settings.html
- Update `README.md`: append 'Releases' chapter
- Update `cmd/hidden_lake/composite`: update Makefiles build/clean
- Update `examples/echo_service`: rename with_stress_test -> prod_test
- Update `examples/traffic_keeper`: append -tags=prod
- Update `test/_data`: change 4096 bit key -> 1024 bit key
- Update `pkg/client/queue`: change receiver of void messages -> random public key
- Update `internal/settings`: delete internal/settings
- Update `cmd/hidden_lake/messenger`: delete fields "message_size_bytes", "work_size_bits", "key_size_bits" from settings

### BUG FIXES

- Update `pkg/network`: fix update network key
- Update `cmd/hidden_lake/messenger`: fix state with network key

<!-- ... -->

## v1.5.17

*August 28, 2023*

### IMPROVEMENTS

- Update `pkg/crypto/entropy`: now used pbkdf2. {[Issue](https://github.com/number571/go-peer/issues/4)}
- Update `pkg/storage`: update test
- Update `cmd/hidden_lake/messenger`: append entropy check password
- Update `pkg/client,pkg/network/conn,pkg/storage`: append comments with algorithm's work
- Update `pkg/network/conn`: readPayload now return error reason

### CHANGES

- Update `cmd/hidden_lake/messenger`: append check of message size
- Update `cmd/hidden_lake/messenger`: deleted HLM<->HLM encryption throw HLS (changed threat model)
- Update `pkg/crypto/entropy`: rename interfaces/functions to keyBuilder 
- Update `cmd/hidden_lake/messenger`: change login="user", password="password" => login="username", password="hello, world!"
- Update `pkg/*/_examples`: update examples for client, network, anonymity packages
- Update `pkg/crypto/symmetric`: new cipher now create cipher.Block interface
- Update `cmd/hidden_lake/messenger/README.md`: fix urls to images
- Update `pkg/network/conn`: deleted FetchPayload method

### BUG FIXES

- Update `pkg/network/conn`: replace 4 conn.Write -> 1 conn.Write. {[Issue](https://github.com/number571/go-peer/issues/5)}
- Update `pkg/network/conn`: append hash check for msgBytes||voidBytes. {[Issue](https://github.com/number571/go-peer/issues/6)}

<!-- ... -->

## v1.5.16

*August 13, 2023*

### IMPROVEMENTS

- Update `docs`: append article decentralized_key_exchange_protocol
- Update `cmd/hidden_lake/README.md`: append to connections "characteristics", "provider" fields
- Update `cmd/hidden_lake/messenger`: append network key updater 
- Update `cmd/hidden_lake`: append http loggers to service, traffic, messenger

### CHANGES

- Update `cfgs[message_size_bytes]`: change from 4KiB to 8KiB
- Update `hidden_lake/*/init_app.go`: append trimsuffix "/" for path value
- Update `composite/mobile/service_messenger`: app.New -> app.NewWithID
- Update `cfgs[messages_capacity]`: change from (1 << 10) to (1 << 20)
- Update `pkg/client/message`: change separator from \n\n to ===
- Update `examples/anon_messenger`: change request.sh -> _request/main.go
- Update `examples/echo_service`: append _request/main.go, move request.sh -> _request/

### BUG FIXES

- Update `hidden_lake/service`: messageSize (4 << 20) -> (4 << 10)
- Update `Makefiles`: append .exe extenstion to windows compile
- Update `hidden_lake/messenger`: edit CDefaultConnectionHLSAddress -> hls_settings.CDefaultHTTPAddress
- Update `cmd/hidden_lake/service,traffic`: update README API
- Update `cmd/hidden_lake/messenger`: fix relation priv_key with HLS (append check IsMyPubKey?)
- Update `cmd/hidden_lake/messenger`: append check in state/update.go for got messages from HLT
- Update `cmd/hidden_lake/messenger`: append E2E encryption of request messages HLM <-> HLM throw HLS
- Update `pkg/client`: fix static size of messages

<!-- ... -->

## v1.5.15

*Jule 31, 2023*

### IMPROVEMENTS

- Update `hidden_lake/messenger`: append RUS language
- Update `hidden_lake/messenger`: append ESP language
- Update `hidden_lake/messenger`: append mobile/android app
- Update `hidden_lake/messenger`: append config editor (Language)
- Update `hidden_lake/messenger`: append connect to storage/backup nodes (HLT) 
- Update `hidden_lake/messenger`: append parallel load traffic from HLTs
- Update `hidden_lake/service`: append to SetPrivKey ephemeral public key
- Update `examples/anon_messenger`: append request.sh for sending text, files
- Update `hidden_lake/service`: replace LimitVoidSize from code to config
- Update `pkg/client/message`: change IMessage from JSON format to JSON/Binary ([]byte), JSON/Hex (string)

### CHANGES

- Rename `hidden_lake/_monolith`: rename _monolith/ -> composite/
- Change `hidden_lake/composite`: changed the order stop apps
- Update `theory_of_the_structure_of_hidden_systems`: updated the schemes in the algebraic model
- Delete `hidden_lake/messenger`: delete config field "traffic"
- Update `*.cfg`: :9571 -> 127.0.0.1:9571, :9582 -> 127.0.0.1:9582, ...
- Update `go.mod`: go1.16 -> go1.17 (reason: fyne/v2 used golang.org/x/sys v0.5.0)
- Update `pkg/anonymity/logbuilder`: append size of messages
- Update `hidden_lake/service`: change receiver ID from PubKey to AliasName
- Update `cfgs[message_size_bytes]`: change from 1MiB to 4KiB (also limit_void_size_bytes)

### BUG FIXES

- Update `hidden_lake/service`: fix bug with -key flag
- Update `hidden_lake/service`: change timeout read/write to queue duration
- Update `hidden_lake/service`: append check for hex_data is an IRequest?
- Update `hidden_lake/service`: fix size of the encrypted messages
- Update `pkg/client/message`: replace lax verification on size of message to strict verification
- Update `pkg`: delete all debug logs

<!-- ... -->

## v1.5.14

*Jule 24, 2023*

### IMPROVEMENTS

- Update `hidden_lake/messenger`: than got new message -> auto scroll to bottom
- Update `hidden_lake/messenger`: append support emoji's text
- Update `hidden_lake/messenger`: append constant chat with ourself
- Update `hidden_lake/messenger`: append support file transfer
- Create `hidden_lake/_monolith`: create service_messenger, service_traffic, service_traffic_messenger

### CHANGES

- Move `CONNECTIONS.md`: merge data of connections with cmd/hidden_lake's README.md 
- Change `hidden_lake/Makefile`: "composite-default: composite-build" -> "composite-default: composite-build composite-run"
- Change `hidden_lake/messenger`: move CChatLimitMessages to config value as "messages_capacity"
- Change `hidden_lake`: move app path's from internal/ to pkg/, move config path's from pkg/ to internal/

<!-- ... -->

## v1.5.13

*Jule 17, 2023*

### IMPROVEMENTS

- Append `examples`: routing to echo_service 
- Create `CONNECTIONS.md`: list of connections to HLT relayers and HLS nodes
- Update `hidden_lake/traffic`: append option 'storage'=(true|false)
- Update `network/anonymous`: change logger -> logBuilder
- Replace `hidden_lake/service`: constants (message_size_bytes, work_size_bits, key_size_bits, queue_period_ms, messages_capacity) to configs .hls, .hlt, .hlm
- Update `hidden_lake/traffic`: append check/push hash messages into database
- Create `.vscode`: append debug running options "Run Hidden Lake" and "Test Echo Service"
- Update `hidden_lake/messenger`: append onlyWritableCharacters into HandleIncomigHTTP 

### CHANGES

- Update `README.md`: delete tree/master suffix in view urls
- Update `hidden_lake/traffic`: delete redirect message to nodes from HTTP handler
- Update `hidden_lake/traffic`: append redirect message to nodes from TCP handler
- Change `examples`: replace middle_hls to middle_hlt
- Change `hidden_lake/adapters`: change recv: hlt-port -> hls-port
- Update `theory_of_the_structure_of_hidden_systems`: Append link to economic reasons

### BUG FIXES

- Update `README.md`: url with images -> _images
- Update `cmd/micro_anon`: change panic error -> print error
- Change `hidden_lake/service`: rename CLogWarnOffResponseFromService, CLogWarnResponseFromService -> CLogInfoOffResponseFromService, CLogInfoResponseFromService
- Update `hidden_lake/messenger`: replace convertToPlain -> escapeOutput function

<!-- ... -->

## v1.5.12

*Jule 08, 2023*

### IMPROVEMENTS

- Append `examples`: daemons installers
- Update `cmd`: append builds os=darwin and arch=arm64 for makefiles
- Update `cmd/micro_anon`: append makefile build to bin/
- Update `bin`: append remove-std option to makefile

### CHANGES

- Makefiles `cmd`: change CGO_ENABLED=1 -> CGO_ENABLED=0
- Tests `make test`: replace race modifier from test-run to test-race
- Package `conn_keeper`: tryConnectToAll "for range map" -> "_, ok := map by key-address"
- Package `network/conn`: append wait read deadline method
- Docs `README.md`: Append 'Calling functions/methods' to code style go-peer
- Package `storage`: Change ISettings, delete ISettings from storage/database
- Update `theory_of_the_structure_of_hidden_systems`: Append information about obfuscating routing algorithm
- Makefule `hidden_lake`: _std_cfg -> _default_cfg
- Move `images`: to projects HLS, HLM, HLT, HLA, MA
- Move `examples/_modules`: to packages network, network/anonimity, client, client/queue

### BUG FIXES

- Makefile `hidden_lake`: _std_cfg -> _std_cfgs
- Package `network`: fix deadlock BroadcastPayload on WritePayload method
- Package `network/conn`: append read/write deadline methods
- Makefiles `cmd`: fix clean commands for (arch, os) names

<!-- ... -->

## v1.5.11

*Jule 03, 2023*

### IMPROVEMENTS

- Projects `tools`: append builds with makefiles
- Update `README.md`: append info about tools applications
- Append `micro_anon`: new project with micro-anonymous network
- Append `examples`: mathematical calculations of resources consumed

### CHANGES

- Package `storage`: deleted in-memory-storage
- Project `hidden_lake/messegner`: fetch request -> broadcast request
- Test `pkg/storage/database`: delete testFailCreate
- Append `hidden_lake/_std_configs`: standart configs for HLS, HLM, HLT
- Append `hidden_lake`: composite build and run applications

### BUG FIXES

- Project `hidden_lake`: makefile docker-default: build -> docker-build

<!-- ... -->

## v1.5.10

*June 25, 2023*

### IMPROVEMENTS

- Projects `hidden_lake`: append build/run with docker
- Update `README.md`: append build and run instructions
- Example `_cmd/echo_service`: append docker-compose examples
- Example `_cmd/anon_messenger`: append docker-compose examples
- Projects `tools`: append encryptor application
- Projects `tools`: append pmanager application
- Projects `tools`: refactoring storage application
- Article `abstract_anonymous_networks`: append info about periods of state in entropy increase networks

### CHANGES

- Project `hidden_lake/messenger`: deleted HLS part from hlm_m application
- Package `encoding`: Serialize function with option (indent/not indent)
- Project `hidden_lake/messenger`: deleted HLS part from hlm application (build and run)
- Package `pkg/storage/database`: replace sqlite3 to leveldb
- Project `mobile_applications`: deleted mobile applications HLS, HLM, HLT
- Module `go.mod`: decrease version from 1.17 to 1.16
- Update `README.md`: append installation, requirements
- Project `hidden_lake/messenger`: rename interface and methods in IState -> IStateManager
- Directory `tools`: moved to cmd/tools

### BUG FIXES

- Project `hidden_lake/messenger`: append checks pStateManager.GetWrapperDB().Get() on nil
- Package `crypto/entropy`: fix range hashes with one input data

<!-- ... -->

## v1.5.9

*June 1, 2023*

### IMPROVEMENTS

- Makefiles `hidden_lake`: append cross-compile (linux, windows, android)
- Makefiles `hidden_lake`: append all-build, all-clean options
- Package `anonymity`: create adapter to interface payload.IPayload
- Article `theory_of_the_structure_of_hidden_systems`: append new section (introduction / economic reasons)
- Functions `NewSettings`: defaultValues -> mustNotNull (strict validation)
- Package `storage`: NewCryptoStorage now get ISettings
- Update `README.md`: append badgers, append table in HLS
- Package `anonymity`: IHandlerF now return ([]byte, error)
- Project `hidden_lake/service`: append response package
- Project `hidden_lake/service`: append response http headers
- Package `pkg/errors`: all packages now use stack errors

### CHANGES

- Package `anonymity`: deleted limit chars in serviceName logger.ILogger
- Package `anonymity`: BroadcastPayload corresponds to the FetchPayload interface
- Service `hidden_lake`: rename CHeaderHLS -> CServiceMask
- Package `storage`: deleted memory storage
- Package `database`: deleted settings value - salt
- Package `test/_data`: reduced the dependency of other packages on this
- Package `internal/api`: JSON (return, result) -> http code, result
- Package `encoding`: Serialize function now Marshal (not MarshalIndent) data
- Example `images`: replace HLS, HLA gifs with requests

### BUG FIXES

- Gitignore `HLT`: change m_hls.apk -> m_hlt.apk 
- Application `HLM_M`: change `gAppHLS == nil || gAppHLM == nil` -> `gAppHLS == nil && gAppHLM == nil`

<!-- ... -->

## v1.5.8

*May 8, 2023*

### INIT
