// Package core is a WebDriver API for Go.
package core

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/service"
	"github.com/sclevine/agouti/core/internal/session"
)

// Chrome returns an instance of a ChromeDriver WebDriver.
func Chrome() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s", address)
	command := []string{"chromedriver", "--silent", "--port=" + port}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &driver{service: service}, nil
}

// PhantomJS returns an instance of a PhantomJS WebDriver.
func PhantomJS() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	url := fmt.Sprintf("http://%s", address)
	command := []string{"phantomjs", fmt.Sprintf("--webdriver=%s", address)}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &driver{service: service}, nil
}

// Selenium returns an instance of a Selenium WebDriver.
func Selenium() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s/wd/hub", address)
	command := []string{"selenium-server", "-port", port}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &driver{service: service}, nil
}

func freeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
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
