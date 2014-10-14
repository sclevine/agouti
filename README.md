Agouti
======

[![Build Status](https://api.travis-ci.org/sclevine/agouti.png?branch=master)](http://travis-ci.org/sclevine/agouti)

Integration testing for Go using Ginkgo and Gomega!

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

The `core` package is a flexible, general-purpose webdriver API for Go. Unlike the `dsl` package, `core` allows unlimited and simultaneous usage of PhantomJS, ChromeDriver, Selenium.

If you plan to use Agouti to write Ginkgo tests, add the start and stop commands for your choice of webdriver in Ginkgo `BeforeSuite` and `AfterSuite` blocks.

See this example `project_suite_test.go` file:
```Go
package project_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"

	"testing"
)

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project Suite")
}

var _ = BeforeSuite(func() {
	StartPhantomJS()
	// OR
	StartChrome()
	// OR
	StartSelenium()
});

var _ = AfterSuite(func() {
	StopWebdriver()
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

Feature("Agouti running on PhantomJS", func() {
	Scenario("Loading a page", func() {
		page := CreatePage()
		page.Size(640, 480)
		page.Navigate(Server.URL)

		Step("find the title of the page", func() {
			Expect(page).To(HaveTitle("Page Title"))
		})

		Step("finds text in a page", func() {
			Expect(page.Find("header")).To(HaveText("Title"))
		})

		Step("asserts that text is not in the header", func() {
			Expect(page.Find("header")).NotTo(HaveText("Not-Title"))
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
			Eventually(Submitted).Should(BeTrue())
		})
	})
})
```
