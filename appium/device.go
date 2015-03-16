package appium

import (
	"fmt"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api/mobile"
	"github.com/sclevine/agouti/internal/element"
)

type mobileSession interface {
	element.Client
	LaunchApp() error
	CloseApp() error
	InstallApp(appPath string) error
	PerformTouch(actions []mobile.Action) error
}

type Device struct {
	*agouti.Page
	session mobileSession
}

func newDevice(session mobileSession, page *agouti.Page) *Device {
	return &Device{
		Page:    page,
		session: session,
	}
}

// Device methods

func (d *Device) LaunchApp() error {
	if err := d.session.LaunchApp(); err != nil {
		return fmt.Errorf("failed to launch app: %s", err)
	}
	return nil
}

func (d *Device) CloseApp() error {
	if err := d.session.CloseApp(); err != nil {
		return fmt.Errorf("failed to close app: %s", err)
	}
	return nil
}

func (d *Device) InstallApp(appPath string) error {
	if err := d.session.InstallApp(appPath); err != nil {
		return fmt.Errorf("failed to install app: %s", err)
	}
	return nil
}

func (d *Device) TouchAction() *TouchAction {
	return NewTouchAction(d.session)
}
