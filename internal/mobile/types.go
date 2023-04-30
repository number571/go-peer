package mobile

type IMobileState interface {
	IsRun() bool
	ToString() string
	SwitchOnSuccess(func() error)

	WithConstructApp(func() error) IMobileState
	WithDestructApp(func() error) IMobileState
}
