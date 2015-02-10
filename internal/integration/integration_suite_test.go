package integration_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
)

var (
	phantomDriver  WebDriver
	chromeDriver   WebDriver
	seleniumDriver WebDriver
	headlessOnly   = os.Getenv("HEADLESS_ONLY") == "true"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	phantomDriver, _ = PhantomJS()
	Expect(phantomDriver.Start()).To(Succeed())

	if !headlessOnly {
		seleniumDriver, _ = Selenium()
		chromeDriver = ChromeDriver()
		Expect(seleniumDriver.Start()).To(Succeed())
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
