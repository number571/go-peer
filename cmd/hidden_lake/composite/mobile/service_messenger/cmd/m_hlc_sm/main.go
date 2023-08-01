package main

import (
	"fmt"
	"os"
	"strings"

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
	gPath string
	gApp  types.ICommand
)

func main() {
	a := app.New()
	w := a.NewWindow("app")

	rootPath := a.Storage().RootURI().Path()
	gPath = strings.Split(rootPath, "/fyne/")[0] + "/hidden_lake"

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

	if !filesystem.OpenFile(gPath).IsExist() {
		if err := os.Mkdir(gPath, 0744); err != nil {
			panic(err)
		}
	}

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
