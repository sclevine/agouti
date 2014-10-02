package agouti

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const PHANTOM_HOST = "127.0.0.1"
const PHANTOM_PORT = 8910

type Scopable interface {
	Within(selector string, scopes ...func(Scopable)) Scopable
	ShouldContainText(text string)
}

type desiredCapabilities map[string]interface{}

type scope struct{}

func (s scope) Within(selector string, scopes ...func(Scopable)) Scopable {
	return scope{}
}

func (s scope) ShouldContainText(text string) {
	gomega.Expect(false).To(gomega.BeTrue())
}

func Scenario(description string, url string, body func()) {
	ginkgo.It(description, body)
}

func Step(description string, body func()) {
	body()
}

func Within(selector string, scopes ...func(Scopable)) Scopable {
	return scope{}.Within(selector, scopes...)
}
