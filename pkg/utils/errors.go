package utils

import "fmt"

func MergeErrors(errors ...error) error {
	var resErr error
	for _, err := range errors {
		if err == nil {
			continue
		}
		if resErr == nil {
			resErr = err
			continue
		}
		resErr = fmt.Errorf("%w, %w", err, resErr)
	}
	return resErr
}
