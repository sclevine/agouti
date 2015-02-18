package core

import "encoding/json"

// A Capabilities instance defines the desired capabilities the WebDriver
// should use to configure a Page.
//
// For example, to open a Firefox page with JavaScript disabled:
//    driver.Page(Use().Browser("firefox").Without("javascriptEnabled"))
// See: https://code.google.com/p/selenium/wiki/DesiredCapabilities
type Capabilities interface {
	// Browser sets the desired browser name - {chrome|firefox|safari|iphone|...}.
	Browser(browser string) Capabilities

	// Version sets the desired browser version (ex. "3.6").
	Version(version string) Capabilities

	// Platform sets the desired browser platform - {WINDOWS|XP|VISTA|MAC|LINUX|UNIX}.
	Platform(platform string) Capabilities

	// With enables the provided feature (ex. "handlesAlerts").
	With(feature string) Capabilities

	// Without disables the provided feature (ex. "javascriptEnabled").
	Without(feature string) Capabilities

	// Custom sets a custom desired capability.
	Custom(key string, value interface{}) Capabilities

	// JSON returns a JSON string representing the desired capabilities.
	JSON() (string, error)
}

// Use returns a Capabilities instance that can be passed to a page.
// All methods called on this instance will modify the original instance.
func Use() Capabilities {
	return capabilities{}
}

type capabilities map[string]interface{}

func (c capabilities) Browser(browser string) Capabilities {
	c["browserName"] = browser
	return c
}

func (c capabilities) Version(version string) Capabilities {
	c["version"] = version
	return c
}

func (c capabilities) Platform(platform string) Capabilities {
	c["platform"] = platform
	return c
}

func (c capabilities) With(feature string) Capabilities {
	c[feature] = true
	return c
}

func (c capabilities) Without(feature string) Capabilities {
	c[feature] = false
	return c
}

func (c capabilities) Custom(key string, value interface{}) Capabilities {
	c[key] = value
	return c
}

func (c capabilities) JSON() (string, error) {
	capabilitiesJSON, err := json.Marshal(c)
	return string(capabilitiesJSON), err
}
