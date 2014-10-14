package chrome_test

import (
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Feature("ChromeDriver", func() {
	Scenario("Chrome", func() {
		page := CreatePage()
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
