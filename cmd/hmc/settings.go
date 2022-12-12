package main

import (
	"github.com/number571/go-peer/cmd/hmc/action"
	hms_database "github.com/number571/go-peer/cmd/hms/database"
	hls_client "github.com/number571/go-peer/pkg/client"
)

var (
	gActions action.IActions
	gWrapper iWrapper
	gClient  hls_client.IClient
	gDB      hms_database.IKeyValueDB
)
