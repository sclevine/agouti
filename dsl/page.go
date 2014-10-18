package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/core"
)

var driver core.WebDriver

// StartPhantomJS starts a PhantomJS WebDriver service for use with CreatePage.
func StartPhantomJS() {
	var err error
	checkWebDriver()
	driver, err = core.PhantomJS()
	checkFailure(err)
	checkFailure(driver.Start())
}

// StartChrome starts a ChromeDriver WebDriver service for use with CreatePage.
func StartChrome() {
	var err error
	checkWebDriver()
	driver, err = core.Chrome()
	checkFailure(err)
	checkFailure(driver.Start())
}

// StartSelenium starts a Selenium WebDriver service for use with CreatePage.
func StartSelenium() {
	var err error
	checkWebDriver()
	driver, err = core.Selenium()
	checkFailure(err)
	checkFailure(driver.Start())
}

// StopWebdriver stops the current running WebDriver.
func StopWebdriver() {
	if driver == nil {
		ginkgo.Fail("WebDriver not started", 1)
	}
	driver.Stop()
	driver = nil
}

// CreatePage creates a new session using the current running WebDriver.
// For Selenium, the browserName argument determines which driver to start the session in.
func CreatePage(browserName ...string) core.Page {
	newPage, err := driver.Page(browserName...)
	checkFailure(err)
	return newPage
}

func checkWebDriver() {
	if driver != nil {
		ginkgo.Fail("WebDriver already started", 2)
	}
}

func checkFailure(err error) {
	if err != nil {
		ginkgo.Fail(err.Error(), 2)
	}
}
