package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/number571/go-peer/internal/mobile"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
)

var (
	gPath string
	gApp  types.ICommand
)

func main() {
	a := app.NewWithID("hidden_lake")
	gPath = a.Storage().RootURI().Path()

	w := a.NewWindow("app")
	w.SetContent(container.New(
		layout.NewCenterLayout(),
		container.NewVBox(
			widget.NewLabel("Hidden Lake Messenger"),
			buttonActions(a),
		),
	))

	w.SetOnClosed(func() { destructApp() })
	w.ShowAndRun()
}

func constructApp() error {
	if gApp == nil {
		var err error
		gApp, err = initApp(gPath)
		if err != nil {
			return err
		}
	}
	return gApp.Run()
}

func destructApp() error {
	if gApp == nil {
		return nil
	}
	return gApp.Stop()
}

func buttonActions(a fyne.App) *widget.Button {
	var (
		button *widget.Button
	)

	state := mobile.NewMobileState(a, pkg_settings.CServiceName).
		WithConstructApp(constructApp).
		WithDestructApp(destructApp)

	button = widget.NewButton(
		state.ToString(),
		func() {
			button.SetText("Processing...")
			state.SwitchOnSuccess(func() error {
				if !state.IsRun() {
					urlPage := fmt.Sprintf("http://%s/about", pkg_settings.CDefaultInterfaceAddress)
					return mobile.OpenURL(a, urlPage)
				}
				return nil
			})
			button.SetText(state.ToString())
		},
	)

	return button
}
