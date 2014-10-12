package dsl

import "github.com/onsi/ginkgo"

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

func Step(text string, callbacks ...func()) {
	ginkgo.By(text, callbacks...)
}

