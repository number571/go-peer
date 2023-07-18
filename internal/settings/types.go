package settings

import "github.com/number571/go-peer/pkg/client/message"

type IConfigSettings interface {
	// IsValid() bool
	message.ISettings
}
