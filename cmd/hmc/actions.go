package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/number571/go-peer/cmd/hmc/action"
	"github.com/number571/go-peer/cmd/hmc/settings"
	"github.com/number571/go-peer/cmd/hms/hmc"
	"github.com/number571/go-peer/modules/inputter"
	"github.com/number571/go-peer/modules/payload"
)

func newActions() action.IActions {
	return action.IActions{
		"": action.NewAction(
			"",
			func() {},
		),
		"exit": action.NewAction(
			"exit from application",
			func() { syscall.Kill(os.Getpid(), syscall.SIGINT) },
		),
		"help": action.NewAction(
			"get help information about this application",
			helpAction,
		),
		"whoami": action.NewAction(
			"get information about authorized user",
			whoamiAction,
		),
		"size": action.NewAction(
			"get count of emails and pages in user database",
			sizeAction,
		),
		"list": action.NewAction(
			"get list of emails by page (10 emails in 1 page)",
			listAction,
		),
		"read": action.NewAction(
			"get full information about 1 email by ID",
			readAction,
		),
		"push": action.NewAction(
			"create new email and push this to servers",
			pushAction,
		),
		"load": action.NewAction(
			"try load all emails to authorized user from servers",
			loadAction,
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

func loadAction() {
	// check connections
	for _, addr := range gWrapper.Config().Original().Connections() {
		client := hmc.NewClient(
			hmc.NewBuilder(gClient),
			hmc.NewRequester(addr),
		)

		// connect to server
		size, err := client.Size()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// infinite loop protection
		if size > settings.CReceiveSize {
			fmt.Println("size of messages > limit")
			continue
		}

		// load new emails
		for i := uint64(0); i < size; i++ {
			msg, err := client.Load(i)
			if err != nil {
				break
			}

			pubKey, pld, err := gClient.Decrypt(msg)
			if err != nil {
				continue
			}

			if gWrapper.Config().Original().F2F() {
				if _, ok := gWrapper.Config().GetNameByPubKey(pubKey); !ok {
					continue
				}
			}

			if pld.Head() != settings.CHeadPayload {
				continue
			}

			if len(strings.Split(string(pld.Body()), settings.CSeparator)) < 2 {
				continue
			}

			err = gDB.Push(gClient.PubKey().Address().Bytes(), msg)
			if err != nil {
				continue
			}
		}
	}
}

func whoamiAction() {
	fmt.Printf("Address:\n%s;\nPublic Key:\n%s;\n",
		gClient.PubKey().Address().String(),
		gClient.PubKey())
}

func sizeAction() {
	size, err := gDB.Size(gClient.PubKey().Address().Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}

	pages := int(math.Ceil(float64(size) / float64(settings.CCountInPage)))

	fmt.Printf("Count: %d\n", size)
	fmt.Printf("Pages: %d\n", pages)
}

func listAction() {
	page := inputter.NewInputter("Page: ").String()
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println(err)
		return
	}

	start := (pageInt - 1) * settings.CCountInPage
	if start < 0 {
		fmt.Println("Page can't be <= 0")
		return
	}

	for pointer := start; pointer < start+settings.CCountInPage; pointer++ {
		msg, err := gDB.Load(gClient.PubKey().Address().Bytes(), uint64(pointer))
		if err != nil {
			break
		}

		pubKey, pld, err := gClient.Decrypt(msg)
		if err != nil {
			panic(err)
		}

		title := []rune(strings.Split(string(pld.Body()), settings.CSeparator)[0])
		if len(title) > settings.CListLenTitle {
			title = []rune(string(title[:settings.CListLenTitle-3]) + "...")
		}

		name, _ := gWrapper.Config().GetNameByPubKey(pubKey)
		fmt.Printf("%s\nID: %d;\nFrom: %s [%s];\nTitle: %s;\n",
			settings.CSeparator,
			pointer,
			pubKey.Address().String(),
			name,
			string(title),
		)
	}
}

func readAction() {
	id := inputter.NewInputter("ID: ").String()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg, err := gDB.Load(gClient.PubKey().Address().Bytes(), uint64(idInt))
	if err != nil {
		return
	}

	pubKey, pld, err := gClient.Decrypt(msg)
	if err != nil {
		panic(err)
	}

	from, _ := gWrapper.Config().GetNameByPubKey(pubKey)
	splited := strings.Split(string(pld.Body()), settings.CSeparator)

	fmt.Printf("%s\nFROM:\n%s [%s]\n%s\nTITLE:\n%s\n%s\nMESSAGE:\n%s\n%s\nPUBLIC_KEY:\n%s\n",
		settings.CSeparator,
		pubKey.Address().String(),
		from,
		settings.CSeparator,
		splited[0],
		settings.CSeparator,
		strings.Join(splited[1:], settings.CSeparator),
		settings.CSeparator,
		pubKey.String(),
	)
}

func pushAction() {
	name := inputter.NewInputter("Receiver: ").String()
	pubKey, ok := gWrapper.Config().GetPubKeyByName(name)
	if !ok {
		fmt.Println("Receiver's public key undefined")
		return
	}

	title := inputter.NewInputter("Title: ").String()
	if title == "" {
		fmt.Println("Title is nil")
		return
	}

	msg := inputter.NewInputter("Message: ").String()
	if msg == "" {
		fmt.Println("Message is nil")
		return
	}

	withoutError := 0
	pushReq := hmc.NewBuilder(gClient).Push(
		pubKey,
		payload.NewPayload(
			settings.CHeadPayload,
			[]byte(fmt.Sprintf("%s%s%s", title, settings.CSeparator, msg)),
		),
	)
	for _, addr := range gWrapper.Config().Original().Connections() {
		err := hmc.NewRequester(addr).Push(pushReq)
		if err != nil {
			fmt.Printf("%s - '%s'\n", addr, err.Error())
			continue
		}
		withoutError++
	}

	fmt.Printf("Message successfully sent (%d servers are received)\n", withoutError)
}
