package main

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	cr "github.com/number571/gopeer/crypto"
)

type DB struct {
	ptr *sql.DB
	mtx sync.Mutex
}

type Sessions struct {
	mpn map[string]*sessionData
	mtx sync.Mutex
}

type User struct {
	Id   int
	Name string
	Pasw []byte
	Priv cr.PrivKey
}

type Email struct {
	Id         int
	SenderName string
	SenderPubl string
	Head       string
	Body       string
	Hash       string
	Time       string
}
