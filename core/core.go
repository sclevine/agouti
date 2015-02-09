// Package core is a WebDriver API for Go.
package core

import (
	"fmt"
	"time"

	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/internal/service"
	"github.com/sclevine/agouti/api/internal/session"
)

// ChromeDriver returns an instance of a ChromeDriver WebDriver.
func ChromeDriver() WebDriver {
	return CustomWebDriver("http://{{.Address}}", []string{"chromedriver", "--silent", "--port={{.Port}}"})
}

// PhantomJS returns an instance of a PhantomJS WebDriver.
// The return error is deprecated and will always be nil.
func PhantomJS() (WebDriver, error) {
	return CustomWebDriver("http://{{.Address}}", []string{"phantomjs", "--webdriver={{.Address}}"}), nil
}

// Selenium returns an instance of a Selenium WebDriver.
// The return error is deprecated and will always be nil.
func Selenium() (WebDriver, error) {
	return CustomWebDriver("http://{{.Address}}/wd/hub", []string{"selenium-server", "-port", "{{.Port}}"}), nil
}

// CustomWebDriver returns an instance of a WebDriver specified by
// a templated URL and command. The URL should be the location of the
// WebDriver Wire Protocol web service brought up by the command. The
// command should be provided as a list of arguments (which are each
// templated). The optional timeout specifies how long to wait for the
// web service to become available. Default timeout is 5 seconds.
//
// Valid template parameters are:
//   {{.Host}} - local address to bind to (usually 127.0.0.1)
//   {{.Port}} - arbitrary free port on the local address
//   {{.Address}} - {{.Host}}:{{.Port}}
//
// Selenium JAR example:
//   command := []string{"java", "-jar", "selenium-server.jar", "-port", "{{.Port}}"}
//   core.CustomWebDriver("http://{{.Address}}/wd/hub", command)
func CustomWebDriver(url string, command []string, timeout ...time.Duration) WebDriver {
	if len(timeout) == 0 {
		timeout = []time.Duration{5 * time.Second}
	}

	driverService := &service.Service{
		URLTemplate: url,
		CmdTemplate: command,
		Timeout:     timeout[0],
	}

	return &driver{service: driverService}
}

// Connect opens a session using the provided WebDriver URL and returns a Page.
func Connect(capabilities Capabilities, url string) (Page, error) {
	pageSession, err := session.Open(url, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to open WebDriver session: %s", err)
	}

	client := &api.Client{Session: pageSession}
	return newPage(client), nil
}

// SauceLabs opens a Sauce Labs session and returns a Page. Does not support Sauce Connect.
func SauceLabs(name, platform, browser, version, username, key string) (Page, error) {
	url := "http://ondemand.saucelabs.com/wd/hub"
	capabilities := capabilities{
		"name":        name,
		"platform":    platform,
		"browserName": browser,
		"version":     version,
		"username":    username,
		"accessKey":   key,
	}

	pageSession, err := session.Open(url, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to Sauce Labs: %s", err)
	}

	client := &api.Client{Session: pageSession}
	return newPage(client), nil
}
