package interrupt

import (
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"
)

// Close all elements in a slice.
func CloseAll(pClosers []types.ICloser) error {
	var err error
	for _, c := range pClosers {
		err = errors.AppendError(err, c.Close())
	}
	return err
}
