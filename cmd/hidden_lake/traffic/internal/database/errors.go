package database

type SIsExistError struct{}

func (p *SIsExistError) Error() string {
	return "message is already exist"
}
