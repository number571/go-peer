package queue

const (
	errPrefix = "pkg/network/anonymity/queue = "
)

type SQueueError struct {
	str string
}

func (err *SQueueError) Error() string { return errPrefix + err.str }

var (
	ErrRunning    = &SQueueError{"queue running"}
	ErrQueueLimit = &SQueueError{"queue limit"}
)
