package integration_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
)

var (
	phantomDriver  *agouti.WebDriver
	chromeDriver   *agouti.WebDriver
	seleniumDriver *agouti.WebDriver
	headlessOnly   = os.Getenv("HEADLESS_ONLY") == "true"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	phantomDriver = agouti.PhantomJS()
	Expect(phantomDriver.Start()).To(Succeed())

	if !headlessOnly {
		seleniumDriver = agouti.Selenium()
		Expect(seleniumDriver.Start()).To(Succeed())

		chromeDriver = agouti.ChromeDriver()
		Expect(chromeDriver.Start()).To(Succeed())
	}
})

var _ = AfterSuite(func() {
	Expect(phantomDriver.Stop()).To(Succeed())

	if !headlessOnly {
		Expect(chromeDriver.Stop()).To(Succeed())
		Expect(seleniumDriver.Stop()).To(Succeed())
	}
})
