package core

import "time"

// Chrome returns an instance of a ChromeDriver WebDriver.
//
// DEPRECATED: Use ChromeDriver instead.
func Chrome() (WebDriver, error) {
	return ChromeDriver(), nil
}

// Connect opens a Page using the provided WebDriver URL.
//
// DEPRECATED: Use NewPage instead.
func Connect(desired Capabilities, url string) (Page, error) {
	return NewPage(url, desired)
}

// CustomWebDriver returns an instance of a WebDriver specified by
// a templated URL and command.
//
// DEPRECATED: Use NewWebDriver instead.
func CustomWebDriver(url string, command []string, timeout ...time.Duration) WebDriver {
	return NewWebDriver(url, command, timeout...)
}
