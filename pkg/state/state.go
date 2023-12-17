package state

import (
	"errors"
	"fmt"
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

func (p *sState) Enable(f IStateFunc) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fEnabled {
		return errors.New("state already enabled")
	}

	if f != nil {
		if err := f(); err != nil {
			return fmt.Errorf("enable state error: %w", err)
		}
	}

	p.fEnabled = true
	return nil
}

func (p *sState) Disable(f IStateFunc) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fEnabled {
		return errors.New("state already disabled")
	}

	if f != nil {
		if err := f(); err != nil {
			return fmt.Errorf("disable state error: %w", err)
		}
	}

	p.fEnabled = false
	return nil
}
