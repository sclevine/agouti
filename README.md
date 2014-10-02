Agouti
======

[![Build Status](https://api.travis-ci.org/sclevine/agouti.png?branch=master)](http://travis-ci.org/sclevine/agouti)

Integration testing for Go using Ginkgo 

Install (OS X):
```
brew install phantomjs
go get github.com/sclevine/agoati
```


Example:

```Go
import . "github.com/sclevine/agouti"

...

Describe("Agouti", func() {
  Scenario("Loads some page", "http://example.com/", func() {
    Step("finds a title", func() {
      Within("header").Within("h1").ShouldContainText("Page Title")
    })

    Within("#some-element", func(scope Scopable) {
      scope.Within("p").ShouldContainText("Foo")

      Step("and finds more text", func() {
        scope.Within("[role=moreText]").ShouldContainText("Bar")
      })
    })
  })
})
```
