package selenium_test

import (
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Feature("Selenium", func() {
	Scenario("Firefox", func() {
		page := CreatePage("firefox")
		page.Size(640, 480)

		Step("navigates to example.com", func() {
			Expect(page.Navigate("http://example.com")).To(Succeed())
		})

		Step("finds the page title", func() {
			Expect(page).To(HaveTitle("Example Domain"))
		})

		Step("finds the header text", func() {
			Expect(page.Find("h1")).To(HaveText("Example Domain"))
		})
	})

	Scenario("Safari", func() {
		page := CreatePage("safari")
		page.Size(640, 480)

		Step("navigates to example.com", func() {
			Expect(page.Navigate("http://example.com")).To(Succeed())
		})

		Step("finds the page title", func() {
			Expect(page).To(HaveTitle("Example Domain"))
		})

		Step("finds the header text", func() {
			Expect(page.Find("h1")).To(HaveText("Example Domain"))
		})
	})
})
