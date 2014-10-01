package agouti_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/sclevine/agouti"
)

var _ = BeforeSuite(func() {
	http.HandleFunc("/page", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("<header><h1>Page Title</h1></header>"))
	})

	go http.ListenAndServe(":42343", nil)
})

var _ = Describe("Agouti DSL", func() {
	Scenario("Loading a page", "http://localhost:42343/page", func() {
		Step("finds text in a page", func() {
			Within("header").ShouldContainText("Page Title")
		})

		Step("allows tests to be scoped by chaining", func() {
			Within("header").Within("h1").ShouldContainText("Page Title")
		})

		Step("allows tests to be scoped by functions", func() {
			Within("header h1", func(scope Scopable) {
				scope.ShouldContainText("Page Title")
			})
		})
	})
})
