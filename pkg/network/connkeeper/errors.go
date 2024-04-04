package connkeeper

const (
	errPrefix = "pkg/network/connkeeper = "
)

type SConnKeeperError struct {
	str string
}

func (err *SConnKeeperError) Error() string { return errPrefix + err.str }

var (
	ErrRunning = &SConnKeeperError{"connkeeper running"}
)
