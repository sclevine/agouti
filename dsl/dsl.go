package dsl

import (
	"github.com/onsi/ginkgo"
	"github.com/sclevine/agouti/page"
)

func Feature(text string, body func()) bool {
	return ginkgo.Describe(text, body)
}

func Background(body interface{}, timeout ...float64) bool {
	return ginkgo.BeforeEach(body, timeout...)
}

func Scenario(description string, body func()) bool {
	return ginkgo.It(description, body)
}

func Step(description string, bodies ...func()) {
	ginkgo.GinkgoWriter.Write([]byte("\n  Step - " + description))
	for _, body := range bodies {
		body()
	}
}

type Page interface {
	Navigate(url string) error
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error
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
