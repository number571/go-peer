package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/internal/mobile"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/types"
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
	w.ShowAndRun()
}

func buttonActions(a fyne.App) *widget.Button {
	var (
		button *widget.Button
		appHLS types.ICommand
		appHLM types.ICommand
	)

	if !filesystem.OpenFile(mobile.CAndroidFullPath).IsExist() {
		if err := os.Mkdir(mobile.CAndroidFullPath, 0644); err != nil {
			panic(err)
		}
	}

	state := mobile.CStateStop
	button = widget.NewButton(
		mobile.ButtonTextFromState(state, settings.CServiceName),
		func() {
			state = mobile.SwitchState(state)
			button.SetText("Processing...")

			switch state {
			case mobile.CStateRun:
				var err error

				appHLS, appHLM, err = initApp(mobile.CAndroidFullPath)
				if err != nil {
					button.SetText(err.Error())
					return
				}

				if err := appHLS.Run(); err != nil {
					button.SetText(err.Error())
					return
				}
				if err := appHLM.Run(); err != nil {
					button.SetText(err.Error())
					return
				}
			case mobile.CStateStop:
				if appHLS != nil {
					if err := appHLS.Stop(); err != nil {
						button.SetText(err.Error())
						return
					}
				}
				if appHLM != nil {
					if err := appHLM.Stop(); err != nil {
						button.SetText(err.Error())
						return
					}
				}
			}

			button.SetText(mobile.ButtonTextFromState(state, settings.CServiceName))
		},
	)

	return button
}
