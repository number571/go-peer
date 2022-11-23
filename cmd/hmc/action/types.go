package action

type IActions map[string]IAction

type IAction interface {
	Description() string
	Do()
}
