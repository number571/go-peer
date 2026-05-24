package layer2

import (
	"reflect"
	"sync"
)

var (
	_ IKeysContainer = &sKeysContainer{}
)

type sKeysContainer struct {
	fMtx       *sync.RWMutex
	fMap       map[string]IParticipantKey
	fList      []IParticipantKey
	lockedType reflect.Type
}

func NewKeysContainer(pKeys ...IParticipantKey) IKeysContainer {
	keysContainer := &sKeysContainer{
		fMtx:  &sync.RWMutex{},
		fMap:  make(map[string]IParticipantKey, 128),
		fList: make([]IParticipantKey, 0, 128),
	}
	for _, v := range pKeys {
		keysContainer.Add(v)
	}
	return keysContainer
}

func (p *sKeysContainer) Get(pHash string) (IParticipantKey, bool) {
	p.fMtx.RLock()
	defer p.fMtx.RUnlock()

	v, ok := p.fMap[pHash]
	return v, ok
}

func (p *sKeysContainer) Del(pHash string) bool {
	p.fMtx.Lock()
	defer p.fMtx.Unlock()

	if _, ok := p.fMap[pHash]; !ok {
		return false
	}

	for i, v := range p.fList {
		if pHash != v.GetHasher().ToString() {
			continue
		}
		p.fList = append(p.fList[:i], p.fList[i+1:]...)
		delete(p.fMap, pHash)
		return true
	}

	panic("cipher not found in list (but exist in map)")
}

func (p *sKeysContainer) Add(pKey IParticipantKey) bool {
	p.fMtx.Lock()
	defer p.fMtx.Unlock()

	valType := reflect.TypeOf(pKey)
	if p.lockedType == nil {
		p.lockedType = valType
	}
	if valType != p.lockedType {
		return false
	}

	strKey := pKey.GetHasher().ToString()
	if _, ok := p.fMap[strKey]; ok {
		return false
	}
	p.fMap[strKey] = pKey

	p.fList = append(p.fList, pKey)
	return true
}

func (p *sKeysContainer) List() []IParticipantKey {
	p.fMtx.RLock()
	defer p.fMtx.RUnlock()

	list := make([]IParticipantKey, len(p.fList))
	copy(list, p.fList)
	return list
}
