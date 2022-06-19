package errors

// used to determine error
const (
	CErrorClient uint64 = iota + 1
	CErrorServer

	CErrorEncode
	CErrorDecode

	CErrorEncrypt
	CErrorDecrypt

	CErrorWrite
	CErrorRead
	CErrorExecute

	CErrorNotEqual
	CErrorUndefined

	CErrorCritical
	CErrorWarning
)
