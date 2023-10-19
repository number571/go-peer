package types

import "github.com/number571/go-peer/pkg/errors"

// Close all elements in a slice.
func CloseAll(pClosers []ICloser) error {
	var err error
	for _, c := range pClosers {
		err = errors.AppendError(err, c.Close())
	}
	return err
}
