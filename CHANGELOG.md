# Changelog

<!-- ... -->

## ~v1.5.10

*??? ??, ????*

In progress...

### IMPROVEMENTS

- Projects `hidden_lake`: append build/run with docker
- Text `README`: append build and run instructions
- Example `_cmd/echo_service`: append docker-compose examples
- Example `_cmd/anon_messenger`: append docker-compose examples

### CHANGES

- Project `hidden_lake/messenger`: deleted HLS part from hlm_m application
- Package `encoding`: Serialize function with option (indent/not indent)
- Project `hidden_lake/messenger`: deleted HLS part from hlm application (build and run)
- Package `pkg/storage/database`: replace sqlite3 to leveldb
- Project `mobile_applications`: deleted mobile applications HLS, HLM, HLT
- Module `go.mod`: decrease version from 1.17 to 1.16

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
