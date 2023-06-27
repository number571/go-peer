# Changelog

<!-- ... -->

## v1.5.11

*??? ??, ????*

In progress...

### IMPROVEMENTS

- Projects `tools`: append builds with makefiles
- Update `README.md`: append info about tools applications

### CHANGES

- Package `storage`: deleted in-memory-storage

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
- Directory `tools`: removed to cmd/tools

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
