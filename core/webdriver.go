package core

import (
	"errors"
	"fmt"

	"github.com/sclevine/agouti/api"
)

// WebDriver controls a Selenium, PhantomJS, or ChromeDriver process.
type WebDriver interface {
	// Start launches the WebDriver process.
	Start() error

	// Stop ends all remaining sessions and stops the WebDriver process.
	Stop() error

	// Page returns a new WebDriver session. The optional config argument
	// configures the returned page. For instance:
	//    driver.Page(Use().Without("javascriptEnabled"))
	// For Selenium, this argument must include a browser. For instance:
	//    seleniumDriver.Page(Use().Browser("safari"))
	Page(config ...Capabilities) (Page, error)
}

type webDriver struct {
	apiWebDriver interface {
		Open(desired map[string]interface{}) (*api.Session, error)
		Start() error
		Stop() error
	}
}

func (d *webDriver) Page(desired ...Capabilities) (Page, error) {
	if len(desired) == 0 {
		desired = append(desired, capabilities{})
	} else if len(desired) > 1 {
		return nil, errors.New("too many arguments")
	}

	session, err := d.apiWebDriver.Open(desired[0].(capabilities))
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %s", err)
	}
	newPage := newPage(session)
	return newPage, nil
}

func (d *webDriver) Start() error {
	if err := d.apiWebDriver.Start(); err != nil {
		return err
	}
	return nil
}

func (d *webDriver) Stop() error {
	return d.apiWebDriver.Stop()
}
