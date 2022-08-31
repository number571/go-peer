package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/number571/go-peer/cmd/hlm/settings"
	"github.com/number571/go-peer/modules/action"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/inputter"

	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func newActions() action.IActions {
	return action.IActions{
		"": action.NewAction(
			"",
			func() {},
		),
		"/exit": action.NewAction(
			"exit from application",
			func() { syscall.Kill(os.Getpid(), syscall.SIGINT) },
		),
		"/help": action.NewAction(
			"get help information about this application",
			helpAction,
		),
		"/whoami": action.NewAction(
			"get information about authorized user",
			whoamiAction,
		),
		"/friends": action.NewAction(
			"get list of friends",
			friendsAction,
		),
		"/channel": action.NewAction(
			"change channel with friend",
			channelAction,
		),
	}
}

func helpAction() {
	type sActionWithKey struct {
		fKey    string
		fAction action.IAction
	}

	actions := []*sActionWithKey{}
	for key, act := range gActions {
		actions = append(actions, &sActionWithKey{
			fKey:    key,
			fAction: act,
		})
	}

	sort.SliceStable(actions, func(i, j int) bool {
		return strings.Compare(actions[i].fKey, actions[j].fKey) < 0
	})

	for _, act := range actions {
		switch act.fKey {
		case "":
			continue
		default:
			fmt.Printf("%s:\t%s\n", act.fKey, act.fAction.Description())
		}
	}
}

func whoamiAction() {
	pubKey, err := gClient.PubKey()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Address:\n%s;\nPublic key:\n%s;\n",
		pubKey.Address().String(),
		pubKey.String())
}

func friendsAction() {
	myPubKey, err := gClient.PubKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	friends, err := gClient.Friends()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(friends) == 0 {
		fmt.Println("List is nil")
		return
	}

	i := 1
	for _, pubKey := range friends {
		if pubKey.Address().String() == myPubKey.Address().String() {
			continue
		}
		fmt.Printf("%d. %s\n", i, pubKey.String())
		i++
	}
}

func channelAction() {
	pubKeyStr := inputter.NewInputter("Receiver (Public key): ").String()
	pubKey := asymmetric.LoadRSAPubKey(pubKeyStr)
	if pubKey == nil {
		fmt.Println("Public key is invalid")
		return
	}

	friends, err := gClient.Friends()
	if err != nil {
		fmt.Println(err)
		return
	}

	isExist := false
	for _, friend := range friends {
		if pubKey.Address().String() == friend.Address().String() {
			isExist = true
			break
		}
	}
	if !isExist {
		fmt.Println("Receiver's public key not in list of friends")
		return
	}

	gChannelPubKey = pubKey
	fmt.Println("Success set channel's public key")
}

func sendActionDefault(msg string) {
	if msg == "" {
		return
	}

	if gChannelPubKey == nil {
		fmt.Println("Public key of channel is not set")
		return
	}

	res, err := gClient.Request(
		gChannelPubKey,
		hls_network.NewRequest(
			"POST",
			settings.CTitlePattern,
			settings.CHandlePush,
		).WithBody([]byte(msg)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	var resp hls_settings.SResponse
	if err := json.Unmarshal(res, &resp); err != nil {
		fmt.Println(err)
		return
	}

	if resp.FReturn != hls_settings.CErrorNone {
		fmt.Println(resp.FResult)
		return
	}
}
