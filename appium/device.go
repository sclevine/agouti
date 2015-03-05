package appium

import (
	"fmt"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api/mobile"
)

type Device struct {
	*agouti.Page
	session deviceSession
}

type deviceSession interface {
	SetEndpoint(thing string) error
	// Appium-related, see python-client:appium/webdriver/webdriver.py
	InstallApp(string) error
	RemoveApp(string) error
	IsAppInstalled(string) (bool, error)
	LaunchApp() error
	CloseApp() error
	GetAppStrings(string) ([]string, error)
	GetCurrentActivity() (string, error)
	Lock() error
	Shake() error
	Reset() error
	OpenNotifications() error
	GetSettings() (map[string]interface{}, error)
	UpdateSettings(map[string]interface{}) error
	ToggleLocationServices() error
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

func newDevice(session *mobile.Session, page *agouti.Page) *Device {
	return &Device{
		Page:    page,
		session: session,
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
