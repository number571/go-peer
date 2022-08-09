package friends

import "github.com/number571/go-peer/crypto/asymmetric"

type IF2F interface {
	iStatus
	iListPubKey
}

type iStatus interface {
	Switch(bool)
	Status() bool
}

type iListPubKey interface {
	InList(asymmetric.IPubKey) bool
	List() []asymmetric.IPubKey
	Append(asymmetric.IPubKey)
	Remove(asymmetric.IPubKey)
}
