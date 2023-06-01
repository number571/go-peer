package types

import "github.com/number571/go-peer/pkg/errors"

// returns last error from slice
func StopAll(pCommands []ICommand) error {
	var err error
	for _, c := range pCommands {
		err = errors.AppendError(err, c.Stop())
	}
	return err
}
