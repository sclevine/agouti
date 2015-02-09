package dsl

import (
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/core"
)

var (
	driver  core.WebDriver
	failDSL func(message string, callerSkip ...int)
)

func init() {
	failDSL = ginkgo.Fail
}

// StartPhantomJS starts a PhantomJS WebDriver service for use with CreatePage.
func StartPhantomJS() {
	checkWebDriverNotStarted()
	driver, _ = core.PhantomJS()
	checkFailure(driver.Start())
}

// StartChrome starts a ChromeDriver WebDriver service for use with CreatePage.
func StartChromeDriver() {
	checkWebDriverNotStarted()
	driver = core.ChromeDriver()
	checkFailure(driver.Start())
}

// StartSelenium starts a Selenium WebDriver service for use with CreatePage.
func StartSelenium() {
	checkWebDriverNotStarted()
	driver, _ = core.Selenium()
	checkFailure(driver.Start())
}

// StopWebDriver stops the current running WebDriver.
func StopWebDriver() {
	if driver == nil {
		failDSL("WebDriver not started", 1)
	}
	driver.Stop()
	driver = nil
}

// CreatePage creates a new session using the current running WebDriver.
// For Selenium, the browserName determines which browser to use for the session.
func CreatePage(browserName ...string) core.Page {
	if driver == nil {
		failDSL("WebDriver not started", 1)
	}
	capabilities := core.Use()
	if len(browserName) > 0 {
		capabilities.Browser(browserName[0])
	}
	newPage, err := driver.Page(capabilities)
	checkFailure(err)
	return newPage
}

// CustomPage creates a new session with a custom set of desired capabilities
// using the current running WebDriver. The core.Use() function may be used
// to generate this set of capabilities. For Selenium, the capabilities
// Browser(string) method sets which browser to use for the session.
func CustomPage(capabilities core.Capabilities) core.Page {
	if driver == nil {
		failDSL("WebDriver not started", 1)
	}
	newPage, err := driver.Page(capabilities)
	checkFailure(err)
	return newPage
}

func checkWebDriverNotStarted() {
	if driver != nil {
		failDSL("WebDriver already started", 2)
	}
}

func checkFailure(err error) {
	if err != nil {
		failDSL(fmt.Sprintf("Agouti failure: %s", err), 2)
	}
}
