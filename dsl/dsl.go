package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/page"
)

func Background(body interface{}, timeout ...float64) bool {
	return ginkgo.BeforeEach(body, timeout...)
}

func Feature(text string, body func()) bool {
	return ginkgo.Describe(text, body)
}

func FFeature(text string, body func()) bool {
	return ginkgo.FDescribe(text, body)
}

func PFeature(text string, body func()) bool {
	return ginkgo.PDescribe(text, body)
}

func XFeature(text string, body func()) bool {
	return ginkgo.XDescribe(text, body)
}

func Scenario(description string, body func(), timeout ...float64) bool {
	return ginkgo.It(description, body, timeout...)
}

func FScenario(description string, body func(), timeout ...float64) bool {
	return ginkgo.FIt(description, body, timeout...)
}

func PScenario(description string, ignored ...interface{}) bool {
	return ginkgo.PIt(description, ignored...)
}

func XScenario(description string, ignored ...interface{}) bool {
	return ginkgo.XIt(description, ignored...)
}

func Step(description string, bodies ...func()) {
	ginkgo.GinkgoWriter.Write([]byte("\n  Step - " + description))
	for _, body := range bodies {
		body()
	}
}

func PStep(description string, _ ...interface{}) {
	ginkgo.GinkgoWriter.Write([]byte("\n  Skipping Step - " + description))
}

func XStep(description string, _ ...interface{}) {
	ginkgo.GinkgoWriter.Write([]byte("\n  Skipping Step - " + description))
}

type Page interface {
	Navigate(url string) error
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error
	DeleteCookieByName(name string) error
	ClearCookies() error
	URL() (string, error)
	Size(height, width int) error
	page.Selection
}

func CreatePage() Page {
	page, err := page.PhantomPage()
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return page
}
