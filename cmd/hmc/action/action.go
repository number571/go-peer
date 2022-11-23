package action

var (
	_ IAction = &sAction{}
)

type sAction struct {
	fDescription string
	fHandler     func()
}

func NewAction(description string, handler func()) IAction {
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
