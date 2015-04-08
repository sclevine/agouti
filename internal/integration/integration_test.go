package integration_test

import . "github.com/onsi/ginkgo"

var _ = Describe("integration tests", func() {
	testPage("PhantomJS", &phantomURL)
	testSelection("PhantomJS", &phantomURL)

	if !headlessOnly {
		testPage("ChromeDriver", &chromeURL)
		testSelection("ChromeDriver", &chromeURL)
		testPage("Firefox", &seleniumURL)
		testSelection("Firefox", &seleniumURL)
	}

	if mobile {
		testMobile("Android", &selendroidURL)
	}
})
