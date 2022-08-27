package utils

type ICloser interface {
	Close() error
}

// returns last error from slice
func CloseAll(cs []ICloser) error {
	var lastErr error
	for _, c := range cs {
		if err := c.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
