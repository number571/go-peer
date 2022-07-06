package main

var (
	_ iAction = &sAction{}
)

type sAction struct {
	fDescription string
	fHandler     func()
}

func newAction(description string, handler func()) iAction {
	return &sAction{
		fDescription: description,
		fHandler:     handler,
	}
}

func (act *sAction) Description() string {
	return act.fDescription
}

func (act *sAction) Do() {
	act.fHandler()
}
