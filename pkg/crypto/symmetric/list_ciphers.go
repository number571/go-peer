package symmetric

import (
	"sync"
)

var (
	_ IListCiphers = &sListCiphers{}
)

type sListCiphers struct {
	fMtx  *sync.RWMutex
	fMap  map[string]interface{}
	fList []ICipher
}

func NewListCiphers(pCiphers ...ICipher) IListCiphers {
	listCiphers := &sListCiphers{
		fMtx:  &sync.RWMutex{},
		fMap:  make(map[string]interface{}, 128),
		fList: make([]ICipher, 0, 128),
	}
	for _, v := range pCiphers {
		listCiphers.Add(v)
	}
	return listCiphers
}

func (p *sListCiphers) Add(pCipher ICipher) bool {
	p.fMtx.Lock()
	defer p.fMtx.Unlock()

	strKey := pCipher.ToString()
	if _, ok := p.fMap[strKey]; ok {
		return false
	}
	p.fMap[strKey] = struct{}{}

	p.fList = append(p.fList, pCipher)
	return true
}

func (p *sListCiphers) Del(pCipher ICipher) bool {
	p.fMtx.Lock()
	defer p.fMtx.Unlock()

	strKey := pCipher.ToString()
	if _, ok := p.fMap[strKey]; !ok {
		return false
	}

	for i, v := range p.fList {
		if strKey != v.ToString() {
			continue
		}
		p.fList = append(p.fList[:i], p.fList[i+1:]...)
		delete(p.fMap, strKey)
		return true
	}

	panic("cipher not found in list (but exist in map)")
}

func (p *sListCiphers) Get() []ICipher {
	p.fMtx.RLock()
	defer p.fMtx.RUnlock()

	list := make([]ICipher, len(p.fList))
	copy(list, p.fList)
	return list
}
