package database

import (
	"database/sql"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"

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
		return nil, err
	}

	if _, err := db.Exec(cTabelKV); err != nil {
		return nil, err
	}

	var saltValue string
	row := db.QueryRow("SELECT value FROM kv WHERE key = $1", cSaltKey)
	if err := row.Scan(&saltValue); err != nil {
		saltValue = encoding.HexEncode(random.NewStdPRNG().GetBytes(symmetric.CAESKeySize))
		_, err := db.Exec("REPLACE INTO kv (key, value) VALUES ($1,$2)", cSaltKey, saltValue)
		if err != nil {
			return nil, err
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
	return err
}

func (p *sSQLiteDB) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	var encValue string
	row := p.fDB.QueryRow("SELECT value FROM kv WHERE key = $1", encoding.HexEncode(p.tryHash(pKey)))
	if err := row.Scan(&encValue); err != nil {
		return nil, err
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
	return err
}

func (p *sSQLiteDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fDB.Close()
}

func (p *sSQLiteDB) tryHash(pKey []byte) []byte {
	if !p.fSettings.GetHashing() {
		return pKey
	}
	return hashing.NewHMACSHA256Hasher(p.fSalt, pKey).ToBytes()
}
