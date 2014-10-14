package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/page"
	"fmt"
)

var browser page.Browser

type AgoutiPage interface {
	page.PageOnly
	page.Selection
}

func StartPhantomJS() {
	var err error
	checkBrowser()
	browser, err = page.PhantomJS()
	checkFailure(err)
	checkFailure(browser.Start())
}

func StartChrome() {
	var err error
	checkBrowser()
	browser, err = page.Chrome()
	checkFailure(err)
	checkFailure(browser.Start())
}

func StartSelenium() {
	var err error
	checkBrowser()
	browser, err = page.Selenium()
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
}

func CreatePage(browserName ...string) AgoutiPage {
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
