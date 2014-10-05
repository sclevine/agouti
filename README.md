Agouti
======

[![Build Status](https://api.travis-ci.org/sclevine/agouti.png?branch=master)](http://travis-ci.org/sclevine/agouti)

Integration testing for Go using Ginkgo 

Install (OS X):
```
brew install phantomjs
go get github.com/sclevine/agoati
```

Make sure to add the `defer CleanupAgouti(SetupAgouti())` to your `project_suite_test.go` file, like so:
```Go
package your_project_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"

	"testing"
)

func TestYourProject(t *testing.T) {
	RegisterFailHandler(Fail)
	defer CleanupAgouti(SetupAgouti())
	RunSpecs(t, "Your Project Suite")
}
```

Example:

```Go
import . "github.com/sclevine/agouti"

...

Feature("Agouti", func() {
  Scenario("Loads some page", func() {
    page := Navigate("http://example.com/")

    Step("finds a title", func() {
      page.Within("header").Within("h1").ShouldContainText("Page Title")
    })

    Within("#some-element", Do(func(someElement Selection) {
      someElement.Within("p").ShouldContainText("Foo")

      Step("and finds more text", func() {
        someElement.Within("[role=moreText]").ShouldContainText("Bar")
      })
    }))
  })
})
```
