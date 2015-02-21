package agouti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/mocks"
)

var _ = Describe("Selectable", func() {
	var (
		bus       *mocks.Bus
		session   *api.Session
		selection *Selection
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		session = &api.Session{bus}
		selection = NewTestSelection(nil, session, "#test")
	})

	type finder interface {
		String() string
		Find(selector string) *Selection
	}

	itShouldSelect := func(method func(string) finder, expected string) {
		It("should apply the appropriate selectors", func() {
			Expect(method("selector").String()[22:]).To(Equal(expected))
		})

		It("should provide the selectable's session to the element repository", func() {
			bus.SendCall.Result = `[{"ELEMENT": "some-id"}]`
			elements, _ := method("a").Find("b").Elements()
			Expect(elements[0].ID).To(Equal("some-id"))
		})
	}

	Describe("#Find", func() {
		method := func(selector string) finder { return selection.Find(selector) }
		itShouldSelect(method, "CSS: selector [single]")
	})

	Describe("#FindByXPath", func() {
		method := func(selector string) finder { return selection.FindByXPath(selector) }
		itShouldSelect(method, "XPath: selector [single]")
	})

	Describe("#FindByLink", func() {
		method := func(selector string) finder { return selection.FindByLink(selector) }
		itShouldSelect(method, `Link: "selector" [single]`)
	})

	Describe("#FindByLabel", func() {
		method := func(selector string) finder { return selection.FindByLabel(selector) }
		itShouldSelect(method, `Label: "selector" [single]`)
	})

	Describe("#FindByButton", func() {
		method := func(selector string) finder { return selection.FindByButton(selector) }
		itShouldSelect(method, `Button: "selector" [single]`)
	})

	Describe("#First", func() {
		method := func(selector string) finder { return selection.First(selector) }
		itShouldSelect(method, "CSS: selector [0]")
	})

	Describe("#FirstByXPath", func() {
		method := func(selector string) finder { return selection.FirstByXPath(selector) }
		itShouldSelect(method, "XPath: selector [0]")
	})

	Describe("#FirstByLink", func() {
		method := func(selector string) finder { return selection.FirstByLink(selector) }
		itShouldSelect(method, `Link: "selector" [0]`)
	})

	Describe("#FirstByLabel", func() {
		method := func(selector string) finder { return selection.FirstByLabel(selector) }
		itShouldSelect(method, `Label: "selector" [0]`)
	})

	Describe("#FirstByButton", func() {
		method := func(selector string) finder { return selection.FirstByButton(selector) }
		itShouldSelect(method, `Button: "selector" [0]`)
	})

	Describe("#All", func() {
		method := func(selector string) finder { return selection.All(selector) }
		itShouldSelect(method, "CSS: selector")
	})

	Describe("#AllByXPath", func() {
		method := func(selector string) finder { return selection.AllByXPath(selector) }
		itShouldSelect(method, "XPath: selector")
	})

	Describe("#AllByLink", func() {
		method := func(selector string) finder { return selection.AllByLink(selector) }
		itShouldSelect(method, `Link: "selector"`)
	})

	Describe("#AllByLabel", func() {
		method := func(selector string) finder { return selection.AllByLabel(selector) }
		itShouldSelect(method, `Label: "selector"`)
	})

	Describe("#AllByButton", func() {
		method := func(selector string) finder { return selection.AllByButton(selector) }
		itShouldSelect(method, `Button: "selector"`)
	})
})
