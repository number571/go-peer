package database

import "github.com/number571/go-peer/pkg/crypto/asymmetric"

var (
	_ IRelation = &sRelation{}
)

type sRelation struct {
	fIAm    asymmetric.IPubKey
	fFriend asymmetric.IPubKey
}

func NewRelation(iam, friend asymmetric.IPubKey) IRelation {
	return &sRelation{
		fIAm:    iam,
		fFriend: friend,
	}
}

func (r *sRelation) IAm() asymmetric.IPubKey {
	return r.fIAm
}

func (r *sRelation) Friend() asymmetric.IPubKey {
	return r.fFriend
}
