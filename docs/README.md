# Docs

## Theoretical works

### 1. Main

1. [Theory of the structure of hidden systems](https://github.com/number571/go-peer/blob/master/docs/theory_of_the_structure_of_hidden_systems.pdf "TotSoHS")
2. [Monolithic cryptographic protocol](https://github.com/number571/go-peer/blob/master/docs/monolithic_cryptographic_protocol.pdf "MCP")
3. [Abstract anonymous networks](https://github.com/number571/go-peer/blob/master/docs/abstract_anonymous_networks.pdf "AAN")
4. [Decentralized key exchange protocol](https://github.com/number571/go-peer/blob/master/docs/decentralized_key_exchange_protocol.pdf "DKEP")
5. [The Hidden Lake anonymous network](https://github.com/number571/go-peer/blob/master/docs/hidden_lake_anonymous_network.pdf "HLAN")

Also, the composition of these works can be found in the book [The general theory of anonymous communications](https://ridero.ru/books/obshaya_teoriya_anonimnykh_kommunikacii/). This book can be purchased in a tangible form on the [Ozon](https://www.ozon.ru/product/obshchaya-teoriya-anonimnyh-kommunikatsiy-vtoroe-izdanie-kovalenko-a-g-1193224608/) and [Wildberries](https://www.wildberries.ru/catalog/177390621/detail.aspx) marketplaces. You can download the book in digital form for free [here](https://github.com/number571/go-peer/blob/master/docs/general_theory_of_anonymous_communications.pdf).

### 2. Presentations

1. [Hidden Lake](https://github.com/number571/go-peer/blob/master/docs/hidden_lake_anonymous_network.pdf "Presentation HLAN")
2. [Grok anonymity](https://github.com/number571/go-peer/blob/master/docs/grok_anonymity.pdf "Presentation GA")

### 3. Habr

1. [Hidden Lake Service](https://habr.com/ru/post/696504/ "Habr HLS")
2. [Hidden Lake Messenger](https://habr.com/ru/post/701488/ "Habr HLM")
3. [Hidden Lake Traffic](https://habr.com/ru/post/717184/ "Habr HLT")
4. [Hidden Lake Adapters](https://habr.com/ru/post/720544/ "Habr HLA")
5. [Micro-Anonymous Network](https://habr.com/ru/articles/745256/ "Habr MA")
6. [Entropy Increase Networks](https://habr.com/ru/articles/743630/ "Habr EIN")
7. [Create node in the Hidden Lake](https://habr.com/ru/articles/765464/ "Habr HL")
8. [Hidden Lake Filesharer](https://habr.com/ru/articles/789968/ "Habr HLF")
9. [Hidden Lake Remoter](https://habr.com/ru/articles/830130/ "Habr HLR")
10. [Hidden Lake Adapter Chatingar](https://habr.com/ru/articles/800965/ "Habr HLA_chatingar")

### 4. Videos

1. [Hidden Lake Messenger](https://www.youtube.com/watch?v=o2J6ewvBKmg "Video HLM")
2. [Create node in the Hidden Lake](https://www.youtube.com/watch?v=GBmQfy16AuA&t=18s "Video HL")

### 5. Schemes

1. [Hidden Lake proto](https://github.com/number571/go-peer/blob/master/docs/hidden_lake_message_layers.svg "HLp")
2. [Hidden Lake modes](https://github.com/number571/go-peer/blob/master/docs/hidden_lake_launch_modes.svg "HLm")

## Default applications ports

1. HLS  = 957x
2. HLT  = 958x
3. HLM  = 959x
4. HLL  = 956x
5. HLE  = 955x
6. HLF  = 954x
7. HLR  = 953x
8. HLD  = 952x

## Code style go-peer

In the course of editing the project, some code styles may be added, some edited. Therefore, the current state of the project may not fully adhere to the code style, but you need to strive for it.

### 1. Prefixes

The name of the global constants must begin with the prefix 'c' or 'C'.
```go
const (
    cInternalConst = 1
    CExternalConst = 2
)
```

The name of the global variables must begin with the prefix 'g' or 'G'. The exception is errors with the prefix 'err' or 'Err'.
```go
var (
    gInternalVariable = 1
    GExternalVariable = 2
)
```

The name of the global structs must begin with the prefix 's' or 'S'. Also fields in the structure must begin with the prefix 'f' or 'F'.
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

The name of the global interfaces must begin with the prefix 'i' or 'I'. Also type functions must begin with the prefix 'i' or 'I'.
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

The name of the function parameters must begin with the prefix 'p' or 'P'. Also method's object must be equal 'p'. The exception of this code style is test files.
```go
func f(pK, pV int) {}
func (p *sObject) m() {}
```

The name of the global constants, variables, structures, fields, interfaces in the test environment must begin with prefix 't' or 'T'.
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
