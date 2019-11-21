/*
Package for create centralized, decentralized and distributed networks.

=======================================================================
SECTION <SETTINGS>:
=======================================================================
Constants: {
    BASE64_SHA256_BYTES: 24
    DATA_SIZE: 3
}

Variables: {
    TEMPLATE: "0.0.0.0"
    CLIENT_NAME: "[CLIENT]"
    SEPARATOR: "\000\001\007\005\000"
    END_BYTES: "\000\005\007\001\000"
    NETWORK_NAME: "GENESIS"
    TITLE_CONNECT: "[TITLE:CONNECT]"
    TITLE_REDIRECT: "[TITLE:REDIRECT]"
    MODE_READ: "[MODE:READ]"
    MODE_SAVE: "[MODE:SAVE]"
    MODE_TEST: "[MODE:TEST]"
    MODE_REMV: "[MODE:REMV]"
    MODE_MERG: "[MODE:MERG]"
    MODE_DISTRIB: "[MODE:DISTRIB]"
    MODE_DECENTR: "[MODE:DECENTR]"
    MODE_READ_MERG: "[MODE:READ][MODE:MERG]"
    MODE_SAVE_MERG: "[MODE:SAVE][MODE:MERG]"
    MODE_DISTRIB_READ: "[MODE:DISTRIB][MODE:READ]"
    MODE_DISTRIB_SAVE: "[MODE:DISTRIB][MODE:SAVE]"
    DEFAULT_HMAC_KEY: []byte("DEFAULT-HMAC-KEY")
    PACK_TIME: 10
    WAIT_TIME: 5
    ROUTE_NUM: 3
    BUFF_SIZE: 1 << 9
    PID_SIZE: 1 << 3
    SESSION_SIZE: 1 << 5
    MAXSIZE_ADDRESS: 1 << 5
    MAXSIZE_PACKAGE: 1 << 10 << 10
    CLIENT_NAME_SIZE: 1 << 4
    IS_NOTHING: true
    IS_DISTRIB: false
    IS_DECENTR: false
    HAS_CRYPTO: false
    HAS_ROUTING: false
    HAS_FRIENDS: false
    CRYPTO_SPEED: false
    HANDLE_ROUTING: false
}

Functions: {
    SettingsSet(settings SettingsType) []uint8
    SettingsGet(key string) interface{}
}
=======================================================================

=======================================================================
SECTION <NETWORK>:
=======================================================================
Functions: {
    NewNode(addr string) *Node
}

Methods: {
    (node *Node) Open() *Node
    (node *Node) Run(handleInit func(*Node), handleServer func(*Node, *Package), handleClient func(*Node, []string)) *Node
    (node *Node) ReadOnly(types ReadonlyType) *Node
    (node *Node) IsHidden(addr string) bool
    (node *Node) IsHandle(addr string) bool
    (node *Node) IsNode(addr string) bool
    (node *Node) IsMyAddress(addr string) bool
    (node *Node) IsMyHashname(hashname string) bool
    (node *Node) AppendToAccessList(access AccessType, addresses ...string)
    (node *Node) DeleteFromAccessList(addresses ...string) *Node
    (node *Node) InAccessList(access AccessType, addr string) bool
    (node *Node) InConnections(addr string) bool
    (node *Node) CreateRedirect(pack *Package) *Package
    (node *Node) Send(pack *Package) *Node
    (node *Node) SendRedirect(pack *Package) *Node
    (node *Node) SendInitRedirect(pack *Package) *Node
    (node *Node) SendToAll(pack *Package) *Node
    (node *Node) SendToAllWithout(pack *Package, sender string) *Node
    (node *Node) MergeConnect(addr string) *Node
    (node *Node) HiddenConnect(addr string) *Node
    (node *Node) ConnectToList(list ...interface{}) *Node
    (node *Node) Connect(data ...string) *Node
    (node *Node) Disconnect(addresses ...string) *Node
    (node *Node) GetConnections(relation RelationType) []string
    (node *Node) Close() *Node
}
=======================================================================

=======================================================================
SECTION <CRYPTO>:
=======================================================================
Functions: {
    ParsePrivate(privData string) *rsa.PrivateKey
    ParsePublic(pubData string) *rsa.PublicKey
    Sign(priv *rsa.PrivateKey, data []byte) []byte
    Verify(pub *rsa.PublicKey, data, sign []byte) error
    HMAC(fHash func([]byte) []byte, data []byte, key []byte) []byte
    GenerateSessionKey(max int) []byte
    EncryptRSA(pub *rsa.PublicKey, data []byte) []byte
    DecryptRSA(priv *rsa.PrivateKey, data []byte) []byte
    EncryptAES(key, data []byte) []byte
    DecryptAES(key, data []byte) []byte
    StringPublic(pub *rsa.PublicKey) string
    StringPrivate(priv *rsa.PrivateKey) string
}

Methods: {
    (node *Node) StringPublic() string
    (node *Node) StringPrivate() string
    (node *Node) ParsePrivate(privData string) *Node
    (node *Node) SetPrivate(priv *rsa.PrivateKey) *Node
    (node *Node) GeneratePrivate(bits int) *Node
    (node *Node) DecryptRSA(data []byte) []byte
    (node *Node) Sign(data []byte) []byte
}
=======================================================================

=======================================================================
SECTION <MODELS>:
=======================================================================
Types: {
    SettingsType map[string]interface{}
    ReadonlyType uint8
    RelationType uint8
    AccessType uint8
}

Constants: {
    ReadAll: 0
    ReadNodes: 1
    ReadClients: 2
    RelationAll: 0
    RelationNode: 1
    RelationHandle: 2
    RelationHidden: 3
    AccessDenied: 0
    AccessAllowed: 1
}

Structures: {
    Node: {
        Hashname string
        Keys Keys
        Setting Setting
        Address Address
        Network Network
    }

    Setting: {
        Mutex *sync.Mutex
        ReadOnly ReadonlyType
        Listen net.Listener
        HandleServer func(*Node, *Package)
        TestConnections map[string]bool
    }

    Keys: {
        Private *rsa.PrivateKey
        Public *rsa.PublicKey
    }

    Network: {
        Addresses map[string]string
        AccessList map[string]AccessType
        Connections map[string]*Connect
    }

    Connect: {
        Relation RelationType
        Hashname string
        Session []byte
        Public *rsa.PublicKey
        Link net.Conn
    }

    Address: {
        IPv4 string
        Port string
    }

    Package: {
        Info Info
        From From
        To To
        Head Head
        Body Body
    }

    type Info: {
        NET string
    }

    From: {
        Hashname string
        Address string
        Public string
    }

    To: {
        Address string
    }

    Head: {
        Title string
        Mode string
    }

    Body: {
        Data [DATA_SIZE]string
        Desc [DATA_SIZE]string
        Time string
        Hash string
        Sign string
    }
}
=======================================================================
*/
package gopeer
