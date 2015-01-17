package core

import (
	"errors"
	"fmt"

	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/session"
)

// WebDriver controls a Selenium, PhantomJS, or ChromeDriver process.
type WebDriver interface {
	// Start launches the WebDriver process.
	Start() error

	// Stop ends all remaining sessions and stops the WebDriver process.
	Stop()

	// Page returns a new WebDriver session. The optional config argument
	// configures the returned page. For instance:
	//    driver.Page(Use().Without("javascriptEnabled"))
	// For Selenium, this argument must include a browser. For instance:
	//    seleniumDriver.Page(Use().Browser("safari"))
	Page(config ...Capabilities) (Page, error)
}

type driver struct {
	service driverService
	pages   []Page
}

type driverService interface {
	URL() (string, error)
	Start() error
	Stop()
}

func (d *driver) Page(config ...Capabilities) (Page, error) {
	if len(config) == 0 {
		config = append(config, capabilities{})
	} else if len(config) > 1 {
		return nil, errors.New("too many arguments")
	}

	url, err := d.service.URL()
	if err != nil {
		return nil, fmt.Errorf("WebDriver not started: %s", err)
	}

	pageSession, err := session.Open(url, config[0])
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %s", err)
	}

	client := &api.Client{Session: pageSession}
	newPage := newPage(client)
	d.pages = append(d.pages, newPage)
	return newPage, nil
}

func (d *driver) Start() error {
	if err := d.service.Start(); err != nil {
		return fmt.Errorf("failed to start service: %s", err)
	}

	return nil
}

func (d *driver) Stop() {
	for _, openPage := range d.pages {
		openPage.Destroy()
	}

	d.service.Stop()
	return
}
