package mobile

import (
	"net/url"

	"fyne.io/fyne/v2"
	"github.com/number571/go-peer/pkg/errors"
)

func OpenURL(pApp fyne.App, pPageURL string) error {
	u, err := url.Parse(pPageURL)
	if err != nil {
		return errors.WrapError(err, "parse url")
	}
	if err := pApp.OpenURL(u); err != nil {
		return errors.WrapError(err, "open url")
	}
	return nil
}
