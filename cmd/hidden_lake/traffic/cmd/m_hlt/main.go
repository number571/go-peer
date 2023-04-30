package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/mobile"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/types"
)

var (
	gAppHLT types.ICommand
)

func main() {
	a := app.New()
	w := a.NewWindow("app")
	w.SetContent(container.New(
		layout.NewCenterLayout(),
		container.NewVBox(
			widget.NewLabel("Hidden Lake Traffic"),
			buttonActions(a),
		),
	))

	w.SetOnClosed(func() { destructApp() })
	w.ShowAndRun()
}

func constructApp() error {
	if gAppHLT == nil {
		var err error
		gAppHLT, err = initApp(mobile.CAndroidFullPath)
		if err != nil {
			return err
		}
	}
	if err := gAppHLT.Run(); err != nil {
		return err
	}
	return nil
}

func destructApp() error {
	if gAppHLT != nil {
		if err := gAppHLT.Stop(); err != nil {
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

	state := mobile.NewState(a, settings.CServiceName).
		WithConstructApp(constructApp).
		WithDestructApp(destructApp)

	button = widget.NewButton(
		state.ToString(),
		func() {
			button.SetText("Processing...")
			state.SwitchOnSuccess(func() error {
				if !state.IsRun() {
					return mobile.OpenURL(a, fmt.Sprintf("http://%s/api/index", hltURL))
				}
				return nil
			})
			button.SetText(state.ToString())
		},
	)

	return button
}
