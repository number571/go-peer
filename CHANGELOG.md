# Changelog

<!-- ... -->

## ~v1.5.9

*May 8, 2023*

Init.

### IMPROVEMENTS

- Makefiles `hidden_lake`: append cross-compile (linux, windows, android)
- Makefiles `hidden_lake`: append all-build, all-clean options
- Package `anonymity`: create adapter to interface payload.IPayload
- Article `theory_of_the_structure_of_hidden_systems`: append new section (introduction / economic reasons)
- NewSettings functions: defaultValues -> mustNotNull (strict validation)
- Package `storage`: NewCryptoStorage now get ISettings

### CHANGES

- Package `anonymity`: deleted limit chars in serviceName logger.ILogger
- Package `anonymity`: BroadcastPayload corresponds to the FetchPayload interface
- Service `hidden_lake`: rename CHeaderHLS -> CServiceMask
- Package `storage`: deleted memory storage
- Package `database`: deleted settings value - salt
- Package `test/_data`: reduced the dependency of other packages on this

### BUG FIXES

- Gitignore (HLT): change m_hls.apk -> m_hlt.apk 

<!-- ... -->

## v1.5.8

*May 8, 2023*

Init.

### IMPROVEMENTS

- Template

### CHANGES

- Template

### BUG FIXES

- Template
