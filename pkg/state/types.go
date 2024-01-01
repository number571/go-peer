package state

type IStateF func() error

type IState interface {
	Enable(IStateF) error
	Disable(IStateF) error
}
