// Package agouti is a universal WebDriver client for Go.
// It extends the agouti/api package to provide a feature-rich interface for
// controlling a Web Browser.
package agouti

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

// NewPage opens a Page using the provided WebDriver URL.
func NewPage(url string, options ...Option) (*Page, error) {
	desiredCapabilities := getOptions(options).desired
	session, err := api.Open(url, desiredCapabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebDriver: %s", err)
	}
	return newPage(session), nil
}

// SauceLabs opens a Sauce Labs session and returns a *Page. Does not support Sauce Connect.
func SauceLabs(name, platform, browser, version, username, accessKey string) (*Page, error) {
	url := fmt.Sprintf("http://%s:%s@ondemand.saucelabs.com/wd/hub", username, accessKey)
	capabilities := Capabilities{
		"name":        name,
		"platform":    platform,
		"browserName": browser,
		"version":     version,
	}

	session, err := api.Open(url, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to Sauce Labs: %s", err)
	}
	return newPage(session), nil
}
