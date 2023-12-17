package state

type IStateFunc func() error

type IState interface {
	Enable(IStateFunc) error
	Disable(IStateFunc) error
}
