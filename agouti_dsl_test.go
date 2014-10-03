package agouti_test

import (
	. "github.com/sclevine/agouti"
	"net/http"
	"net/http/httptest"
)

var server *httptest.Server

var _ = BeforeSuite(func() {
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("<header><h1>Page Title</h1></header>"))
	}))
})

var _ = AfterSuite(func() {
	server.Close()
})

var _ = Feature("Agouti DSL", func() {
	Scenario("Loading a page", func() {
		page := Navigate(server.URL)

		Step("finds text in a page", func() {
			page.Within("header").ShouldContainText("Page Title")
		})

		Step("allows tests to be scoped by chaining", func() {
			page.Within("header").Within("h1").ShouldContainText("Page Title")
		})

		Step("allows tests to be scoped by functions", func() {
			page.Within("header h1", func(h1 *Selection) {
				h1.ShouldContainText("Page Title")
			})
		})
	})
})
