package integration_test

import (
	"bytes"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
)

var (
	phantomDriver    = agouti.PhantomJS()
	chromeDriver     = agouti.ChromeDriver()
	seleniumDriver   = agouti.Selenium(agouti.Browser("firefox"))
	selendroidDriver = agouti.Selendroid("selendroid-standalone-0.15.0-with-dependencies.jar")

	phantomURL, chromeURL, seleniumURL, selendroidURL string

	headlessOnly = os.Getenv("HEADLESS_ONLY") == "true"
	mobile       = os.Getenv("MOBILE") == "true"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	Expect(phantomDriver.Start()).To(Succeed())
	if !headlessOnly {
		Expect(chromeDriver.Start()).To(Succeed())
		Expect(seleniumDriver.Start()).To(Succeed())
	}
	if mobile {
		Expect(selendroidDriver.Start()).To(Succeed())
	}
	phantomURL = phantomDriver.URL()
	chromeURL = chromeDriver.URL()
	seleniumURL = seleniumDriver.URL()
	selendroidURL = seleniumDriver.URL()
	return serialize(phantomURL, chromeURL, seleniumURL, selendroidURL)
}, func(data []byte) {
	urls := deserialize(data)
	phantomURL = urls[0]
	chromeURL = urls[1]
	seleniumURL = urls[2]
	selendroidURL = urls[3]
})

func serialize(urls ...string) []byte {
	var byteURLs [][]byte
	for _, url := range urls {
		byteURLs = append(byteURLs, []byte(url))
	}
	return bytes.Join(byteURLs, []byte{0})
}

func deserialize(urlData []byte) []string {
	byteURLs := bytes.Split(urlData, []byte{0})
	var urls []string
	for _, byteURL := range byteURLs {
		urls = append(urls, string(byteURL))
	}
	return urls
}

var _ = SynchronizedAfterSuite(func() {}, func() {
	Expect(phantomDriver.Stop()).To(Succeed())
	if !headlessOnly {
		Expect(chromeDriver.Stop()).To(Succeed())
		Expect(seleniumDriver.Stop()).To(Succeed())
	}
	if mobile {
		Expect(selendroidDriver.Stop()).To(Succeed())
	}
})
