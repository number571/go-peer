package gopeer

import (
	"testing"
)

func TestNewListener(t *testing.T) {
	var errorIsExist bool
	listener := NewListener(settings.IS_CLIENT)
	if listener == nil {
		errorIsExist = true
		t.Errorf("listener == nil [NewListener(IS_CLIENT)]")
	}
	if listener.address.ipv4+listener.address.port != settings.IS_CLIENT {
		errorIsExist = true
		t.Errorf("listener.address.ipv4 + listener.address.port != settings.IS_CLIENT")
	}
	listener = NewListener("")
	if listener != nil {
		errorIsExist = true
		t.Errorf("listener != nil [NewListener('')]")
	}
	address := "ipv4:port"
	listener = NewListener(address)
	if listener == nil {
		errorIsExist = true
		t.Errorf("listener == nil [NewListener(address)]")
	}
	if listener.address.ipv4+listener.address.port != address {
		errorIsExist = true
		t.Errorf("listener.address.ipv4 + listener.address.port != address")
	}
	if !errorIsExist {
		t.Logf("NewListener() success")
	}
}

func TestListenerOpen(t *testing.T) {
	var (
		listener     = new(Listener)
		errorIsExist bool
	)
	defer func() {
		listener.Close()
	}()
	listener = NewListener(settings.IS_CLIENT)
	listener = listener.Open(nil)
	if listener != nil {
		errorIsExist = true
		t.Errorf("listener != nil [listener.Open(nil)]")
	}
	nodeKey, nodeCert := GenerateCertificate(settings.NETWORK, settings.KEY_SIZE)
	listener = NewListener(settings.IS_CLIENT)
	listener = listener.Open(&Certificate{
		Cert: []byte(nodeCert),
		Key:  []byte(nodeKey),
	})
	if listener == nil || listener.listen != nil {
		errorIsExist = true
		t.Errorf("listener == nil || listener.listen != nil [1]")
	}
	listener = NewListener(":7070")
	listener = listener.Open(&Certificate{
		Cert: []byte(nodeCert),
		Key:  []byte(nodeKey),
	})
	if listener == nil || listener.listen == nil {
		errorIsExist = true
		t.Errorf("listener == nil || listener.listen == nil [2]")
	}
	if !errorIsExist {
		t.Logf("listener.Open() success")
	}
}
