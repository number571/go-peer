package asymmetric

type IListPubKeyChains interface {
	AllPubKeyChains() []IPubKeyChain
	GetPubKeyChain(ISignPubKey) (IPubKeyChain, bool)
	AddPubKeyChain(IPubKeyChain)
	DelPubKeyChain(IPubKeyChain)
}

type IPrivKeyChain interface {
	ToString() string
	GetPubKeyChain() IPubKeyChain
	GetKEncPrivKey() IKEncPrivKey
	GetSignPrivKey() ISignPrivKey
}

type IPubKeyChain interface {
	ToString() string
	GetKEncPubKey() IKEncPubKey
	GetSignPubKey() ISignPubKey
}

type IKEncPrivKey interface {
	ToBytes() []byte
	GetPubKey() IKEncPubKey
	Decapsulate([]byte) ([]byte, error)
}

type IKEncPubKey interface {
	Encapsulate() ([]byte, []byte, error)
	ToBytes() []byte
}

type ISignPrivKey interface {
	ToBytes() []byte
	GetPubKey() ISignPubKey
	SignBytes([]byte) []byte
}

type ISignPubKey interface {
	ToBytes() []byte
	VerifyBytes([]byte, []byte) bool
}
