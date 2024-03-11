package common

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
)

func GetMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   "",
		FWorkSizeBits: 22,
	})
}
