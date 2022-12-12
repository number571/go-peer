package closer

import "github.com/number571/go-peer/pkg/types"

// returns last error from slice
func CloseAll(cs []types.ICloser) error {
	var lastErr error
	for _, c := range cs {
		if err := c.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
