package dsl_test

import (
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/dsl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DSL sanity checks", func() {
	Feature("starting and stopping WebDrivers", func() {
		Background(func() {
			// TODO: test background
		})

		AfterEach(func() {
			StopWebDriver()
		})

		Scenario("PhantomJS", func() {
			StartPhantomJS()
			page := CreatePage()
			Destroy(page)
		})

		Scenario("ChromeDriver", func() {
			StartChromeDriver()
			page := CreatePage()
			Destroy(page)
		})

		Scenario("Selenium", func() {
			StartSelenium()
			page := CreatePage("firefox")
			Destroy(page)

			Step("using CustomPage", func() {
				page = CustomPage(Use().Browser("chrome"))
				Destroy(page)
			})
		})
	})

	XFeature("this Describe is pending (using X)", func() {
		Scenario("so this does not run", func() {
			Fail("failed to pend spec")
		})
	})

	PFeature("this Describe is pending (using P)", func() {
		Scenario("so this does not run", func() {
			Fail("failed to pend spec")
		})
	})

	XScenario("this is pending (using X) and does not run", func() {
		Fail("failed to pend spec")
	})

	PScenario("this is pending (using P) and does not run", func() {
		Fail("failed to pend spec")
	})

})
