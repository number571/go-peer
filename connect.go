package gopeer

import (
	"crypto/rsa"
)

// Return client's hashname.
func (connect *Connect) Hashname() string {
	return connect.hashname
}

// Return client's public.
func (connect *Connect) Public() *rsa.PublicKey {
	x := *connect.public
	return &x
}

// Return public key of intermediate node.
func (connect *Connect) Throw() *rsa.PublicKey {
	return connect.throwClient
}

// Return client's address.
func (connect *Connect) Address() string {
	return connect.address
}

// Return client's session.
func (connect *Connect) Session() []byte {
	return connect.session
}

// Return client's certificate.
func (connect *Connect) Certificate() []byte {
	return connect.certificate
}
