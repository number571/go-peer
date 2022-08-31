package main

import (
	hms_database "github.com/number571/go-peer/cmd/hms/database"
	"github.com/number571/go-peer/modules/action"
	"github.com/number571/go-peer/modules/client"
)

const (
	cSeparator    = "=============================================="
	cCountInPage  = 10 // count of messages in one page
	cWorkSize     = 20
	cListLenTitle = 50                 // length of title in list
	cAKeySize     = 4096               // size of asymmetric key
	cReceiveSize  = 8192               // count of messages from server
	cHeadPayload  = 0x5710017500000001 // head of payload
)

var (
	gActions action.IActions
	gWrapper iWrapper
	gClient  client.IClient
	gDB      hms_database.IKeyValueDB
)
