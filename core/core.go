// Agouti core is a general-purpose WebDriver API for Golang
package core

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/service"
	"github.com/sclevine/agouti/core/internal/session"
	"github.com/sclevine/agouti/core/internal/types"
	"github.com/sclevine/agouti/core/internal/webdriver"
)

type Selection types.Selection
type MultiSelection types.MultiSelection
type Page types.Page

// WebDriver represents a Selenium, PhantomJS, or ChromeDriver process
type WebDriver interface {
	// Start launches the WebDriver process
	Start() error

	// Stop ends all remaining sessions and stops the WebDriver process
	Stop()

	// Page returns a new WebDriver session.
	// For Selenium, browserName is the type of browser ("firefox", "safari", "chrome", etc.)
	Page(browserName ...string) (types.Page, error)
}

// Chrome returns an instance of a ChromeDriver WebDriver
func Chrome() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s", address)
	command := []string{"chromedriver", "--silent", "--port=" + port}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &webdriver.Driver{Service: service}, nil
}

// PhantomJS returns an instance of a PhantomJS WebDriver
func PhantomJS() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	url := fmt.Sprintf("http://%s", address)
	command := []string{"phantomjs", fmt.Sprintf("--webdriver=%s", address)}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &webdriver.Driver{Service: service}, nil
}

// Selenium returns an instance of a Selenium WebDriver
func Selenium() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s/wd/hub", address)
	command := []string{"selenium-server", "-port", port}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &webdriver.Driver{Service: service}, nil
}

// SauceLabs returns a Page with a Sauce Labs session
func SauceLabs(name, platform, browser, version, username, key string) (Page, error) {
	url := "http://ondemand.saucelabs.com/wd/hub"
	capabilities := map[string]interface{}{
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
	return &page.Page{Client: client}, nil
}

func freeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}
