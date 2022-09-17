package database

import "github.com/number571/go-peer/modules/crypto/asymmetric"

type sRelation struct {
	fIam    asymmetric.IPubKey
	fFriend asymmetric.IPubKey
}

func NewRelation(iam, friend asymmetric.IPubKey) IRelation {
	return &sRelation{
		fIam:    iam,
		fFriend: friend,
	}
}

func (rel *sRelation) IAm() asymmetric.IPubKey {
	return rel.fIam
}

func (rel *sRelation) Friend() asymmetric.IPubKey {
	return rel.fFriend
}
