package database

type SIsExistError struct{}
type SIsNotExistError struct{}

func (p *SIsExistError) Error() string {
	return "message is already exist"
}

func (p *SIsNotExistError) Error() string {
	return "message is not exist"
}
