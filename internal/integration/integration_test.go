package integration_test

import (
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/internal/integration"
	. "github.com/sclevine/agouti/matchers"
	"time"
)

var _ = Feature("Agouti running on PhantomJS", func() {
	Scenario("Loading a page", func() {
		page := CreatePage()
		page.Size(640, 480)
		page.Navigate(Server.URL)

		Step("find the title of the page", func() {
			Expect(page).To(HaveTitle("Page Title"))
		})

		Step("finds a header in the page", func() {
			Expect(page.Find("header")).To(BeFound())
		})

		Step("finds text in the header", func() {
			Expect(page.Find("header")).To(HaveText("Title"))
		})

		Step("asserts that text is not in the header", func() {
			Expect(page.Find("header")).NotTo(HaveText("Not-Title"))
		})

		Step("asserts on the visibility of elements", func() {
			Expect(page.Find("header h1")).To(BeVisible())
			Expect(page.Find("header h2")).NotTo(BeVisible())
		})

		Step("allows tests to be scoped by chaining", func() {
			Expect(page.Find("header").Find("h1")).To(HaveText("Title"))
		})

		Step("allows assertions that wait for matchers to be true", func() {
			Expect(page.Find("#some_element")).NotTo(HaveText("some text"))
			Eventually(page.Find("#some_element"), 4*time.Second).Should(HaveText("some text"))
			Consistently(page.Find("#some_element")).Should(HaveText("some text"))
		})

		Step("allows entering values into fields", func() {
			Expect(page.Find("#some_input").Fill("some other value")).To(Succeed())
		})

		Step("allows retrieving attributes by name", func() {
			Expect(page.Find("#some_input")).To(HaveAttribute("value", "some other value"))
		})

		Step("allows asserting on whether a CSS style exists", func() {
			Expect(page.Find("#some_element")).To(HaveCSS("color", "rgba(0, 0, 255, 1)"))
		})

		Step("allows clicking on a link", func() {
			Expect(page.Find("a").Click()).To(Succeed())
			Expect(page.URL()).To(ContainSubstring("#new_page"))
		})

		Step("allows double-clicking on an element", func() {
			selection := page.Find("#double_click")
			Expect(selection.DoubleClick()).To(Succeed())
			Expect(selection).To(HaveText("double-click success"))
		})

		Step("allows checking a checkbox", func() {
			checkbox := page.Find("#some_checkbox")
			Expect(checkbox.Check()).To(Succeed())
			Expect(checkbox).To(BeSelected())
		})

		Step("allows selecting an option", func() {
			selection := page.Find("#some_select")
			Expect(selection.Select("second option")).To(Succeed())
			Expect(selection.Find("option:last-child")).To(BeSelected())
		})

		Step("allows executing arbitrary javascript", func() {
			arguments := map[string]interface{}{"elementID": "some_element"}
			var result string
			Expect(page.RunScript("return document.getElementById(elementID).innerHTML;", arguments, &result)).To(Succeed())
			Expect(result).To(Equal("some text"))
		})

		Step("allows submitting a form", func() {
			Expect(page.Find("#some_form").Submit()).To(Succeed())
			Eventually(Submitted).Should(BeTrue())
		})
	})
})
