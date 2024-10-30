package state

import (
	"errors"
	"sync"
)

var (
	_ IState = &sState{}
)

type sState struct {
	fEnabled bool
	fMutex   sync.Mutex
}

func NewBoolState() IState {
	return &sState{}
}

func (p *sState) Enable(f IStateF) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fEnabled {
		return ErrStateEnabled
	}

	if f != nil {
		if err := f(); err != nil {
			return errors.Join(ErrFuncEnable, err)
		}
	}

	p.fEnabled = true
	return nil
}

func (p *sState) Disable(f IStateF) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fEnabled {
		return ErrStateDisabled
	}

	if f != nil {
		if err := f(); err != nil {
			return errors.Join(ErrFuncDisable, err)
		}
	}

	p.fEnabled = false
	return nil
}
