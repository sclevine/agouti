package agouti

import (
	"fmt"
	"github.com/sclevine/agouti/api"
)

// A WebDriver controls a Selenium, PhantomJS, or ChromeDriver process.
// This struct embeds api.WebDriver, which provides Start and Stop methods
// for starting and stopping the process.
type WebDriver struct {
	*api.WebDriver
}

// NewWebDriver returns an instance of a WebDriver specified by
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
//   agouti.NewWebDriver("http://{{.Address}}/wd/hub", command)
func NewWebDriver(url string, command []string, options ...Option) *WebDriver {
	apiWebDriver := api.NewWebDriver(url, command)
	defaultConfig := &config{timeout: apiWebDriver.Timeout}
	apiWebDriver.Timeout = defaultConfig.apply(options).timeout
	return &WebDriver{apiWebDriver}
}

// ChromeDriver returns an instance of a ChromeDriver WebDriver.
func ChromeDriver() *WebDriver {
	return NewWebDriver("http://{{.Address}}", []string{"chromedriver", "--silent", "--port={{.Port}}"})
}

// PhantomJS returns an instance of a PhantomJS WebDriver.
// The return error is deprecated and will always be nil.
func PhantomJS() *WebDriver {
	return NewWebDriver("http://{{.Address}}", []string{"phantomjs", "--webdriver={{.Address}}"})
}

// Selenium returns an instance of a Selenium WebDriver.
// The return error is deprecated and will always be nil.
func Selenium() *WebDriver {
	return NewWebDriver("http://{{.Address}}/wd/hub", []string{"selenium-server", "-port", "{{.Port}}"})
}

// NewPage returns a *Page that corresponds to a new WebDriver session.
// Any provided options configure the page. For instance:
//    driver.NewPage(agouti.Desired(agouti.NewCapabilities().Without("javascriptEnabled")))
// For Selenium, this argument must include a browser. For instance:
//    seleniumDriver.NewPage(agouti.Desired(agouti.NewCapabilities().Browser("safari")))
func (w *WebDriver) NewPage(options ...Option) (*Page, error) {
	desiredCapabilities := getOptions(options).desired

	session, err := w.Open(desiredCapabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %s", err)
	}

	return newPage(session), nil
}
