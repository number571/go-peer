package interrupt

import (
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

// Stop all elements in a slice.
func StopAll(pCommands []types.ICommand) error {
	errList := make([]error, 0, len(pCommands))
	for _, c := range pCommands {
		if err := c.Stop(); err != nil {
			errList = append(errList, err)
		}
	}
	return utils.MergeErrors(errList...)
}
