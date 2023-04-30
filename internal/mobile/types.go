package mobile

type IMobileState interface {
	IsRun() bool
	ToString() string
	SwitchOnSuccess(func(bool) error)

	WithConstructApp(func() error) IMobileState
	WithDestructApp(func() error) IMobileState
}
