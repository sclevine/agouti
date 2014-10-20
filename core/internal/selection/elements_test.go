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
		selection types.Selection
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
			err          error
		)

		BeforeEach(func() {
			firstParent = &mocks.Element{}
			secondParent = &mocks.Element{}
			children = []*mocks.Element{&mocks.Element{}, &mocks.Element{}, &mocks.Element{}, &mocks.Element{}}
			firstParent.GetElementsCall.ReturnElements = []types.Element{children[0], children[1]}
			secondParent.GetElementsCall.ReturnElements = []types.Element{children[2], children[3]}
			client.GetElementsCall.ReturnElements = []types.Element{firstParent, secondParent}
		})

		Context("when all elements are successful retrieved", func() {
			BeforeEach(func() {
				err = selection.All("parents").AllByXPath("children").Click()
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
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

		Context("when all indexed elements are successful retrieved", func() {
			BeforeEach(func() {
				err = selection.All("parents").At(1).AllByXPath("children").At(1).Click()
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
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

		Context("when there is no selection", func() {
			It("should return an error", func() {
				err := selection.Click()
				Expect(err).To(MatchError("failed to select '': empty selection"))
			})
		})

		Context("when retrieving the parent elements fails", func() {
			BeforeEach(func() {
				selection = selection.All("parents")
				client.GetElementsCall.Err = errors.New("some error")
			})

			It("should return the error", func() {
				err := selection.Click()
				Expect(err).To(MatchError("failed to select 'CSS: parents': some error"))
			})
		})

		Context("when retrieving any of the child elements fails", func() {
			BeforeEach(func() {
				selection = selection.All("parents").AllByXPath("children")
				secondParent.GetElementsCall.Err = errors.New("some error")
			})

			It("should return the error", func() {
				err := selection.Click()
				Expect(err).To(MatchError("failed to select 'CSS: parents | XPath: children': some error"))
			})
		})

		Context("when the parent selection index is out of range", func() {
			It("should return an error with the index and maximum index", func() {
				selection = selection.All("parents").At(2)
				Expect(selection.Click()).To(MatchError("failed to select 'CSS: parents [2]': element index out of range (>1)"))
			})
		})

		Context("when child selection indices are out of range", func() {
			It("should return an error with the index and maximum index", func() {
				selection = selection.All("parents").At(0).All("children").At(2)
				Expect(selection.Click()).To(MatchError("failed to select 'CSS: parents [0] | CSS: children [2]': element index out of range (>1)"))
			})
		})
	})

	Describe("retrieving at least one element", func() {
		Context("when the client retrieves zero elements", func() {
			It("should fail with an error indicating there were no elements", func() {
				selection = selection.All("#selector")
				client.GetElementsCall.ReturnElements = []types.Element{}
				Expect(selection.Click()).To(MatchError("failed to select 'CSS: #selector': no elements found"))
			})
		})
	})

	Describe("retrieving exactly one element", func() {
		Context("when the client retrieves zero elements", func() {
			It("should fail with an error indicating there are no elements", func() {
				selection = selection.All("#selector")
				client.GetElementsCall.ReturnElements = []types.Element{}
				_, err := selection.Text()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': no elements found"))
			})
		})

		Context("when the client retrieves more than one element", func() {
			It("should fail with an error indicating there are too many elements", func() {
				selection = selection.All("#selector")
				client.GetElementsCall.ReturnElements = []types.Element{&mocks.Element{}, &mocks.Element{}}
				_, err := selection.Text()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': method does not support multiple elements (2)"))
			})
		})
	})
})
