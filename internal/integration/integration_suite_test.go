package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/internal/integration"

	"os"
	"testing"
)

var (
	phantomDriver  WebDriver
	chromeDriver   WebDriver
	seleniumDriver WebDriver
	headlessOnly   bool
)

func TestIntegration(t *testing.T) {
	headlessOnly = os.Getenv("HEADLESS_ONLY") == "true"
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

	Server.Start()
})

var _ = AfterSuite(func() {
	Server.Close()
	phantomDriver.Stop()

	if !headlessOnly {
		chromeDriver.Stop()
		seleniumDriver.Stop()
	}
})
