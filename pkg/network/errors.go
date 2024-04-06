package network

const (
	errPrefix = "pkg/network = "
)

type SNetworkError struct {
	str string
}

func (err *SNetworkError) Error() string {
	return errPrefix + err.str
}

var (
	ErrNoConnections        = &SNetworkError{"no connections"}
	ErrWriteTimeout         = &SNetworkError{"write timeout"}
	ErrBroadcastMessage     = &SNetworkError{"broadcast message"}
	ErrCreateListener       = &SNetworkError{"create listener"}
	ErrListenerAccept       = &SNetworkError{"listener accept"}
	ErrHasLimitConnections  = &SNetworkError{"has limit connections"}
	ErrConnectionIsExist    = &SNetworkError{"connection already exist"}
	ErrConnectionIsNotExist = &SNetworkError{"connection is not exist"}
	ErrCloseConnection      = &SNetworkError{"close connection"}
	ErrAddConnections       = &SNetworkError{"add connection"}
)
