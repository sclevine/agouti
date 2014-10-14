package dsl

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/core"
)

var browser core.Browser

func StartPhantomJS() {
	var err error
	checkBrowser()
	browser, err = core.PhantomJS()
	checkFailure(err)
	checkFailure(browser.Start())
}

func StartChrome() {
	var err error
	checkBrowser()
	browser, err = core.Chrome()
	checkFailure(err)
	checkFailure(browser.Start())
}

func StartSelenium() {
	var err error
	checkBrowser()
	browser, err = core.Selenium()
	checkFailure(err)
	checkFailure(browser.Start())
}

func StopWebdriver() {
	if browser == nil {
		ginkgo.Fail("browser not started", 1)
	}
	if err := browser.Stop(); err != nil {
		fmt.Println(err)
	}
	browser = nil
}

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
