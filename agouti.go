// Package agouti is a universal WebDriver client for Go.
// It extends the agouti/api package to provide a feature-rich interface for
// controlling a Web Browser.
package agouti

import "fmt"

// PhantomJS returns an instance of a PhantomJS WebDriver.
func PhantomJS() *WebDriver {
	return NewWebDriver("http://{{.Address}}", []string{"phantomjs", "--webdriver={{.Address}}"})
}

// ChromeDriver returns an instance of a ChromeDriver WebDriver.
func ChromeDriver() *WebDriver {
	return NewWebDriver("http://{{.Address}}", []string{"chromedriver", "--silent", "--port={{.Port}}"})
}

// Selenium returns an instance of a Selenium WebDriver.
func Selenium() *WebDriver {
	return NewWebDriver("http://{{.Address}}/wd/hub", []string{"selenium-server", "-port", "{{.Port}}"})
}

// SauceLabs opens a Sauce Labs session and returns a *Page. Does not support Sauce Connect.
func SauceLabs(name, platform, browser, version, username, accessKey string) (*Page, error) {
	url := fmt.Sprintf("http://%s:%s@ondemand.saucelabs.com/wd/hub", username, accessKey)
	capabilities := NewCapabilities().Browser(name).Platform(platform).Version(version)
	capabilities["name"] = name

	return NewPage(url, Desired(capabilities))
}
