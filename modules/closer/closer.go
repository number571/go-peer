package closer

import "github.com/number571/go-peer/modules"

// returns last error from slice
func CloseAll(cs []modules.ICloser) error {
	var lastErr error
	for _, c := range cs {
		if err := c.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
