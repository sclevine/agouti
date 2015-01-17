package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
)

var _ = Describe("Utils", func() {
	var (
		selection         *Selection
		client            *mocks.Client
		elementRepository *mocks.ElementRepository
		element           *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		elementRepository = &mocks.ElementRepository{}
		emptySelection := &Selection{Client: client, Elements: elementRepository}
		selection = emptySelection.AppendCSS("#selector")
		element = &mocks.Element{}
	})

	Describe("#Count", func() {
		BeforeEach(func() {
			elementRepository.GetCall.ReturnElements = []Element{element, element}
		})

		It("should request elements from the client using the provided selector", func() {
			selection.Count()
			Expect(elementRepository.GetCall.Selectors).To(Equal([]Selector{Selector{Type: "css selector", Value: "#selector"}}))
		})

		Context("when the client succeeds in retrieving the elements", func() {
			It("should successfully return the text", func() {
				count, err := selection.Count()
				Expect(count).To(Equal(2))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the the client fails to retrieve the elements", func() {
			It("should return an error", func() {
				elementRepository.GetCall.Err = errors.New("some error")
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#EqualsElement", func() {
		var (
			otherElementRepository *mocks.ElementRepository
			otherSelection         *Selection
			otherElement           *api.Element
		)

		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = element
			otherElementRepository = &mocks.ElementRepository{}
			emptySelection := &Selection{Elements: otherElementRepository}
			otherSelection = emptySelection.AppendCSS("#other_selector")
			otherElement = &api.Element{}
			otherElementRepository.GetExactlyOneCall.ReturnElement = otherElement
		})

		It("should compare the selection elements for equality", func() {
			selection.EqualsElement(otherSelection)
			Expect(element.IsEqualToCall.Element).To(Equal(otherElement))
		})

		It("should successfully return true if they are equal", func() {
			element.IsEqualToCall.ReturnEquals = true
			equal, err := selection.EqualsElement(otherSelection)
			Expect(equal).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should successfully return false if they are not equal", func() {
			element.IsEqualToCall.ReturnEquals = false
			equal, err := selection.EqualsElement(otherSelection)
			Expect(equal).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is an error retrieving elements from the selection", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})

		Context("when there is an error retrieving elements from the other selection", func() {
			It("should return an error", func() {
				otherElementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).To(MatchError("failed to select 'CSS: #other_selector': some error"))
			})
		})

		Context("when the client fails to compare the elements", func() {
			It("should return an error", func() {
				element.IsEqualToCall.Err = errors.New("some error")
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).To(MatchError("failed to compare 'CSS: #selector' to 'CSS: #other_selector': some error"))
			})
		})
	})

	Describe("#SwitchToFrame", func() {
		var apiElement *api.Element

		BeforeEach(func() {
			apiElement = &api.Element{}
			elementRepository.GetExactlyOneCall.ReturnElement = apiElement
		})

		It("should successfully switch to the frame indicated by the selection", func() {
			Expect(selection.SwitchToFrame()).To(Succeed())
			Expect(client.FrameCall.Frame).To(Equal(apiElement))
		})

		Context("when there is an error retrieving exactly one element", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				err := selection.SwitchToFrame()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})

		Context("when the client fails to switch frames", func() {
			It("should return an error", func() {
				client.FrameCall.Err = errors.New("some error")
				err := selection.SwitchToFrame()
				Expect(err).To(MatchError("failed to switch to frame 'CSS: #selector': some error"))
			})
		})
	})
})
