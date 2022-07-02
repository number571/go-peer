package utils

// returns last error from slice
func CloseAll(cs []ICloser) error {
	var err error
	for _, c := range cs {
		e := c.Close()
		if e != nil {
			err = e
		}
	}
	return err
}
