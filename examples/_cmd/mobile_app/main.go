package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	stateStop = 0
	stateRun  = 1
)

func switchState(state int) int {
	if state == stateStop {
		return stateRun
	}
	return stateStop
}

func textState(state int) string {
	if state == stateStop {
		return "Run TEST"
	}
	return "Stop TEST"
}

func main() {
	a := app.New()
	w := a.NewWindow("app")

	var button *widget.Button

	state := stateStop
	button = widget.NewButton(textState(state), func() {
		state = switchState(state)
		button.SetText(textState(state))
	})

	w.SetContent(container.New(
		layout.NewCenterLayout(),
		container.NewVBox(
			widget.NewLabel("Test Application"),
			button,
		),
	))
	w.ShowAndRun()
}
