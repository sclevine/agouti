Agouti
======

[![Build Status](https://api.travis-ci.org/sclevine/agouti.png?branch=master)](http://travis-ci.org/sclevine/agouti)

Integration testing for Go using Ginkgo 

Install:
```bash
$ go get github.com/sclevine/agouti
```
To use with PhantomJS (OS X):
```bash
$ brew install phantomjs
```
To use with ChromeDriver (OS X):
```bash
$ brew install chromedriver
```
To use with Selenium Webdriver (OS X):
```bash
$ brew install selenium-server-standalone
```
To use the `matcher` package, which provides Gomega matchers:
```bash
$ go get github.com/onsi/gomega
```
To use the `dsl` package, which defines tests that can be run with Ginkgo:
```bash
$ go get github.com/onsi/ginkgo/ginkgo
```

If you use the `dsl` package, note that:
 * `Feature` is a Ginkgo `Describe`
 * `Scenario` is a Ginkgo `It`
 * `Background` is a Ginkgo `BeforeEach`
 * `Step` is a Ginkgo `By`

Feel free to import Ginkgo and use any of its container blocks instead! Agouti is 100% compatible with Ginkgo and Gomega.

If you plan to use Agouti to write Ginkgo tests, add the start and stop commands for your choice of webdriver in Ginkgo `BeforeSuite` and `AfterSuite` blocks.

See this example `project_suite_test.go` file:
```Go
package project_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/page"

	"testing"
)

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project Suite")
}

var _ = BeforeSuite(func() {
	StartChrome()
	// and/or
	StartPhantom()
	// and/or
	StartSelenium()
});

var _ = AfterSuite(func() {
	StopChrome()
	// and/or
	StopPhantom()
	// and/or
	StopSelenium()
});
```

Example:

```Go
import (
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
)

...

Feature("Agouti", func() {
	Scenario("Loading a page", func() {
		page := CreatePage() // for PhantomJS
		// page := CreatePage("chrome") // for Chrome via ChromeDriver
		// page := CreatePage("firefox") // for Firefox via Selenium
		// page := CreatePage("safari") // for Safari via Selenium
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
```
