package mobile

import (
	"net/url"

	"fyne.io/fyne/v2"
)

func OpenURL(pApp fyne.App, pPageURL string) error {
	u, err := url.Parse(pPageURL)
	if err != nil {
		return err
	}
	return pApp.OpenURL(u)
}
