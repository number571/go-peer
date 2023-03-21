# Docs

## Theoretical works

### 1. Main

1. [Theory of the structure of hidden systems](https://github.com/number571/go-peer/blob/master/docs/theory_of_the_structure_of_hidden_systems.pdf "TotSoHS")
2. [Monolithic cryptographic protocol](https://github.com/number571/go-peer/blob/master/docs/monolithic_cryptographic_protocol.pdf "MCP")
3. [Abstract anonymous networks](https://github.com/number571/go-peer/blob/master/docs/abstract_anonymous_networks.pdf "AAN")

### 2. Habr

1. [Hidden Lake Service](https://habr.com/ru/post/696504/ "Habr HLS")
2. [Hidden Lake Messenger](https://habr.com/ru/post/701488/ "Habr HLM")
3. [Hidden Lake Traffic](https://habr.com/ru/post/717184/ "Habr HLT")
4. [Hidden Lake Adapters](https://habr.com/ru/post/720544/ "Habr HLA")

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

The name of the global variables must begin with the prefix 'g' or 'G'.
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

The name of the global constants, variables, structures, fields, interfaces in the test enviroment must begin with prefix 't' or 'T'.
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
); err != nil { /* not allowed multiply args in if block */
    // ...
}
```

### 4. Interface declaration

When a type is bound to an interface, it must be explicitly specified like this.
```go
var (
	_ types.ICommand = &sApp{}
)
```
