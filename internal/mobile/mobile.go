package mobile

import (
	"fyne.io/fyne/v2"
)

var (
	_ IMobileState = &sMobileState{}
)

type sMobileState struct {
	fApp          fyne.App
	fIsRun        bool
	fServiceName  string
	fConstructApp func() error
	fDestructApp  func() error
}

func NewState(pApp fyne.App, pServiceName string) IMobileState {
	return &sMobileState{
		fApp:         pApp,
		fIsRun:       false,
		fServiceName: pServiceName,
	}
}

func (p *sMobileState) WithConstructApp(pConstructApp func() error) IMobileState {
	p.fConstructApp = pConstructApp
	return p
}

func (p *sMobileState) WithDestructApp(pDestructApp func() error) IMobileState {
	p.fDestructApp = pDestructApp
	return p
}

func (p *sMobileState) IsRun() bool {
	return p.fIsRun
}

func (p *sMobileState) ToString() string {
	if !p.fIsRun {
		return "Start " + p.fServiceName
	}
	return "Stop " + p.fServiceName
}

func (p *sMobileState) SwitchOnSuccess(pDo func() error) {
	if err := p.doDefault(); err != nil {
		p.sendNotification(err)
		return
	}
	if err := pDo(); err != nil {
		p.sendNotification(err)
		return
	}
	p.fIsRun = !p.fIsRun
}

func (p *sMobileState) sendNotification(pErr error) {
	if p.fIsRun {
		p.fApp.SendNotification(fyne.NewNotification("runApp error", pErr.Error()))
		return
	}
	p.fApp.SendNotification(fyne.NewNotification("stopApp error", pErr.Error()))
}

func (p *sMobileState) doDefault() error {
	switch !p.fIsRun {
	case true:
		if p.fConstructApp != nil {
			if err := p.fConstructApp(); err != nil {
				return err
			}
		}
	case false:
		if p.fDestructApp != nil {
			if err := p.fDestructApp(); err != nil {
				return err
			}
		}
	}
	return nil
}
