package selenium_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"

	"testing"
)

func TestSelenium(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Selenium Suite")
}

var _ = BeforeSuite(func() {
	StartSelenium()
})

var _ = AfterSuite(func() {
	StopWebdriver()
})
