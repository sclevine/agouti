package appium

import (
	"fmt"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/internal/target"
)

type mobileSession interface {
	LaunchApp() error
	CloseApp() error
	InstallApp(appPath string) error
	PerformTouch(actions []interface{}) error
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

// Override Find to find by A11y instead of CSS selector
// Note: embedded overrides don't need the same signature - just method name.
func (d *Device) Find(id string) *Selection {
	return d.addSelector(target.A11yID, id)
}

func (d *Device) FindByID(resourceID string) *Selection {
	return d.addSelector(target.A11yID, resourceID) // needs new target type
}

func (d *Device) FindByClass(class string) *Selection {
	return d.addSelector(target.Class, class)
}

func (d *Device) FindByXPath(xPath string) *Selection {
	return d.addSelector(target.XPath, xPath)
}

// Make this behave differently on different devices
// Consider appium.Android() and appium.IOS() just like
// agouti has agouti.PhantomJS() and agouti.Selenium().
func (d *Device) FindByUI(uiQuery string) *Selection {
	return d.addSelector(target.AndroidAut, uiQuery)
}

func (d *Device) addSelector(selectorType target.Type, value string) *Selection {
	selectors := target.Selectors{}.Append(selectorType, value)
	selection := d.WithSelectors(agouti.Selectors(selectors))
	return &Selection{selection, d.session}
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

// Don't return anything from the mobile package. Bring TouchAction up
// to the appium level.
func (d *Device) TouchAction() *TouchAction {
	return NewTouchAction(d.session)
}
