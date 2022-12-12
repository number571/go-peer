package friends

import "github.com/number571/go-peer/pkg/crypto/asymmetric"

type IF2F interface {
	InList(asymmetric.IPubKey) bool
	List() []asymmetric.IPubKey
	Append(asymmetric.IPubKey)
	Remove(asymmetric.IPubKey)
}
