package quantum

type IListPubKeyChains interface {
	AllPubKeyChains() []IPubKeyChain
	GetPubKeyChain(ISignerPubKey) (IPubKeyChain, bool)
	AddPubKeyChain(IPubKeyChain)
	DelPubKeyChain(IPubKeyChain)
}

type IPrivKeyChain interface {
	ToString() string
	GetPubKeyChain() IPubKeyChain
	GetKEMPrivKey() IKEMPrivKey
	GetSignerPrivKey() ISignerPrivKey
}

type IPubKeyChain interface {
	ToString() string
	GetKEMPubKey() IKEMPubKey
	GetSignerPubKey() ISignerPubKey
}

type IKEMPrivKey interface {
	ToBytes() []byte
	GetPubKey() IKEMPubKey
	Decapsulate([]byte) ([]byte, error)
}

type IKEMPubKey interface {
	Encapsulate() ([]byte, []byte, error)
	ToBytes() []byte
}

type ISignerPrivKey interface {
	ToBytes() []byte
	GetPubKey() ISignerPubKey
	SignBytes([]byte) []byte
}

type ISignerPubKey interface {
	ToBytes() []byte
	VerifyBytes([]byte, []byte) bool
}
