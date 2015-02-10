package selection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
)

var _ = Describe("Selection", func() {
	var (
		selection         *Selection
		session           *mocks.Session
		elementRepository *mocks.ElementRepository
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		elementRepository = &mocks.ElementRepository{}
		selection = &Selection{Session: session, Elements: elementRepository}
	})

	Describe(".NewSelection", func() {
		It("should return an empty selection with the provided session and element repository", func() {
			realSession := &api.Session{}
			selection = NewSelection(realSession)
			Expect(selection.Session).To(Equal(realSession))
			Expect(selection.Elements.(*ElementRepository).Client).To(Equal(realSession))
		})
	})

	Describe("#AppendCSS", func() {
		Context("when the selection ends with an unindexed CSS selector", func() {
			It("should modify the last CSS selector to include the new selector", func() {
				Expect(selection.AppendCSS("#selector").AppendCSS("#subselector").String()).To(Equal("CSS: #selector #subselector"))
			})

			It("should propagate the session and element repository", func() {
				Expect(selection.AppendCSS("#selector").AppendCSS("#subselector").Session).To(Equal(selection.Session))
				Expect(selection.AppendCSS("#selector").AppendCSS("#subselector").Elements).To(Equal(selection.Elements))
			})
		})

		It("should propagate the session and element repository in all other cases", func() {
			Expect(selection.AppendCSS("#selector").Session).To(Equal(selection.Session))
			Expect(selection.AppendCSS("#selector").Elements).To(Equal(selection.Elements))
		})

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
		It("should propagate the session and element repository", func() {
			Expect(selection.AppendXPath("//selector").Session).To(Equal(selection.Session))
			Expect(selection.AppendXPath("//selector").Elements).To(Equal(selection.Elements))
		})

		It("should add a new XPath selector to the selection", func() {
			Expect(selection.AppendXPath("//selector").String()).To(Equal("XPath: //selector"))
		})
	})

	Describe("#AppendLink", func() {
		It("should propagate the session and element repository", func() {
			Expect(selection.AppendLink("some text").Session).To(Equal(selection.Session))
			Expect(selection.AppendLink("some text").Elements).To(Equal(selection.Elements))
		})

		It("should add a new 'link text' selector to the selection", func() {
			Expect(selection.AppendLink("some text").String()).To(Equal(`Link: "some text"`))
		})
	})

	Describe("#AppendLabeled", func() {
		It("should propagate the session and element repository", func() {
			Expect(selection.AppendLabeled("some text").Session).To(Equal(selection.Session))
			Expect(selection.AppendLabeled("some text").Elements).To(Equal(selection.Elements))
		})

		It("should add a new XPath label-lookup selector to the selection", func() {
			Expect(selection.AppendLabeled("some text").String()).To(Equal(`XPath: //input[@id=(//label[normalize-space()="some text"]/@for)] | //label[normalize-space()="some text"]/input`))
		})
	})

	Describe("#AppendButton", func() {
		It("should propagate the session and element repository", func() {
			Expect(selection.AppendButton("some text").Session).To(Equal(selection.Session))
			Expect(selection.AppendButton("some text").Elements).To(Equal(selection.Elements))
		})

		It("should add a new XPath label-lookup selector to the selection", func() {
			Expect(selection.AppendButton("some text").String()).To(Equal(`XPath: //input[@type="submit" or @type="button"][normalize-space(@value)="some text"] | //button[normalize-space()="some text"]`))
		})
	})

	Describe("#At", func() {
		Context("when called on a selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(selection.At(1).String()).To(Equal(""))
			})

			It("should propagate the session and element repository", func() {
				Expect(selection.At(0).Session).To(Equal(selection.Session))
				Expect(selection.At(0).Elements).To(Equal(selection.Elements))
			})
		})

		Context("when called on a selection with selectors", func() {
			It("should select an index of the current selection", func() {
				Expect(selection.AppendCSS("#selector").At(1).String()).To(Equal("CSS: #selector [1]"))
			})

			It("should propagate the session and element repository", func() {
				Expect(selection.AppendCSS("#selector").At(0).Session).To(Equal(selection.Session))
				Expect(selection.AppendCSS("#selector").At(0).Elements).To(Equal(selection.Elements))
			})
		})
	})

	Describe("#Single", func() {
		Context("when called on a selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(selection.Single().String()).To(Equal(""))
			})

			It("should propagate the session and element repository", func() {
				Expect(selection.Single().Session).To(Equal(selection.Session))
				Expect(selection.Single().Elements).To(Equal(selection.Elements))
			})
		})

		Context("when called on a selection with selectors", func() {
			It("should select a single element of the current selection", func() {
				Expect(selection.AppendCSS("#selector").Single().String()).To(Equal("CSS: #selector [single]"))
			})

			It("should propagate the session and element repository", func() {
				Expect(selection.AppendCSS("#selector").Single().Session).To(Equal(selection.Session))
				Expect(selection.AppendCSS("#selector").Single().Elements).To(Equal(selection.Elements))
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
