// Agouti core is a general-purpose WebDriver API for Golang
package core

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/service"
	"github.com/sclevine/agouti/core/internal/session"
	"github.com/sclevine/agouti/core/internal/webdriver"
)

type Selection interface {
	Find(selector string) *selection.Selection
	FindByXPath(selector string) *selection.Selection
	FindByLink(text string) *selection.Selection
	FindByLabel(text string) *selection.Selection
	All(selector string) *selection.MultiSelection
	AllByXPath(selector string) *selection.MultiSelection
	AllByLink(text string) *selection.MultiSelection
	AllByLabel(text string) *selection.MultiSelection
	String() string
	Count() (int, error)
	Click() error
	DoubleClick() error
	Fill(text string) error
	Text() (string, error)
	Attribute(attribute string) (string, error)
	CSS(property string) (string, error)
	Check() error
	Uncheck() error
	Selected() (bool, error)
	Visible() (bool, error)
	Enabled() (bool, error)
	Select(text string) error
	Submit() error
	EqualsElement(comparable interface{}) (bool, error)
}

type MultiSelection interface {
	Selection
	At(index int) Selection
	Single() Selection
}

type Page interface {
	Destroy() error
	Navigate(url string) error
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error
	DeleteCookie(name string) error
	ClearCookies() error
	URL() (string, error)
	Size(width, height int) error
	Screenshot(filename string) error
	Title() (string, error)
	HTML() (string, error)
	RunScript(body string, arguments map[string]interface{}, result interface{}) error
	PopupText() (string, error)
	EnterPopupText(text string) error
	ConfirmPopup() error
	CancelPopup() error
	Forward() error
	Back() error
	Refresh() error
	Find(selector string) *selection.Selection
	FindByXPath(selector string) *selection.Selection
	FindByLink(text string) *selection.Selection
	FindByLabel(text string) *selection.Selection
	All(selector string) *selection.MultiSelection
	AllByXPath(selector string) *selection.MultiSelection
	AllByLink(text string) *selection.MultiSelection
	AllByLabel(text string) *selection.MultiSelection
}

type Capabilities interface {
	Browser(browser string) session.Capabilities
	Version(version string) session.Capabilities
	Platform(platform string) session.Capabilities
	With(feature string) session.Capabilities
	Without(feature string) session.Capabilities
	JSON() string
}

// WebDriver represents a Selenium, PhantomJS, or ChromeDriver process
type WebDriver interface {
	// Start launches the WebDriver process
	Start() error

	// Stop ends all remaining sessions and stops the WebDriver process
	Stop()

	// Page returns a new WebDriver session. The optional capabilities
	// argument allows for specification of the desired browser capabilities.
	// For Selenium, the capabilities argument should specify a browser.
	Page(capabilities ...session.JSONable) (*page.Page, error)
}

// Use returns a Capabilities object with the following chainable methods:
//  Browser(string) - {chrome|firefox|safari|iphone|...} - browser name
//  Version(string) - ex. "3.6" - browser version
//  Platform(string) - {WINDOWS|XP|VISTA|MAC|LINUX|UNIX} - browser platform
//  With(feature string) - enable the specified feature
//  Without(feature string) - disable the specified feature
// The Map() method returns a map[string]interface{}
func Use() Capabilities {
	return session.Capabilities{}
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
	capabilities := session.Capabilities{
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

// Connect returns a Page with a generic session open on the provided URL
func Connect(capabilities Capabilities, url string) (Page, error) {
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
