// Agouti dsl implements a Capybara-like DSL for writing acceptance tests
// It is provided entirely for convenience and not necessary for writing tests using Agouti core and matchers
package dsl

import "github.com/onsi/ginkgo"

// Background is equivalent to Ginkgo BeforeEach
func Background(body interface{}, timeout ...float64) bool {
	return ginkgo.BeforeEach(body, timeout...)
}

// Feature is equivalent to Ginkgo Describe
func Feature(text string, body func()) bool {
	return ginkgo.Describe(text, body)
}

// FFeature is equilavent to a Ginkgo FDescribe (Focused Describe)
func FFeature(text string, body func()) bool {
	return ginkgo.FDescribe(text, body)
}

// PFeature is equilavent to a Ginkgo PDescribe (Pending Describe)
func PFeature(text string, body func()) bool {
	return ginkgo.PDescribe(text, body)
}

// XFeature is equilavent to a Ginkgo XDescribe (Pending Describe)
func XFeature(text string, body func()) bool {
	return ginkgo.XDescribe(text, body)
}

// Scenario is equivalent to a Ginkgo It
func Scenario(description string, body func(), timeout ...float64) bool {
	return ginkgo.It(description, body, timeout...)
}

// FScenario is equivalent to a Ginkgo FIt (Focused It)
func FScenario(description string, body func(), timeout ...float64) bool {
	return ginkgo.FIt(description, body, timeout...)
}

// PScenario is equivalent to a Ginkgo PIt (Pending It)
func PScenario(description string, ignored ...interface{}) bool {
	return ginkgo.PIt(description, ignored...)
}

// XScenario is equivalent to a Ginkgo XIt (Pending It)
func XScenario(description string, ignored ...interface{}) bool {
	return ginkgo.XIt(description, ignored...)
}

// Step is equivalent to Ginkgo By
func Step(text string, callbacks ...func()) {
	ginkgo.By(text, callbacks...)
}
