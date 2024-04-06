package utils

const (
	errPrefix = "cmd/hidden_lake/applications/filesharer/internal/utils = "
)

type SUtilsError struct {
	str string
}

func (err *SUtilsError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRespSizeGeLimit = &SUtilsError{"response size >= limit message size"}
	ErrGetSettingsHLS  = &SUtilsError{"get settings from HLS"}
	ErrGetSizeInBase64 = &SUtilsError{"get size in base64"}
)
