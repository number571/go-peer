package database

import (
	"database/sql"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

const (
	cSaltKey = "__SALT__"
	cTabelKV = `
		CREATE TABLE IF NOT EXISTS kv (
			key text unique, 
			value text
		);
		`
)

var (
	_ IKeyValueDB = &sSQLiteDB{}
)

type sSQLiteDB struct {
	fMutex    sync.Mutex
	fSalt     []byte
	fDB       *sql.DB
	fSettings ISettings
	fCipher   symmetric.ICipher
}

func NewSQLiteDB(pSett ISettings) (IKeyValueDB, error) {
	db, err := sql.Open("sqlite3", pSett.GetPath())
	if err != nil {
		return nil, errors.WrapError(err, "open database")
	}

	if _, err := db.Exec(cTabelKV); err != nil {
		return nil, errors.WrapError(err, "insert KV table")
	}

	var saltValue string
	row := db.QueryRow("SELECT value FROM kv WHERE key = $1", cSaltKey)
	if err := row.Scan(&saltValue); err != nil {
		saltValue = encoding.HexEncode(random.NewStdPRNG().GetBytes(symmetric.CAESKeySize))
		_, err := db.Exec("REPLACE INTO kv (key, value) VALUES ($1,$2)", cSaltKey, saltValue)
		if err != nil {
			return nil, errors.WrapError(err, "insert salt into database")
		}
	}

	return &sSQLiteDB{
		fSalt:     encoding.HexDecode(saltValue),
		fDB:       db,
		fSettings: pSett,
		fCipher:   symmetric.NewAESCipher(pSett.GetCipherKey()),
	}, nil
}

func (p *sSQLiteDB) GetSettings() ISettings {
	return p.fSettings
}

func (p *sSQLiteDB) Set(pKey []byte, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	_, err := p.fDB.Exec(
		"REPLACE INTO kv (key, value) VALUES ($1,$2)",
		encoding.HexEncode(p.tryHash(pKey)),
		encoding.HexEncode(doEncrypt(p.fCipher, pValue)),
	)
	if err != nil {
		return errors.WrapError(err, "insert key/value to database")
	}
	return nil
}

func (p *sSQLiteDB) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	var encValue string
	row := p.fDB.QueryRow("SELECT value FROM kv WHERE key = $1", encoding.HexEncode(p.tryHash(pKey)))
	if err := row.Scan(&encValue); err != nil {
		return nil, errors.WrapError(err, "read value by key")
	}

	return tryDecrypt(
		p.fCipher,
		encoding.HexDecode(encValue),
	)
}

func (p *sSQLiteDB) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	_, err := p.fDB.Exec(
		"DELETE FROM kv WHERE key = $1",
		encoding.HexEncode(p.tryHash(pKey)),
	)
	if err != nil {
		return errors.WrapError(err, "delete value by key")
	}
	return nil
}

func (p *sSQLiteDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Close(); err != nil {
		return errors.WrapError(err, "close database")
	}
	return nil
}

func (p *sSQLiteDB) tryHash(pKey []byte) []byte {
	if !p.fSettings.GetHashing() {
		return pKey
	}
	return hashing.NewHMACSHA256Hasher(p.fSalt, pKey).ToBytes()
}
