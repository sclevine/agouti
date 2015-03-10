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

// Agouti wrapped selectors

func (d *Device) Find(selector string) *Selection {
	return d.wrapSelection(d.Page.Find(selector))
}

func (d *Device) FindByXPath(xPath string) *Selection {
	return d.wrapSelection(d.Page.FindByXPath(xPath))
}

func (d *Device) FindByLink(text string) *Selection {
	return d.wrapSelection(d.Page.FindByLink(text))
}

func (d *Device) All(selector string) *MultiSelection {
	return d.wrapMultiSelection(d.Page.All(selector))
}

// Appium-specific selectors

// Finds by Accessibility ID, under Android and iOS
func (d *Device) FindByA11yID(id string) *Selection {
	return d.newSelection(d.addSelector(target.A11yID, id).Single())
}

// Finds by class, Appium searches native view classes with this method.
func (d *Device) FindByClass(class string) *Selection {
	return d.newSelection(d.addSelector(target.Class, class).Single())
}

func (d *Device) FindByAndroidUI(uiautomatorQuery string) *Selection {
	return d.newSelection(d.addSelector(target.AndroidAut, uiautomatorQuery).Single())
}

func (d *Device) FindByiOSUI(uiautomationQuery string) *Selection {
	return d.newSelection(d.addSelector(target.IOSAut, uiautomationQuery).Single())
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

// Selection helpers
func (d *Device) addSelector(selectorType target.Type, value string) target.Selectors {
	return target.Selectors(target.Selectors{}.Append(selectorType, value))
}
func (d *Device) newSelection(selectors target.Selectors) *Selection {
	return &Selection{d.WithSelectors(agouti.Selectors(selectors)), d.session}
}
func (d *Device) newMultiSelection(selectors target.Selectors) *MultiSelection {
	return &MultiSelection{*d.newSelection(selectors), d.newSelection}
}
func (d *Device) wrapSelection(selection *agouti.Selection) *Selection {
	return &Selection{selection, d.session}
}
func (d *Device) wrapMultiSelection(selection *agouti.MultiSelection) *MultiSelection {
	return &MultiSelection{*d.wrapSelection(&selection.Selection), d.newSelection}
}
