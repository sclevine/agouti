package selection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
)

var _ = Describe("Selection", func() {
	var (
		client    *mocks.Client
		selection *Selection
		element   *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		selection = &Selection{Client: client}
		element = &mocks.Element{}
	})

	Describe("#Find", func() {
		It("should select a single element by CSS selector", func() {
			Expect(selection.Find("#selector").String()).To(Equal("CSS: #selector [single]"))
		})
	})

	Describe("#FindByXPath", func() {
		It("should select a single element by XPath", func() {
			Expect(selection.FindByXPath("//selector").String()).To(Equal("XPath: //selector [single]"))
		})
	})

	Describe("#FindByLink", func() {
		It("should select a single element by link text", func() {
			Expect(selection.FindByLink("some text").String()).To(Equal(`Link: "some text" [single]`))
		})
	})

	Describe("#FindByLabel", func() {
		It("should select a single element by label", func() {
			Expect(selection.FindByLabel("some label").String()).To(MatchRegexp("XPath: .+input.+ \\[single\\]"))
		})
	})

	Describe("#First", func() {
		It("should select the first element by CSS selector", func() {
			Expect(selection.First("#selector").String()).To(Equal("CSS: #selector [0]"))
		})
	})

	Describe("#FirstByXPath", func() {
		It("should select the first element by XPath", func() {
			Expect(selection.FirstByXPath("//selector").String()).To(Equal("XPath: //selector [0]"))
		})
	})

	Describe("#FirstByLink", func() {
		It("should select the first element by link text", func() {
			Expect(selection.FirstByLink("some text").String()).To(Equal(`Link: "some text" [0]`))
		})
	})

	Describe("#FirstByLabel", func() {
		It("should select the first element by label", func() {
			Expect(selection.FirstByLabel("some label").String()).To(MatchRegexp("XPath: .+input.+ \\[0\\]"))
		})
	})

	Describe("#All", func() {
		Context("when there is no selection", func() {
			It("should add a new CSS selector to the selection", func() {
				Expect(selection.All("#selector").String()).To(Equal("CSS: #selector"))
			})
		})

		Context("when the selection ends with an non-CSS selector", func() {
			It("should add a new selector to the selection", func() {
				xpath := selection.AllByXPath("//selector")
				Expect(xpath.All("#subselector").String()).To(Equal("XPath: //selector | CSS: #subselector"))
			})
		})

		Context("when the selection ends with an unindexed CSS selector", func() {
			It("should modify the last css selector to include the new selector", func() {
				Expect(selection.All("#selector").All("#subselector").String()).To(Equal("CSS: #selector #subselector"))
			})
		})

		Context("when the selection ends with an indexed selector", func() {
			It("should add a new selector to the selection", func() {
				Expect(selection.All("#selector").At(0).All("#subselector").String()).To(Equal("CSS: #selector [0] | CSS: #subselector"))
			})
		})

		Context("when the selection ends with a single-element-only selector", func() {
			It("should add a new selector to the selection", func() {
				Expect(selection.All("#selector").Single().All("#subselector").String()).To(Equal("CSS: #selector [single] | CSS: #subselector"))
			})
		})
	})

	Describe("#AllByXPath", func() {
		It("should add a new XPath selector to the selection", func() {
			Expect(selection.AllByXPath("//selector").String()).To(Equal("XPath: //selector"))
		})
	})

	Describe("#AllByLink", func() {
		It("should add a new 'link text' selector to the selection", func() {
			Expect(selection.AllByLink("some text").String()).To(Equal(`Link: "some text"`))
		})
	})

	Describe("#AllByLabel", func() {
		It("should add an XPath selector for finding by label", func() {
			Expect(selection.AllByLabel("label name").String()).To(Equal(`XPath: //input[@id=(//label[normalize-space(text())="label name"]/@for)] | //label[normalize-space(text())="label name"]/input`))
		})
	})

	Describe("selectors are always copied", func() {
		Context("when two CSS selections are created from the same XPath parent", func() {
			It("should not overwrite the first created child", func() {
				parent := selection.AllByXPath("//one").AllByXPath("//two").AllByXPath("//parent")
				firstChild := parent.All("#firstChild")
				parent.All("#secondChild")
				Expect(firstChild.String()).To(Equal("XPath: //one | XPath: //two | XPath: //parent | CSS: #firstChild"))
			})
		})
	})
})
