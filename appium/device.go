package appium

import (
	"fmt"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api/mobile"
)

type Device struct {
	*agouti.Page
	session *mobile.Session
}

func (d *Device) DeviceMethod(thing string) error {
	if err := d.session.SetEndpoint(thing); err != nil {
		return fmt.Errorf("failed to do stuff: %s", err)
	}
	return nil
}

// override finder methods
func (d *Device) Find(selector string) *Selection {
	return &Selection{d.Page.Find(selector), d.session}
}

func (d *Device) FindByID(resourceId string) *Selection {
	return &Selection{&agouti.Selection{}, d.session}
}

func newDevice(session *mobile.Session, page *agouti.Page) *Device {
	return &Device{
		Page:       page,
		session:    session,
	}
}

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

func (d *Device) TouchAction() *mobile.TouchAction {
	return &mobile.TouchAction{
		Session: d.session,
	}
}

func (d *Device) PerformMultiTouch(actions ...*mobile.TouchAction) {
}

func (s *Device) FindID(selector string) *agouti.Selection {
	return nil
}
func (s *Device) FindByXPath() *agouti.Selection {
	return nil
}
func (s *Device) FindA11yID(id string) *agouti.Selection {
	return nil
}
func (s *Device) FindByClass(class string) *agouti.Selection {
	return nil
}
func (s *Device) FindiOS(uiautomationQuery string) *agouti.Selection {
	return nil
}
func (s *Device) FindAndroid(uiautomatorQuery string) *agouti.Selection {
	return nil
}
