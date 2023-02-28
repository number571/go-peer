package types

// returns last error from slice
func StopAllCommands(cs []ICommand) error {
	var lastErr error
	for _, c := range cs {
		if err := c.Stop(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
