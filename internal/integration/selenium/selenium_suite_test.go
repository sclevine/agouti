package selenium_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"

	"os"
	"testing"
)

func TestSelenium(t *testing.T) {
	RegisterFailHandler(Fail)
	if os.Getenv("HEADLESS_ONLY") != "true" {
		RunSpecs(t, "Selenium Suite")
	}
}

var _ = BeforeSuite(func() {
	StartSelenium()
})

var _ = AfterSuite(func() {
	StopWebdriver()
})
