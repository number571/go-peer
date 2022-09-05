package main

import (
	hms_database "github.com/number571/go-peer/cmd/hms/database"
	"github.com/number571/go-peer/modules/action"
	"github.com/number571/go-peer/modules/client"
)

var (
	gActions action.IActions
	gWrapper iWrapper
	gClient  client.IClient
	gDB      hms_database.IKeyValueDB
)
