package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/core"
)

var browser core.Browser

// StartPhantomJS starts a PhantomJS WebDriver service for use with CreatePage.
func StartPhantomJS() {
	var err error
	checkBrowser()
	browser, err = core.PhantomJS()
	checkFailure(err)
	checkFailure(browser.Start())
}

// StartChrome starts a ChromeDriver WebDriver service for use with CreatePage.
func StartChrome() {
	var err error
	checkBrowser()
	browser, err = core.Chrome()
	checkFailure(err)
	checkFailure(browser.Start())
}

// StartSelenium starts a Selenium WebDriver service for use with CreatePage.
func StartSelenium() {
	var err error
	checkBrowser()
	browser, err = core.Selenium()
	checkFailure(err)
	checkFailure(browser.Start())
}

// StopWebdriver stops the current running WebDriver.
func StopWebdriver() {
	if browser == nil {
		ginkgo.Fail("browser not started", 1)
	}
	browser.Stop()
	browser = nil
}

// CreatePage creates a new session using the current running WebDriver.
// For Selenium, the browserName argument determines which browser to start the session in.
func CreatePage(browserName ...string) core.Page {
	newPage, err := browser.Page(browserName...)
	checkFailure(err)
	return newPage
}

func checkBrowser() {
	if browser != nil {
		ginkgo.Fail("browser already started", 2)
	}
}

func checkFailure(err error) {
	if err != nil {
		ginkgo.Fail(err.Error(), 2)
	}
}
