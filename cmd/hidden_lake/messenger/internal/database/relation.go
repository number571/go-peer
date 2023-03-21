package database

import "github.com/number571/go-peer/pkg/crypto/asymmetric"

var (
	_ IRelation = &sRelation{}
)

type sRelation struct {
	fIAm    asymmetric.IPubKey
	fFriend asymmetric.IPubKey
}

func NewRelation(pIAm, pFriend asymmetric.IPubKey) IRelation {
	return &sRelation{
		fIAm:    pIAm,
		fFriend: pFriend,
	}
}

func (p *sRelation) IAm() asymmetric.IPubKey {
	return p.fIAm
}

func (p *sRelation) Friend() asymmetric.IPubKey {
	return p.fFriend
}
