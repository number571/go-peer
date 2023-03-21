package types

// returns last error from slice
func StopAll(pCommands []ICommand) error {
	var lastErr error
	for _, c := range pCommands {
		if err := c.Stop(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
