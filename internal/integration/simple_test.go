package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var server *httptest.Server
var submitted bool

var _ = BeforeSuite(func() {
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		submitted = request.Method == "POST"
		html, _ := ioutil.ReadFile("test_page.html")
		response.Write(html)
	}))
})

var _ = AfterSuite(func() {
	server.Close()
})

var _ = Feature("Agouti", func() {
	Scenario("Loading a page with a cookie and clicking", func() {
		page := CreatePage()
		page.Size(640, 480)
		page.Navigate(server.URL)

		Step("finds text in a page", func() {
			Expect(page.Find("header")).To(ContainText("Page Title"))
		})

		Step("asserts that text is not in a page", func() {
			Expect(page).NotTo(ContainText("Page Not-Title"))
			Expect(page.Find("header")).NotTo(ContainText("Page Not-Title"))
		})

		Step("allows tests to be scoped by chaining", func() {
			Expect(page.Find("header").Find("h1")).To(ContainText("Page Title"))
		})

		Step("allows assertions that wait for matchers to be true", func() {
			Expect(page.Find("#some_element")).NotTo(ContainText("some text"))
			Eventually(page.Find("#some_element"), 4*time.Second).Should(ContainText("some text"))
			Consistently(page.Find("#some_element")).Should(ContainText("some text"))
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

		Step("allows submitting a form", func() {
			Expect(page.Find("#some_form").Submit()).To(Succeed())
			Eventually(submitted).Should(BeTrue())
		})
	})
})
