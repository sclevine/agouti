package integration_test

import (
	"time"

	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/internal/integration"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Feature("Agouti running on PhantomJS", func() {
	Scenario("Loading a page", func() {
		page := CreatePage()
		page.Size(640, 480)
		page.Navigate(Server.URL)

		Step("finds the title of the page", func() {
			Expect(page).To(HaveTitle("Page Title"))
		})

		Step("finds a header in the page", func() {
			Expect(page.Find("header")).To(BeFound())
		})

		Step("finds text in the header", func() {
			Expect(page.Find("header")).To(HaveText("Title"))
		})

		Step("matches text in the header", func() {
			Expect(page.Find("header")).To(MatchText("T.+e"))
		})

		Step("finds an element by label text", func() {
			Expect(page.FindByLabel("Some Label")).To(HaveAttribute("value", "some labeled value"))
		})

		Step("finds an element embedded in a label", func() {
			Expect(page.FindByLabel("Some Container Label")).To(HaveAttribute("value", "some embedded value"))
		})

		Step("asserts that text is not in the header", func() {
			Expect(page.Find("header")).NotTo(HaveText("Not-Title"))
		})

		Step("asserts on the visibility of elements", func() {
			Expect(page.Find("header h1")).To(BeVisible())
			Expect(page.Find("header h2")).NotTo(BeVisible())
		})

		Step("allows referring to an element by selection index", func() {
			Expect(page.Find("option").At(0)).To(HaveText("first option"))
			Expect(page.Find("select").At(1).Find("option").At(0)).To(HaveText("third option"))
		})

		Step("allows tests to be scoped by chaining", func() {
			Expect(page.Find("header").Find("h1")).To(HaveText("Title"))
		})

		Step("allows locating elements by XPath", func() {
			Expect(page.Find("header").FindXPath("//h1")).To(HaveText("Title"))
		})

		Step("allows assertions that wait for matchers to be true", func() {
			Expect(page.Find("#some_element")).NotTo(HaveText("some text"))
			Eventually(page.Find("#some_element"), 4*time.Second).Should(HaveText("some text"))
			Consistently(page.Find("#some_element")).Should(HaveText("some text"))
		})

		Step("allows serializing the current page HTML", func() {
			Expect(page.HTML()).To(ContainSubstring(`<div id="some_element" class="some-element" style="color: blue;">some text</div>`))
		})

		Step("allows entering values into fields", func() {
			Fill(page.Find("#some_input"), "some other value")
		})

		Step("allows retrieving attributes by name", func() {
			Expect(page.Find("#some_input")).To(HaveAttribute("value", "some other value"))
		})

		Step("allows asserting on whether a CSS style exists", func() {
			Expect(page.Find("#some_element")).To(HaveCSS("color", "rgba(0, 0, 255, 1)"))
			Expect(page.Find("#some_element")).To(HaveCSS("color", "rgb(0, 0, 255)"))
			Expect(page.Find("#some_element")).To(HaveCSS("color", "blue"))
		})

		Step("allows double-clicking on an element", func() {
			selection := page.Find("#double_click")
			DoubleClick(selection)
			Expect(selection).To(HaveText("double-click success"))
		})

		Step("allows checking a checkbox", func() {
			checkbox := page.Find("#some_checkbox")
			Check(checkbox)
			Expect(checkbox).To(BeSelected())
		})

		Step("allows selecting an option by text", func() {
			selection := page.Find("#some_select")
			Select(selection, "second option")
			Expect(selection.Find("option:last-child")).To(BeSelected())
		})

		Step("allows executing arbitrary javascript", func() {
			arguments := map[string]interface{}{"elementID": "some_element"}
			var result string
			Expect(page.RunScript("return document.getElementById(elementID).innerHTML;", arguments, &result)).To(Succeed())
			Expect(result).To(Equal("some text"))
		})

		Step("allows comparing two selections for equality", func() {
			Expect(page.Find("#some_element")).To(EqualElement(page.FindXPath("//div[@class='some-element']")))
		})

		Step("allows submitting a form", func() {
			Submit(page.Find("#some_form"))
			Eventually(Submitted).Should(BeTrue())
		})

		Step("allows clicking on a link", func() {
			Click(page.Find("a"))
			Expect(page.URL()).To(ContainSubstring("#new_page"))
		})

		Step("allows navigating through browser history", func() {
			Expect(page.Back()).To(Succeed())
			Expect(page.URL()).NotTo(ContainSubstring("#new_page"))
			Expect(page.Forward()).To(Succeed())
			Expect(page.URL()).To(ContainSubstring("#new_page"))
		})

		Step("allows refreshing the page", func() {
			checkbox := page.Find("#some_checkbox")
			Check(checkbox)
			Expect(page.Refresh()).To(Succeed())
			Expect(checkbox).NotTo(BeSelected())
		})
	})
})
