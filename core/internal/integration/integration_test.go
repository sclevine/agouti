package integration_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Deprecated integration tests", func() {
	var (
		page      Page
		submitted bool
		server    *httptest.Server
	)

	BeforeEach(func() {
		server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			if request.Method == "POST" {
				submitted = true
			}
			html, _ := ioutil.ReadFile("test_page.html")
			response.Write(html)
		}))

		var err error
		page, err = phantomDriver.Page()
		Expect(err).NotTo(HaveOccurred())
		Expect(page.Size(640, 480)).To(Succeed())
		Expect(page.Navigate(server.URL)).To(Succeed())
	})

	AfterEach(func() {
		page.Destroy()
		server.Close()
	})

	It("should support finding the page title and URL", func() {
		Expect(page).To(HaveTitle("Page Title"))
		Expect(page).To(HaveURL(server.URL + "/"))
	})

	It("should support finding page elements", func() {
		By("finding a header in the page", func() {
			Expect(page.Find("header")).To(BeFound())
			Expect(page.Find("not-a-header")).NotTo(BeFound())
		})

		By("finding text in the header", func() {
			Expect(page.Find("header")).To(HaveText("Title"))
		})

		By("asserting that text is not in the header", func() {
			Expect(page.Find("header")).NotTo(HaveText("Not-Title"))
		})

		By("referring to an element by selection index", func() {
			Expect(page.All("option").At(0)).To(HaveText("first option"))
			Expect(page.All("select").At(1).First("option")).To(HaveText("third option"))
		})

		By("matching text in the header", func() {
			Expect(page.Find("header")).To(MatchText("T.+e"))
		})

		By("scoping selections by chaining", func() {
			Expect(page.Find("header").Find("h1")).To(HaveText("Title"))
		})

		By("locating elements by XPath", func() {
			Expect(page.Find("header").FindByXPath("//h1")).To(HaveText("Title"))
		})

		By("comparing two selections for equality", func() {
			Expect(page.Find("#some_element")).To(EqualElement(page.FindByXPath("//div[@class='some-element']")))
		})
	})

	It("should support selecting multiple elements", func() {
		By("asserting on their state", func() {
			Expect(page.All("select").All("option")).To(BeVisible())
			Expect(page.All("h1,h2")).NotTo(BeVisible())
		})
	})

	It("should support finding form elements by label", func() {
		By("finding an element by label text", func() {
			Expect(page.FindByLabel("Some Label")).To(HaveAttribute("value", "some labeled value"))
		})

		By("finding an element embedded in a label", func() {
			Expect(page.FindByLabel("Some Container Label")).To(HaveAttribute("value", "some embedded value"))
		})
	})

	It("should support finding button elements by text", func() {
		By("finding a <button>", func() {
			Expect(page.FindByButton("Some Button")).To(HaveAttribute("name", "some button name"))
		})

		By("finding an <input> button", func() {
			Expect(page.FindByButton("Some Input Button")).To(HaveAttribute("type", "button"))
			Expect(page.FindByButton("Some Submit Button")).To(HaveAttribute("type", "submit"))
		})
	})

	It("should support asserting on element properties", func() {
		By("finding visible elements", func() {
			Expect(page.Find("header h1")).To(BeVisible())
			Expect(page.Find("header h2")).NotTo(BeVisible())
		})

		By("finding enabled elements", func() {
			Expect(page.Find("#some_checkbox")).To(BeEnabled())
			Expect(page.Find("#some_disabled_checkbox")).NotTo(BeEnabled())
		})

		By("finding the active element", func() {
			Expect(page.Find("#some_checkbox")).NotTo(BeActive())
			Expect(page.Find("#some_checkbox").Click()).To(Succeed())
			Expect(page.Find("#some_checkbox")).To(BeActive())
		})
	})

	It("should support running JavaScript and making assertions on the DOM", func() {
		By("waiting for matchers to be true", func() {
			Expect(page.Find("#some_element")).NotTo(HaveText("some text"))
			Eventually(page.Find("#some_element"), "4s").Should(HaveText("some text"))
			Consistently(page.Find("#some_element")).Should(HaveText("some text"))
		})

		By("serializing the current page HTML", func() {
			Expect(page.HTML()).To(ContainSubstring(`>some text</div>`))
		})

		By("executing arbitrary javascript", func() {
			arguments := map[string]interface{}{"elementID": "some_element"}
			var result string
			Expect(page.RunScript("return document.getElementById(elementID).innerHTML;", arguments, &result)).To(Succeed())
			Expect(result).To(Equal("some text"))
		})
	})

	It("should support filling out fields and asserting on their values", func() {
		By("entering values into fields", func() {
			Expect(page.Find("#some_input").Fill("some other value")).To(Succeed())
		})

		By("retrieving attributes by name", func() {
			Expect(page.Find("#some_input")).To(HaveAttribute("value", "some other value"))
		})
	})

	It("should support matching CSS styles", func() {
		Expect(page.Find("#some_element")).To(HaveCSS("color", "rgba(0, 0, 255, 1)"))
		Expect(page.Find("#some_element")).To(HaveCSS("color", "rgb(0, 0, 255)"))
		Expect(page.Find("#some_element")).To(HaveCSS("color", "blue"))
	})

	It("should support form actions", func() {
		By("double-clicking on an element", func() {
			selection := page.Find("#double_click")
			Expect(selection.DoubleClick()).To(Succeed())
			Expect(selection).To(HaveText("double-click success"))
		})

		By("checking a checkbox", func() {
			checkbox := page.Find("#some_checkbox")
			Expect(checkbox.Check()).To(Succeed())
			Expect(checkbox).To(BeSelected())
		})

		By("selecting an option by text", func() {
			selection := page.Find("#some_select")
			Expect(selection.Select("second option")).To(Succeed())
			Expect(selection.Find("option:last-child")).To(BeSelected())
		})

		By("submitting a form", func() {
			Expect(page.Find("#some_form").Submit()).To(Succeed())
			Eventually(func() bool { return submitted }).Should(BeTrue())
		})
	})

	It("should support links and navigation", func() {
		By("clicking on a link", func() {
			Expect(page.FindByLink("Click Me").Click()).To(Succeed())
			Expect(page.URL()).To(ContainSubstring("#new_page"))
		})

		By("navigating through browser history", func() {
			Expect(page.Back()).To(Succeed())
			Expect(page.URL()).NotTo(ContainSubstring("#new_page"))
			Expect(page.Forward()).To(Succeed())
			Expect(page.URL()).To(ContainSubstring("#new_page"))
		})

		By("refreshing the page", func() {
			checkbox := page.Find("#some_checkbox")
			Expect(checkbox.Check()).To(Succeed())
			Expect(page.Refresh()).To(Succeed())
			Expect(checkbox).NotTo(BeSelected())
		})
	})

	It("should support retrieving logs", func() {
		Eventually(page).Should(HaveLoggedInfo("some log"))
		Expect(page).NotTo(HaveLoggedError())
		Eventually(page, "4s").Should(HaveLoggedError("ReferenceError: Can't find variable: doesNotExist\n  (anonymous function)"))
	})

	It("should support switching frames", func() {
		By("switching to an iframe", func() {
			Expect(page.Find("#frame").SwitchToFrame()).To(Succeed())
			Expect(page.Find("body")).To(MatchText("Example Domain"))
		})

		By("switching back to the default frame by referring to the root frame", func() {
			Expect(page.SwitchToRootFrame()).To(Succeed())
			Expect(page.Find("body")).NotTo(MatchText("Example Domain"))
		})
	})

	It("should support switching windows", func() {
		Expect(page.Find("#new_window").Click()).To(Succeed())
		windows, _ := page.WindowCount()
		Expect(windows).To(Equal(2))

		By("switching windows", func() {
			Expect(page.SwitchToWindow("new window")).To(Succeed())
			Expect(page.Find("header")).NotTo(BeFound())
			Expect(page.NextWindow()).To(Succeed())
			Expect(page.Find("header")).To(BeFound())
		})

		By("closing windows", func() {
			Expect(page.CloseWindow()).To(Succeed())
			windows, _ := page.WindowCount()
			Expect(windows).To(Equal(1))
		})
	})
})
