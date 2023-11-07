package interrupt

import (
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"
)

// Stop all elements in a slice.
func StopAll(pCommands []types.ICommand) error {
	var err error
	for _, c := range pCommands {
		err = errors.AppendError(err, c.Stop())
	}
	return err
}
