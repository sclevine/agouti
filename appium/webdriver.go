package appium

import (
	"fmt"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api/mobile"
)

type WebDriver struct {
	driver *agouti.WebDriver
}

func New(options ...Option) *WebDriver {
	capabilities := config{}.merge(options).desired
	agoutiWebDriver := agouti.NewWebDriver("http://{{.Address}}/wd/hub", []string{"appium", "-p", "{{.Port}}"}, agouti.Desired(capabilities))
	return &WebDriver{agoutiWebDriver}
}

func (w *WebDriver) NewDevice(options ...Option) (*Device, error) {
	capabilities := config{}.merge(options).desired
	page, err := w.driver.NewPage(agouti.Desired(capabilities))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebDriver: %s", err)
	}
	mobileSession := &mobile.Session{page.Session()}

	return newDevice(mobileSession, page), nil
}

func (w *WebDriver) Start() error {
	return w.driver.Start()
}

func (w *WebDriver) Stop() error {
	return w.driver.Stop()
}

func (w *WebDriver) Debug(state bool) {
	//w.driver.Debug(state)
}
