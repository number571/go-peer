package utils

import (
	error_chain "github.com/g8rswimmer/error-chain"
)

func MergeErrors(pErrors ...error) error {
	errChain := error_chain.New()
	for _, err := range pErrors {
		if err == nil {
			continue
		}
		errChain.Add(err)
	}
	if len(errChain.Errors()) == 0 {
		return nil
	}
	return errChain
}
