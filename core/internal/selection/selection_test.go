package selection_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
)

var _ = Describe("Selection", func() {
	var (
		client    *mocks.Client
		selection types.Selection
		element   *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		selection = &Selection{Client: client}
		element = &mocks.Element{}
	})

	Describe("#At", func() {
		Context("when called on a selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(selection.(types.MultiSelection).At(1).String()).To(Equal(""))
			})
		})

		Context("when called on a selection with selectors ", func() {
			It("should select an index of the current selectino", func() {
				Expect(selection.All("#selector").At(1).String()).To(Equal("CSS: #selector [1]"))
			})
		})
	})

	Describe("#Find", func() {
		It("should select the first of all of the elements by CSS", func() {
			Expect(selection.Find("#selector").String()).To(Equal("CSS: #selector [0]"))
		})
	})

	Describe("#FindByXPath", func() {
		It("should select the first of all of the elements by XPath", func() {
			Expect(selection.FindByXPath("//selector").String()).To(Equal("XPath: //selector [0]"))
		})
	})

	Describe("#FindByLink", func() {
		It("should select the first of all of the elements by link text", func() {
			Expect(selection.FindByLink("some text").String()).To(Equal(`Link: "some text" [0]`))
		})
	})

	Describe("#FindByLabel", func() {
		It("should select the first of all of the elements by label", func() {
			Expect(selection.FindByLabel("some label").String()).To(MatchRegexp("XPath: .+input.+ \\[0\\]"))
		})
	})

	Describe("#All", func() {
		Context("when there is no selection", func() {
			It("should add a new css selector to the selection", func() {
				Expect(selection.All("#selector").String()).To(Equal("CSS: #selector"))
			})
		})

		Context("when the selection ends with an non-css selector", func() {
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

	Describe("#Count", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element, element}
			selection = selection.All("#selector")
		})

		It("should request elements from the client using the provided selector", func() {
			selection.Count()
			Expect(client.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "#selector"}))
		})

		Context("when the client succeeds in retrieving the elements", func() {
			It("should return the text", func() {
				count, _ := selection.Count()
				Expect(count).To(Equal(2))
			})

			It("should not return an error", func() {
				_, err := selection.Count()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the the client fails to retrieve the elements", func() {
			BeforeEach(func() {
				client.GetElementsCall.Err = errors.New("some error")
			})

			It("should return an error", func() {
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#EqualsElement", func() {
		var (
			otherClient    *mocks.Client
			otherSelection types.Selection
			otherElement   *mocks.Element
		)

		BeforeEach(func() {
			selection = selection.All("#selector")
			client.GetElementsCall.ReturnElements = []types.Element{element}
			otherClient = &mocks.Client{}
			otherSelection = &Selection{Client: otherClient}
			otherSelection = otherSelection.All("#other_selector")
			otherElement = &mocks.Element{}
			otherClient.GetElementsCall.ReturnElements = []types.Element{otherElement}
		})

		Context("when multiple elements are selected from the selection", func() {
			It("should return an error with the number of elements", func() {
				client.GetElementsCall.ReturnElements = []types.Element{element, element}
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).To(MatchError("failed to select 'CSS: #selector': method does not support multiple elements (2)"))
			})
		})

		Context("when multiple elements are selected from the other selection", func() {
			It("should return an error with the number of elements", func() {
				otherClient.GetElementsCall.ReturnElements = []types.Element{element, element}
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).To(MatchError("failed to select 'CSS: #other_selector': method does not support multiple elements (2)"))
			})
		})

		It("should compare the selection elements for equality", func() {
			selection.EqualsElement(otherSelection)
			Expect(element.IsEqualToCall.Element).To(Equal(otherElement))
		})

		Context("if the provided element is not a *Selection", func() {
			It("should return an error", func() {
				_, err := selection.EqualsElement("not a selection")
				Expect(err).To(MatchError("provided object is not a selection"))
			})
		})

		Context("if the client fails to compare the elements", func() {
			It("should return an error", func() {
				element.IsEqualToCall.Err = errors.New("some error")
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).To(MatchError("failed to compare 'CSS: #selector' to 'CSS: #other_selector': some error"))
			})
		})

		Context("if the client succeeds in comparing the elements", func() {
			It("should return true if they are equal", func() {
				element.IsEqualToCall.ReturnEquals = true
				equal, _ := selection.EqualsElement(otherSelection)
				Expect(equal).To(BeTrue())
			})

			It("should return false if they are not equal", func() {
				element.IsEqualToCall.ReturnEquals = false
				equal, _ := selection.EqualsElement(otherSelection)
				Expect(equal).To(BeFalse())
			})

			It("should not return an error", func() {
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
