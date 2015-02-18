package dsl

import "github.com/sclevine/agouti"

var driver *agouti.WebDriver

// StartPhantomJS starts a PhantomJS WebDriver service for use with CreatePage.
func StartPhantomJS() {
	checkWebDriverNotStarted()
	driver = agouti.PhantomJS()
	checkFailure(driver.Start())
}

// StartChrome starts a ChromeDriver WebDriver service for use with CreatePage.
func StartChromeDriver() {
	checkWebDriverNotStarted()
	driver = agouti.ChromeDriver()
	checkFailure(driver.Start())
}

// StartSelenium starts a Selenium WebDriver service for use with CreatePage.
func StartSelenium() {
	checkWebDriverNotStarted()
	driver = agouti.Selenium()
	checkFailure(driver.Start())
}

// StopWebDriver stops the current running WebDriver.
func StopWebDriver() {
	if driver == nil {
		globalFailHandler("WebDriver not started", 1)
	}
	driver.Stop()
	driver = nil
}

// CreatePage creates a new session using the current running WebDriver.
// For Selenium, the browserName determines which browser to use for the session.
func CreatePage(browserName ...string) *agouti.Page {
	if driver == nil {
		globalFailHandler("WebDriver not started", 1)
	}
	capabilities := agouti.NewCapabilities()
	if len(browserName) > 0 {
		capabilities.Browser(browserName[0])
	}
	newPage, err := driver.NewPage(agouti.Desired(capabilities))
	checkFailure(err)
	return newPage
}

// CustomPage creates a new session with a custom set of desired capabilities
// using the current running WebDriver. The agouti.Use() function may be used
// to generate this set of capabilities. For Selenium, the capabilities
// Browser(string) method sets which browser to use for the session.
func CustomPage(capabilities agouti.Capabilities) *agouti.Page {
	if driver == nil {
		globalFailHandler("WebDriver not started", 1)
	}
	newPage, err := driver.NewPage(agouti.Desired(capabilities))
	checkFailure(err)
	return newPage
}

func checkWebDriverNotStarted() {
	if driver != nil {
		globalFailHandler("WebDriver already started", 2)
	}
}
