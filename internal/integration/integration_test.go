package integration_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/internal/integration"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Feature("Agouti running on PhantomJS", func() {
	var page Page

	Background(func() {
		page = CreatePage()
		page.Size(640, 480)
		page.Navigate(Server.URL)
	})

	AfterEach(func() {
		page.Destroy()
	})

	Scenario("finding the page title and URL", func() {
		Expect(page).To(HaveTitle("Page Title"))
		Expect(page).To(HaveURL(Server.URL + "/"))
	})

	Scenario("finding page elements", func() {
		Step("finding a header in the page", func() {
			Expect(page.Find("header")).To(BeFound())
			Expect(page.Find("not-a-header")).NotTo(BeFound())
		})

		Step("finding text in the header", func() {
			Expect(page.Find("header")).To(HaveText("Title"))
		})

		Step("asserting that text is not in the header", func() {
			Expect(page.Find("header")).NotTo(HaveText("Not-Title"))
		})

		Step("referring to an element by selection index", func() {
			Expect(page.All("option").At(0)).To(HaveText("first option"))
			Expect(page.All("select").At(1).First("option")).To(HaveText("third option"))
		})

		Step("matching text in the header", func() {
			Expect(page.Find("header")).To(MatchText("T.+e"))
		})

		Step("scoping selections by chaining", func() {
			Expect(page.Find("header").Find("h1")).To(HaveText("Title"))
		})

		Step("locating elements by XPath", func() {
			Expect(page.Find("header").FindByXPath("//h1")).To(HaveText("Title"))
		})

		Step("comparing two selections for equality", func() {
			Expect(page.Find("#some_element")).To(EqualElement(page.FindByXPath("//div[@class='some-element']")))
		})
	})

	Scenario("selecting multiple elements", func() {
		Step("asserting on their state", func() {
			Expect(page.All("select").All("option")).To(BeVisible())
			Expect(page.All("h1,h2")).NotTo(BeVisible())
		})
	})

	Scenario("finding form elements by label", func() {
		Step("finding an element by label text", func() {
			Expect(page.FindByLabel("Some Label")).To(HaveAttribute("value", "some labeled value"))
		})

		Step("finding an element embedded in a label", func() {
			Expect(page.FindByLabel("Some Container Label")).To(HaveAttribute("value", "some embedded value"))
		})
	})

	Scenario("element properties", func() {
		Step("finding visible elements", func() {
			Expect(page.Find("header h1")).To(BeVisible())
			Expect(page.Find("header h2")).NotTo(BeVisible())
		})

		Step("finding enabled elements", func() {
			Expect(page.Find("#some_checkbox")).To(BeEnabled())
			Expect(page.Find("#some_disabled_checkbox")).NotTo(BeEnabled())
		})

		Step("finding the active element", func() {
			Expect(page.Find("#some_checkbox")).NotTo(BeActive())
			Click(page.Find("#some_checkbox"))
			Expect(page.Find("#some_checkbox")).To(BeActive())
		})
	})

	Scenario("asynchronous javascript and DOM assertions", func() {
		Step("waiting for matchers to be true", func() {
			Expect(page.Find("#some_element")).NotTo(HaveText("some text"))
			Eventually(page.Find("#some_element"), 4*time.Second).Should(HaveText("some text"))
			Consistently(page.Find("#some_element")).Should(HaveText("some text"))
		})

		Step("serializing the current page HTML", func() {
			Expect(page.HTML()).To(ContainSubstring(`>some text</div>`))
		})

		Step("executing arbitrary javascript", func() {
			arguments := map[string]interface{}{"elementID": "some_element"}
			var result string
			Expect(page.RunScript("return document.getElementById(elementID).innerHTML;", arguments, &result)).To(Succeed())
			Expect(result).To(Equal("some text"))
		})
	})

	Scenario("filling fields and asserting on their values", func() {
		Step("entering values into fields", func() {
			Fill(page.Find("#some_input"), "some other value")
		})

		Step("retrieving attributes by name", func() {
			Expect(page.Find("#some_input")).To(HaveAttribute("value", "some other value"))
		})
	})

	Scenario("CSS styles", func() {
		Expect(page.Find("#some_element")).To(HaveCSS("color", "rgba(0, 0, 255, 1)"))
		Expect(page.Find("#some_element")).To(HaveCSS("color", "rgb(0, 0, 255)"))
		Expect(page.Find("#some_element")).To(HaveCSS("color", "blue"))
	})

	Scenario("form actions", func() {
		Step("double-clicking on an element", func() {
			selection := page.Find("#double_click")
			DoubleClick(selection)
			Expect(selection).To(HaveText("double-click success"))
		})

		Step("checking a checkbox", func() {
			checkbox := page.Find("#some_checkbox")
			Check(checkbox)
			Expect(checkbox).To(BeSelected())
		})

		Step("selecting an option by text", func() {
			selection := page.Find("#some_select")
			Select(selection, "second option")
			Expect(selection.Find("option:last-child")).To(BeSelected())
		})

		Step("submitting a form", func() {
			Submit(page.Find("#some_form"))
			Eventually(Submitted).Should(BeTrue())
		})
	})

	Scenario("links and navigation", func() {
		Step("allows clicking on a link", func() {
			Click(page.FindByLink("Click Me"))
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

	Scenario("logging", func() {
		Eventually(page).Should(HaveLoggedInfo("some log"))
		Expect(page).NotTo(HaveLoggedError())
		Eventually(page, 4*time.Second).Should(HaveLoggedError("ReferenceError: Can't find variable: doesNotExist\n  (anonymous function)"))
	})

	Scenario("frames", func() {
		Step("allows switching to an iframe", func() {
			SwitchToFrame(page.Find("#frame"))
			Expect(page.Find("body")).To(MatchText("Example Domain"))
		})

		// PhantomJS does not seem to comply to /session/:id/frame/parent

		//Step("allows switching back to the default frame by referring to the parent frame", func() {
		//	SwitchToParentFrame(page)
		//	Expect(page.Find("body")).NotTo(MatchText("Example Domain"))
		//})
		//SwitchToFrame(page.Find("#frame"))

		Step("allows switching back to the default frame by referring to the root frame", func() {
			SwitchToRootFrame(page)
			Expect(page.Find("body")).NotTo(MatchText("Example Domain"))
		})
	})

	Scenario("windows", func() {
		Click(page.Find("#new_window"))
		windows, _ := page.WindowCount()
		Expect(windows).To(Equal(2))

		Step("allows switching windows", func() {
			Expect(page.SwitchToWindow("new window")).To(Succeed())
			Expect(page.Find("header")).NotTo(BeFound())
			Expect(page.NextWindow()).To(Succeed())
			Expect(page.Find("header")).To(BeFound())
		})

		Step("allows closing windows", func() {
			Expect(page.CloseWindow()).To(Succeed())
			windows, _ := page.WindowCount()
			Expect(windows).To(Equal(1))
		})
	})

	// PhantomJS has trouble with alert boxes, even with handlesAlerts
	PScenario("popup boxes", func() {
		Step("allows interacting with alert popups", func() {
			Click(page.Find("#popup_alert"))
			Expect(page).To(HavePopupText("some alert"))
			Expect(page.ConfirmPopup()).To(Succeed())
		})

		Step("allows interacting with confirm boxes", func() {
			var confirmed bool

			Click(page.Find("#popup_confirm"))
			Expect(page.ConfirmPopup()).To(Succeed())
			Expect(page.RunScript("return confirmed;", nil, &confirmed)).To(Succeed())
			Expect(confirmed).To(BeTrue())

			Click(page.Find("#popup_confirm"))
			Expect(page.RunScript("return confirmed;", nil, &confirmed)).To(Succeed())
			Expect(page.CancelPopup()).To(Succeed())
			Expect(confirmed).To(BeFalse())
		})

		Step("allows interacting with prompt boxes", func() {
			var promptText string

			Click(page.Find("#popup_prompt"))
			Expect(page.EnterPopupText("banana")).To(Succeed())
			Expect(page.ConfirmPopup()).To(Succeed())
			Expect(page.RunScript("return promptText;", nil, &promptText)).To(Succeed())
			Expect(promptText).To(Equal("banana"))
		})
	})
})
