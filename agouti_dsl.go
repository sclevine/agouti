package agouti

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/phantom"
	"github.com/sclevine/agouti/webdriver"
	"time"
	"github.com/sclevine/agouti/page"
)

const PHANTOM_HOST = "127.0.0.1"
const PHANTOM_PORT = 8910

var phantomService *phantom.Service

type Page interface {
	Selection
}

type Selection interface {
	FinalSelection
	Within(selector string, bodies ...func(Selection)) Selection
}

type FinalSelection interface {
	ShouldContainText(text string)
}

type Cookie struct {
	Name string
	Value interface{}
	Path string
	Domain string
	Secure bool
	HTTPOnly bool
	Expiry int64
}

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
func Navigate(url string, cookies []Cookie) Page {
	session, err := phantomService.CreateSession()
	if err != nil {
		ginkgo.Fail(err.Error()) // TODO: test error
	}

	driver := &webdriver.Driver{session}

	for _, cookie := range cookies {
		if err := driver.SetCookie(&webdriver.Cookie(cookie)); err != nil {
			ginkgo.Fail(err.Error())
		}
	}

	if err := driver.Navigate(url); err != nil {
		ginkgo.Fail(err.Error())// TODO: test error
	}

	return &page.Page{driver}
}
