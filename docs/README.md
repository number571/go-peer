# Docs

## Material

### 1. Papers

1. [Theory of the structure of hidden systems](https://github.com/number571/go-peer/blob/master/docs/theory_of_the_structure_of_hidden_systems.pdf "TotSoHS")
2. [Monolithic cryptographic protocol](https://github.com/number571/go-peer/blob/master/docs/monolithic_cryptographic_protocol.pdf "MCP")
3. [Abstract anonymous networks](https://github.com/number571/go-peer/blob/master/docs/abstract_anonymous_networks.pdf "AAN")
4. [Decentralized key exchange protocol](https://github.com/number571/go-peer/blob/master/docs/decentralized_key_exchange_protocol.pdf "DKEP")

Also, the composition of these works can be found in the book [The general theory of anonymous communications](https://ridero.ru/books/obshaya_teoriya_anonimnykh_kommunikacii/). This book can be purchased in a tangible form on the [Ozon](https://www.ozon.ru/product/obshchaya-teoriya-anonimnyh-kommunikatsiy-vtoroe-izdanie-kovalenko-a-g-1193224608/) and [Wildberries](https://www.wildberries.ru/catalog/177390621/detail.aspx) marketplaces. You can download the book in digital form for free [here](https://github.com/number571/go-peer/blob/master/docs/general_theory_of_anonymous_communications.pdf).

### 2. Presentations

1. [Grok anonymity](https://github.com/number571/go-peer/blob/master/docs/grok_anonymity.pdf "Presentation GA")

### 3. Videos

1. [DC/EI/QB networks & HLM](https://www.youtube.com/watch?v=o2J6ewvBKmg)

## Code style

In the course of editing the project, some code styles may be added, some edited. Therefore, the current state of the project may not fully adhere to the code style, but you need to strive for it.

### 1. Prefixes

The name of the global constants must begin with the prefix 'c' (internal) or 'C' (external).
```go
const (
    cInternalConst = 1
    CExternalConst = 2
)
```

The name of the global variables must begin with the prefix 'g' (internal) or 'G' (external). The exception is errors with the prefix 'err' or 'Err'.
```go
var (
    gInternalVariable = 1
    GExternalVariable = 2
)
```

The name of the global structs must begin with the prefix 's' (internal) or 'S' (external). Also fields in the structure must begin with the prefix 'f' or 'F'.
```go
type (
    sInternalStruct struct{
        fInternalField int 
    }
    SExternalStruct struct{
        FExternalField int
    }
)
```

The name of the global interfaces must begin with the prefix 'i' (internal) or 'I' (external). Also type functions must begin with the prefix 'i' or 'I'.
```go
type (
    iInternalInterface interface{}
    IExternalInterface interface{}
)

type (
    iInternalFunc func()
    iExternalFunc func()
)
```

The name of the function parameters must begin with the prefix 'p'. Also method's object must be equal 'p'. The exception of this code style is test files.
```go
func f(pK, pV int) {}
func (p *sObject) m() {}
```

The name of the global constants, variables, structures, fields, interfaces in the test environment must begin with prefix 't' (internal) or 'T' (external).
```go
const (
    tcInternalConst = 1
    TcExternalConst = 2
)

var (
    tgInternalVariable = 1
    TgExternalVariable = 2
)

type (
    tsInternalStruct struct{
        tfInternalField int 
    }
    TsExternalStruct struct{
        TfInternalField int 
    }
)

type (
    tiInternalInterface interface{}
    TiExternalInterface interface{}
)

type (
    tiInternalFunc func()
    TiExternalFunc func()
)
```

### 2. Function / Methods names

Functions and methods should consist of two parts, where the first is a verb, the second is a noun. Standart names: Get, Set, Add, Del and etc. Example
```go
type IClient interface {
	GetIndex() (string, error)

	GetPubKey() (asymmetric.IPubKey, error)
	SetPrivKey(asymmetric.IPrivKey) error

	GetOnlines() ([]string, error)
	DelOnline(string) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(string, asymmetric.IPubKey) error
	DelFriend(string) error

	GetConnections() ([]string, error)
	AddConnection(string) error
	DelConnection(string) error

	BroadcastRequest(asymmetric.IPubKey, request.IRequest) error
	FetchRequest(asymmetric.IPubKey, request.IRequest) ([]byte, error)
}
```

### 3. If blocks

The following is allowed.
```go
if err := f(); err != nil {
    // ...
}

err := g(
    a,
    b,
    c,
)
if err != nil {
    // ...
}
```

The following is not allowed.
```go
if v {
    // ...
} else { /* eradicate the 'else' block */
    // ...
}

err := f() /* may be in if block */
if err != nil {
    // ...
}

if err := g(
    a,
    b,
    c,
); err != nil { /* not allowed multiply line-args in if block */
    // ...
}
```

### 4. Interface declaration

When a type is bound to an interface, it must be explicitly specified like this.
```go
var (
	_ types.IRunner = &sApp{}
)
```

### 5. Calling functions/methods

External simple getter functions/methods should not be used inside the package.
```go
func (p *sObject) GetSettings() ISettings {
	return p.fSettings
}
func (p *sObject) GetValue() IValue {
    p.fMutex.Lock()
    defer p.fMutex.Unlock()

    return p.fValue
}
...
func (p *sObject) DoSomething() {
	_ = p.fSettings // correct
    _ = p.GetSettings() // incorrect

    // incorrect
    p.fMutex.Lock()
    _ = p.fValue
    p.fMutex.Unlock()

    _ = p.GetValue() // correct
}
```

### 6. Args/Returns interfaces

It is not allowed to use global structures in function arguments or when returning. Interfaces should be used instead of structures.

The following is allowed.
```go
func doObject(_ IObject) {}
func newObject() IObject {}
```

The following is not allowed.
```go
func doObject(_ *SObject) {}
func newObject() *SObject {}
```
