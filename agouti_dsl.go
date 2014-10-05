package agouti

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/page"
	"github.com/sclevine/agouti/phantom"
	"github.com/sclevine/agouti/webdriver"
	"time"
)

const PHANTOM_HOST = "127.0.0.1"
const PHANTOM_PORT = 8910

var phantomService *phantom.Service

type Page interface {
	page.Selection
}

type Selection page.Selection
type FinalSelection page.FinalSelection

type Do func(Selection)
type DoFinal func(FinalSelection)

func (f Do) Call(selection page.Selection) {
	f(selection)
}
func (f DoFinal) Call(selection page.FinalSelection) {
	f(selection)
}

type Cookie webdriver.Cookie

func SetupAgouti() bool {
	phantomService = &phantom.Service{Host: PHANTOM_HOST, Port: PHANTOM_PORT, Timeout: time.Second}
	phantomService.Start()
	return true
}

func CleanupAgouti(ignored bool) bool {
	phantomService.Stop()
	return true
}

func Feature(text string, body func()) bool {
	return ginkgo.Describe(text, body)
}

func Background(body interface{}, timeout ...float64) bool {
	return ginkgo.BeforeEach(body, timeout...)
}

func Scenario(description string, body func()) bool {
	return ginkgo.It(description, body)
}

func Step(description string, body func()) {
	ginkgo.GinkgoWriter.Write([]byte("\n  Step - " + description))
	body()
}

// TODO: strip cookies out and test
func Navigate(url string, cookies ...[]Cookie) Page {
	session, err := phantomService.CreateSession()
	if err != nil {
		ginkgo.Fail(err.Error()) // TODO: test error
	}

	driver := &webdriver.Driver{session}

	if len(cookies) > 0 {
		for _, cookie := range cookies[0] {
			driverCookie := webdriver.Cookie(cookie)
			if err := driver.SetCookie(&driverCookie); err != nil {
				ginkgo.Fail(err.Error())
			}
		}
	}

	if err := driver.Navigate(url); err != nil {
		ginkgo.Fail(err.Error()) // TODO: test error
	}

	return &page.Page{driver, ginkgo.Fail}
}
