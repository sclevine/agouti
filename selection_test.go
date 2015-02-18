package agouti_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/mocks"
	"github.com/sclevine/agouti/internal/target"
)

var _ = Describe("Selection", func() {
	var (
		firstElement  *mocks.Element
		secondElement *api.Element
	)

	BeforeEach(func() {
		firstElement = &mocks.Element{}
		secondElement = &api.Element{}
	})

	Describe("#String", func() {
		It("should return a string representation of the selection", func() {
			selection := NewTestMultiSelection(nil, nil, "#selector")
			Expect(selection.AllByXPath("#subselector").String()).To(Equal("CSS: #selector | XPath: #subselector"))
		})
	})

	Describe("#Count", func() {
		var (
			selection         *MultiSelection
			elementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			elementRepository = &mocks.ElementRepository{}
			selection = NewTestMultiSelection(elementRepository, nil, "#selector")
			elementRepository.GetCall.ReturnElements = []element.Element{firstElement, secondElement}
		})

		It("should request elements from the session using the provided selector", func() {
			selection.Count()
			Expect(elementRepository.GetCall.Selectors).To(Equal(target.Selectors{target.Selector{Type: "css selector", Value: "#selector"}}))
		})

		Context("when the session succeeds in retrieving the elements", func() {
			It("should successfully return the text", func() {
				Expect(selection.Count()).To(Equal(2))
			})
		})

		Context("when the the session fails to retrieve the elements", func() {
			It("should return an error", func() {
				elementRepository.GetCall.Err = errors.New("some error")
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#EqualsElement", func() {
		var (
			firstSelection          *Selection
			secondSelection         *Selection
			firstElementRepository  *mocks.ElementRepository
			secondElementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			firstElementRepository = &mocks.ElementRepository{}
			firstElementRepository.GetExactlyOneCall.ReturnElement = firstElement
			firstSelection = NewTestSelection(firstElementRepository, nil, "#first_selector")

			secondElementRepository = &mocks.ElementRepository{}
			secondElementRepository.GetExactlyOneCall.ReturnElement = secondElement
			secondSelection = NewTestSelection(secondElementRepository, nil, "#second_selector")
		})

		It("should compare the selection elements for equality", func() {
			firstSelection.EqualsElement(secondSelection)
			Expect(firstElement.IsEqualToCall.Element).To(Equal(secondElement))
		})

		It("should successfully return true if they are equal", func() {
			firstElement.IsEqualToCall.ReturnEquals = true
			Expect(firstSelection.EqualsElement(secondSelection)).To(BeTrue())
		})

		It("should successfully return false if they are not equal", func() {
			firstElement.IsEqualToCall.ReturnEquals = false
			Expect(firstSelection.EqualsElement(secondSelection)).To(BeFalse())
		})

		Context("when the provided object is a *MultiSelection", func() {
			It("should not fail", func() {
				multiSelection := NewTestMultiSelection(secondElementRepository, nil, "#multi_selector")
				Expect(firstSelection.EqualsElement(multiSelection)).To(BeFalse())
				Expect(firstElement.IsEqualToCall.Element).To(Equal(secondElement))
			})
		})

		Context("when the provided object is not a type of selection", func() {
			It("should return an error", func() {
				_, err := firstSelection.EqualsElement("not a selection")
				Expect(err).To(MatchError("must be *Selection or *MultiSelection"))
			})
		})

		Context("when there is an error retrieving elements from the selection", func() {
			It("should return an error", func() {
				firstElementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to select 'CSS: #first_selector [single]': some error"))
			})
		})

		Context("when there is an error retrieving elements from the other selection", func() {
			It("should return an error", func() {
				secondElementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to select 'CSS: #second_selector [single]': some error"))
			})
		})

		Context("when the session fails to compare the elements", func() {
			It("should return an error", func() {
				firstElement.IsEqualToCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to compare 'CSS: #first_selector [single]' to 'CSS: #second_selector [single]': some error"))
			})
		})
	})
})
