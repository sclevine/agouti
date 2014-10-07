package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"net/http"
	"net/http/httptest"
)

var server *httptest.Server

var _ = BeforeSuite(func() {
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte(`<html><body><header><h1>Page Title</h1></header><a href="#new_page">Click Me</a></body></html>`))
	}))
})

var _ = AfterSuite(func() {
	server.Close()
})

var _ = Feature("Agouti", func() {
	Scenario("Loading a page with a cookie and clicking", func() {
		cookie := CreateCookie("theName", 42, "/my-path", "example.com", false, false, 1412358590)
		page := CreatePage()
		page.Navigate(server.URL).SetCookie(cookie)

		Step("finds text in a page", func() {
			page.Should().ContainText("Page Title")
			page.Within("header").Should().ContainText("Page Title")
		})

		Step("assets that text is not in a page", func() {
			page.ShouldNot().ContainText("Page Not-Title")
			page.Within("header").ShouldNot().ContainText("Page Not-Title")
		})

		Step("allows tests to be scoped by chaining", func() {
			page.Within("header").Within("h1").Should().ContainText("Page Title")
		})

		Step("allows tests to be scoped by functions", func() {
			page.Within("header h1", Do(func(h1 Selection) {
				h1.Should().ContainText("Page Title")
			}))
		})

		Step("allows clicking on a link", func() {
			page.Within("a").Click()
			Expect(page.URL()).To(ContainSubstring("#new_page"))
		})
	})
})
