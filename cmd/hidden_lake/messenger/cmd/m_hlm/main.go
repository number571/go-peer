package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/number571/go-peer/internal/mobile"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
)

var (
	gAppHLM types.ICommand
)

func main() {
	a := app.New()
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
	if gAppHLM == nil {
		var err error
		gAppHLM, err = initApp(mobile.CAndroidFullPath)
		if err != nil {
			return err
		}
	}
	if err := gAppHLM.Run(); err != nil {
		return err
	}
	return nil
}

func destructApp() error {
	if gAppHLM != nil {
		if err := gAppHLM.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func buttonActions(a fyne.App) *widget.Button {
	var (
		button *widget.Button
	)

	if !filesystem.OpenFile(mobile.CAndroidFullPath).IsExist() {
		if err := os.Mkdir(mobile.CAndroidFullPath, 0744); err != nil {
			panic(err)
		}
	}

	state := mobile.NewState(a, pkg_settings.CServiceName).
		WithConstructApp(constructApp).
		WithDestructApp(destructApp)

	button = widget.NewButton(
		state.ToString(),
		func() {
			button.SetText("Processing...")
			state.SwitchOnSuccess(func() error {
				if !state.IsRun() {
					return mobile.OpenURL(a, fmt.Sprintf("http://%s/about", hlmURL))
				}
				return nil
			})
			button.SetText(state.ToString())
		},
	)

	return button
}
