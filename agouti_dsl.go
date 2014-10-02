package agouti

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/sclevine/agouti/phantom"
	"time"
)

const PHANTOM_HOST = "127.0.0.1"
const PHANTOM_PORT = 8910

var beforeSuiteBody = func() {}
var afterSuiteBody = func() {}

func init() {
	var phantomService *phantom.Service

	ginkgo.BeforeSuite(func() {
		phantomService = &phantom.Service{Host: PHANTOM_HOST, Port: PHANTOM_PORT, Timeout: time.Second}
		phantomService.Start()
		beforeSuiteBody()
	})

	// TODO: figure out how to test
	ginkgo.AfterSuite(func() {
		afterSuiteBody()
		phantomService.Stop()
	})
}

func BeforeSuite(body func()) bool {
	// TODO: panic if already defined, conform to ginkgo interface
	beforeSuiteBody = body
	return true
}

func AfterSuite(body func()) bool {
	// TODO: panic if already defined, conform to ginkgo interface
	afterSuiteBody = body
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
	body()
}

func Navigate(url string) {

}

type Scopable interface {
	Within(selector string, scopes ...func(Scopable)) Scopable
	ShouldContainText(text string)
}

func Within(selector string, scopes ...func(Scopable)) Scopable {
	return scope{}.Within(selector, scopes...)
}

type scope struct{}

func (s scope) Within(selector string, scopes ...func(Scopable)) Scopable {
	return scope{}
}

func (s scope) ShouldContainText(text string) {
	gomega.Expect(true).To(gomega.BeTrue())
}
