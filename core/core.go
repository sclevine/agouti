// Agouti core is a general-purpose WebDriver API for Golang
package core

import (
	"fmt"
	"github.com/sclevine/agouti/core/internal/service"
	"github.com/sclevine/agouti/core/internal/types"
	"github.com/sclevine/agouti/core/internal/webdriver"
	"net"
	"strings"
	"time"
)

type Selection types.Selection
type Page types.Page

// WebDriver represents a Selenium, PhantomJS, or Chrome (via ChromeDriver) WebDriver process
type WebDriver interface {
	// Start launches the WebDriver process
	Start() error

	// Stop ends all remaining sessions and stops the WebDriver process
	Stop()

	// Page returns a new WebDriver session.
	// For Selenium, browserName is the type of browser ("firefox", "safari", "chrome", etc.)
	Page(browserName ...string) (types.Page, error)
}

// Chrome returns an instance of a Chrome WebDriver via ChromeDriver
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

func freeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}
