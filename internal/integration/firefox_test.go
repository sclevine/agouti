package integration_test

import (
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
	. "github.com/sclevine/agouti/page"
)

var _ = PFeature("Agouti running on Firefox", func() {
	Scenario("Loading a page", func() {
		StartSelenium()
		defer StopSelenium()

		page := CreatePage("firefox")
		page.Size(640, 480)
		page.Navigate(server.URL)

		Step("find the title of the page", func() {
			Expect(page).To(HaveTitle("Page Title"))
		})
	})
})
