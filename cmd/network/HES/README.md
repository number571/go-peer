# HES

> Hidden email service. Version 1.1.11.

### Home page
<img src="/cmd/HES/userside/images/HES1.png" alt="HomePage"/>

### Characteristics
1. End to end encryption;
2. Supported tor connections;
3. Symmetric algorithm: AES-CBC;
4. Asymmetric algorithm: RSA-OAEP, RSA-PSS;
5. Hash function: SHA256;

### Account page
<img src="/cmd/HES/userside/images/HES4.png" alt="AccountPage"/>

### Compile and run
```
$ make
> go build client.go gconsts.go cdatabase.go cmodels.go csessions.go
> go build server.go gconsts.go sdatabase.go sconfig.go
$ ./server -open="localhost:8080" &
$ ./client -open="localhost:7545"
```

### List of emails page
<img src="/cmd/HES/userside/images/HES7.png" alt="ListOfEmailsPage"/>

### DB and CFG files
> Database and config files are creates when the application starts.

#### Server side db (server.db)
```sql
/* recv = hash(public_key) */
/* hash = hash(data) */
/* data = encrypt(email) */
CREATE TABLE IF NOT EXISTS emails (
	id      INTEGER,
	recv    VARCHAR(255),
	hash    VARCHAR(255) UNIQUE,
	data    TEXT,
	addtime DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);
```

#### Server side cfg (server.cfg)
```go
type CFG struct {
	Pasw  string      `json:"pasw"`
	Conns [][2]string `json:"conns"`
}
```

#### Client side db (client.db)
```sql
/* !key_pasw = hash(password, salt)^25 */
/* hashn     = hash(nickname) */
/* hashp     = hash(!key_pasw, nickname) */
/* priv      = encrypt[!key_pasw](private_key) */
CREATE TABLE IF NOT EXISTS users (
	id   INTEGER,
	f2f  BOOLEAN,
	hashn VARCHAR(255) UNIQUE,
	hashp VARCHAR(255),
	salt VARCHAR(255),
	priv TEXT,
	PRIMARY KEY(id)
);
/* hashn = hash(nickname, !key_pasw) */
/* hashp = hash(public_key, !key_pasw) */
/* name  = encrypt[!key_pasw](nickname) */
/* publ  = encrypt[!key_pasw](public_key) */
CREATE TABLE IF NOT EXISTS contacts (
	id      INTEGER,
	id_user INTEGER,
	hashn   VARCHAR(255) UNIQUE,
	hashp   VARCHAR(255) UNIQUE,
	name    NVARCHAR(255),
	publ    TEXT,
	PRIMARY KEY(id),
	FOREIGN KEY(id_user) REFERENCES users(id) ON DELETE CASCADE
);
/* hash = hash(host, !key_pasw) */
/* host = encrypt[!key_pasw](host) */
/* pasw = encrypt[!key_pasw](pasw) */
CREATE TABLE IF NOT EXISTS connects (
	id      INTEGER,
	id_user INTEGER,
	hash    VARCHAR(255) UNIQUE,
	host    VARCHAR(255),
	pasw    VARCHAR(255),
	PRIMARY KEY(id),
	FOREIGN KEY(id_user) REFERENCES users(id) ON DELETE CASCADE
);
/* hash    = hash(pack_hash, !key_pasw) */
/* spubl   = encrypt[!key_pasw](public_key) */
/* sname   = encrypt[!key_pasw](nickname) */
/* head    = encrypt[!key_pasw](title) */
/* body    = encrypt[!key_pasw](message) */
/* addtime = encrypt[!key_pasw](time_rec) */
CREATE TABLE IF NOT EXISTS emails (
	id      INTEGER,
	id_user INTEGER,
	deleted BOOLEAN DEFAULT 0,
	hash    VARCHAR(255) UNIQUE,
	spubl   TEXT,
	sname   NVARCHAR(255),
	head    NVARCHAR(255),
	body    TEXT,
	addtime TEXT,
	PRIMARY KEY(id),
	FOREIGN KEY(id_user) REFERENCES users(id) ON DELETE CASCADE
);
```

### Email page
<img src="/cmd/HES/userside/images/HES8.png" alt="EmailPage"/>
