package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
)

var _ = Describe("Elements", func() {
	var (
		client    *mocks.Client
		selection *Selection
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		selection = &Selection{Client: client}
	})

	Describe("retrieving any number of elements", func() {
		var (
			firstParent  *mocks.Element
			secondParent *mocks.Element
			children     []*mocks.Element
		)

		BeforeEach(func() {
			firstParent = &mocks.Element{}
			secondParent = &mocks.Element{}
			children = []*mocks.Element{&mocks.Element{}, &mocks.Element{}, &mocks.Element{}, &mocks.Element{}}
			firstParent.GetElementsCall.ReturnElements = []types.Element{children[0], children[1]}
			secondParent.GetElementsCall.ReturnElements = []types.Element{children[2], children[3]}
			client.GetElementsCall.ReturnElements = []types.Element{firstParent, secondParent}
		})

		Context("when all elements are successfully retrieved", func() {
			BeforeEach(func() {
				Expect(selection.All("parents").AllByXPath("children").Click()).To(Succeed())
			})

			It("should retrieve the parent elements using the client", func() {
				Expect(client.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "parents"}))
			})

			It("should retrieve the child elements of the parent selector", func() {
				Expect(firstParent.GetElementsCall.Selector).To(Equal(types.Selector{Using: "xpath", Value: "children"}))
				Expect(secondParent.GetElementsCall.Selector).To(Equal(types.Selector{Using: "xpath", Value: "children"}))
			})

			It("should click on all of the child elements", func() {
				Expect(children[0].ClickCall.Called).To(BeTrue())
				Expect(children[1].ClickCall.Called).To(BeTrue())
				Expect(children[2].ClickCall.Called).To(BeTrue())
				Expect(children[3].ClickCall.Called).To(BeTrue())
			})
		})

		Context("when all non-zero-indexed elements are successfully retrieved", func() {
			BeforeEach(func() {
				Expect(selection.All("parents").At(1).AllByXPath("children").At(1).Click()).To(Succeed())
			})

			It("should retrieve the parent elements using the client", func() {
				Expect(client.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "parents", Index: 1, Indexed: true}))
			})

			It("should retrieve the child elements of the parent selector", func() {
				Expect(firstParent.GetElementsCall.Selector.Using).To(BeEmpty())
				Expect(secondParent.GetElementsCall.Selector).To(Equal(types.Selector{Using: "xpath", Value: "children", Index: 1, Indexed: true}))
			})

			It("should click on only the selected child elements", func() {
				Expect(children[0].ClickCall.Called).To(BeFalse())
				Expect(children[1].ClickCall.Called).To(BeFalse())
				Expect(children[2].ClickCall.Called).To(BeFalse())
				Expect(children[3].ClickCall.Called).To(BeTrue())
			})
		})

		Context("when all zero-indexed elements are successfully retrieved", func() {
			BeforeEach(func() {
				firstParent.GetElementCall.ReturnElement = children[0]
				client.GetElementCall.ReturnElement = firstParent
				Expect(selection.All("parents").At(0).AllByXPath("children").At(0).Click()).To(Succeed())
			})

			It("should retrieve the first parent element using the client", func() {
				Expect(client.GetElementCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "parents", Index: 0, Indexed: true}))
			})

			It("should retrieve the first child element of the parent selector", func() {
				Expect(firstParent.GetElementCall.Selector).To(Equal(types.Selector{Using: "xpath", Value: "children", Index: 0, Indexed: true}))
			})

			It("should click on only the selected child elements", func() {
				Expect(children[0].ClickCall.Called).To(BeTrue())
				Expect(children[1].ClickCall.Called).To(BeFalse())
				Expect(children[2].ClickCall.Called).To(BeFalse())
				Expect(children[3].ClickCall.Called).To(BeFalse())
			})
		})

		Context("when single-element-only elements are successfully retrieved", func() {
			BeforeEach(func() {
				firstParent.GetElementsCall.ReturnElements = []types.Element{children[0]}
				client.GetElementsCall.ReturnElements = []types.Element{firstParent}
				Expect(selection.All("parents").Single().AllByXPath("children").Single().Click()).To(Succeed())
			})

			It("should retrieve the parent element using the client", func() {
				Expect(client.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "parents", Single: true}))
			})

			It("should retrieve the child element of the parent selector", func() {
				Expect(firstParent.GetElementsCall.Selector).To(Equal(types.Selector{Using: "xpath", Value: "children", Single: true}))
			})

			It("should click on the selected child element", func() {
				Expect(children[0].ClickCall.Called).To(BeTrue())
			})
		})

		Context("when there is no selection", func() {
			It("should return an error", func() {
				Expect(selection.Click()).To(MatchError("failed to select '': empty selection"))
			})
		})

		Context("when retrieving the parent elements fails", func() {
			It("should return an error", func() {
				bothParents := selection.All("parents")
				client.GetElementsCall.Err = errors.New("some error")
				Expect(bothParents.Click()).To(MatchError("failed to select 'CSS: parents': some error"))
			})
		})

		Context("when retrieving any of the child elements fails", func() {
			It("should return an error", func() {
				allChildren := selection.All("parents").AllByXPath("children")
				secondParent.GetElementsCall.Err = errors.New("some error")
				Expect(allChildren.Click()).To(MatchError("failed to select 'CSS: parents | XPath: children': some error"))
			})
		})

		Context("when a single-element-only parent selection refers to multiple parents", func() {
			It("should return an error", func() {
				allChildren := selection.All("parents").Single().AllByXPath("children")
				Expect(allChildren.Click()).To(MatchError("failed to select 'CSS: parents [single] | XPath: children': ambiguous find"))
			})
		})

		Context("when a single-element-only parent selection refers to no parents", func() {
			It("should return an error", func() {
				noChildren := selection.All("parents").Single().AllByXPath("children")
				client.GetElementsCall.ReturnElements = []types.Element{}
				Expect(noChildren.Click()).To(MatchError("failed to select 'CSS: parents [single] | XPath: children': element not found"))
			})
		})

		Context("when any single-element-only child selection refers to multiple child elements", func() {
			It("should return an error", func() {
				allChildren := selection.All("parents").AllByXPath("children").Single()
				firstParent.GetElementsCall.ReturnElements = []types.Element{children[0]}
				Expect(allChildren.Click()).To(MatchError("failed to select 'CSS: parents | XPath: children [single]': ambiguous find"))
			})
		})

		Context("when any single-element-only child selection refers to no child elements", func() {
			It("should return an error", func() {
				noChild := selection.All("parents").AllByXPath("children").Single()
				firstParent.GetElementsCall.ReturnElements = []types.Element{}
				Expect(noChild.Click()).To(MatchError("failed to select 'CSS: parents | XPath: children [single]': element not found"))
			})
		})

		Context("when the parent selection index is out of range", func() {
			It("should return an error", func() {
				noParent := selection.All("parents").At(2)
				Expect(noParent.Click()).To(MatchError("failed to select 'CSS: parents [2]': element index out of range"))
			})
		})

		Context("when child selection indices are out of range", func() {
			It("should return an error", func() {
				noChild := selection.All("parents").At(1).All("children").At(2)
				Expect(noChild.Click()).To(MatchError("failed to select 'CSS: parents [1] | CSS: children [2]': element index out of range"))
			})
		})

		Context("when a zero-indexed parent selection element does not exist", func() {
			It("should return an error", func() {
				client.GetElementCall.Err = errors.New("some error")
				noParent := selection.All("parents").At(0)
				Expect(noParent.Click()).To(MatchError("failed to select 'CSS: parents [0]': some error"))
			})
		})

		Context("when a zero-indexed child selection element does not exist", func() {
			It("should return an error", func() {
				firstParent.GetElementCall.Err = errors.New("some error")
				client.GetElementCall.ReturnElement = firstParent
				noChild := selection.All("parents").At(0).All("children").At(0)
				Expect(noChild.Click()).To(MatchError("failed to select 'CSS: parents [0] | CSS: children [0]': some error"))
			})
		})
	})

	Describe("retrieving at least one element", func() {
		Context("when the client retrieves zero elements", func() {
			It("should fail with an error", func() {
				empty := selection.All("#selector")
				client.GetElementsCall.ReturnElements = []types.Element{}
				Expect(empty.Click()).To(MatchError("failed to select 'CSS: #selector': no elements found"))
			})
		})
	})

	Describe("retrieving exactly one element", func() {
		Context("when the client retrieves zero elements", func() {
			It("should return an error", func() {
				empty := selection.All("#selector")
				client.GetElementsCall.ReturnElements = []types.Element{}
				_, err := empty.Text()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': no elements found"))
			})
		})

		Context("when the client retrieves more than one element", func() {
			It("should return an error", func() {
				multiple := selection.All("#selector")
				client.GetElementsCall.ReturnElements = []types.Element{&mocks.Element{}, &mocks.Element{}}
				_, err := multiple.Text()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': method does not support multiple elements (2)"))
			})
		})
	})
})
