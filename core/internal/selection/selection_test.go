package selection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
)

var _ = Describe("Selection", func() {
	var (
		selection *Selection
		client    *mocks.Client
		element   *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		selection = &Selection{Client: client}
		element = &mocks.Element{}
	})

	Describe("#AppendCSS", func() {
		Context("when there is no selection", func() {
			It("should add a new CSS selector to the selection", func() {
				Expect(selection.AppendCSS("#selector").String()).To(Equal("CSS: #selector"))
			})
		})

		Context("when the selection ends with an non-CSS selector", func() {
			It("should add a new selector to the selection", func() {
				xpath := selection.AppendXPath("//selector")
				Expect(xpath.AppendCSS("#subselector").String()).To(Equal("XPath: //selector | CSS: #subselector"))
			})
		})

		Context("when the selection ends with an unindexed CSS selector", func() {
			It("should modify the last CSS selector to include the new selector", func() {
				Expect(selection.AppendCSS("#selector").AppendCSS("#subselector").String()).To(Equal("CSS: #selector #subselector"))
			})
		})

		Context("when the selection ends with an indexed selector", func() {
			It("should add a new selector to the selection", func() {
				Expect(selection.AppendCSS("#selector").At(0).AppendCSS("#subselector").String()).To(Equal("CSS: #selector [0] | CSS: #subselector"))
			})
		})

		Context("when the selection ends with a single-element-only selector", func() {
			It("should add a new selector to the selection", func() {
				Expect(selection.AppendCSS("#selector").Single().AppendCSS("#subselector").String()).To(Equal("CSS: #selector [single] | CSS: #subselector"))
			})
		})
	})

	Describe("#AppendXPath", func() {
		It("should add a new XPath selector to the selection", func() {
			Expect(selection.AppendXPath("//selector").String()).To(Equal("XPath: //selector"))
		})
	})

	Describe("#AppendLink", func() {
		It("should add a new 'link text' selector to the selection", func() {
			Expect(selection.AppendLink("some text").String()).To(Equal(`Link: "some text"`))
		})
	})

	Describe("#AppendLabeled", func() {
		It("should add a new XPath label-lookup selector to the selection", func() {
			Expect(selection.AppendLabeled("some text").String()).To(Equal(`XPath: //input[@id=(//label[normalize-space(text())="some text"]/@for)] | //label[normalize-space(text())="some text"]/input`))
		})
	})

	Describe("#At", func() {
		Context("when called on a selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(selection.At(1).String()).To(Equal(""))
			})
		})

		Context("when called on a selection with selectors", func() {
			It("should select an index of the current selection", func() {
				Expect(selection.AppendCSS("#selector").At(1).String()).To(Equal("CSS: #selector [1]"))
			})
		})
	})

	Describe("#Single", func() {
		Context("when called on a selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(selection.Single().String()).To(Equal(""))
			})
		})

		Context("when called on a selection with selectors", func() {
			It("should select a single element of the current selection", func() {
				Expect(selection.AppendCSS("#selector").Single().String()).To(Equal("CSS: #selector [single]"))
			})
		})
	})

	Describe("selectors are always copied", func() {
		Context("when two CSS selections are created from the same XPath parent", func() {
			It("should not overwrite the first created child", func() {
				parent := selection.AppendXPath("//one").AppendXPath("//two").AppendXPath("//parent")
				firstChild := parent.AppendCSS("#firstChild")
				parent.AppendCSS("#secondChild")
				Expect(firstChild.String()).To(Equal("XPath: //one | XPath: //two | XPath: //parent | CSS: #firstChild"))
			})
		})
	})
})
