package agouti

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/page"
	"github.com/sclevine/agouti/phantom"
	"github.com/sclevine/agouti/webdriver"
	"time"
	"fmt"
)

const PHANTOM_HOST = "127.0.0.1"
const PHANTOM_PORT = 8910

var phantomService *phantom.Service

type Page page.Page
type Selection page.Selection
type FinalSelection page.FinalSelection

type Do func(Selection)

func (f Do) Call(selection page.Selection) {
	f(selection)
}

type Cookie webdriver.Cookie

func SetupAgouti() bool {
	phantomService = &phantom.Service{Host: PHANTOM_HOST, Port: PHANTOM_PORT, Timeout: 3 * time.Second}
	if err := phantomService.Start(); err != nil {
		panic("Agouti failed to start phantomjs: " + err.Error())
	}
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

func CreatePage() Page {
	session, err := phantomService.CreateSession()
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	driver := &webdriver.Driver{session}

	return Page(page.NewPage(driver, ginkgo.Fail))
}

func CreateCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) webdriver.Cookie {
	return webdriver.Cookie{name, value, path, domain, secure, httpOnly, expiry}
}
