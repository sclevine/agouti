package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
)

var _ = Describe("Utils", func() {
	var (
		selection *Selection
		client    *mocks.Client
		element   *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		emptySelection := &Selection{Client: client}
		selection = emptySelection.AppendCSS("#selector")
		element = &mocks.Element{}
	})

	Describe("#Count", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element, element}
		})

		It("should request elements from the client using the provided selector", func() {
			selection.Count()
			Expect(client.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "#selector"}))
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
				client.GetElementsCall.Err = errors.New("some error")
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#EqualsElement", func() {
		var (
			otherClient    *mocks.Client
			otherSelection *Selection
			otherElement   *mocks.Element
		)

		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
			otherClient = &mocks.Client{}
			emptySelection := &Selection{Client: otherClient}
			otherSelection = emptySelection.AppendCSS("#other_selector")
			otherElement = &mocks.Element{}
			otherClient.GetElementsCall.ReturnElements = []types.Element{otherElement}
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

		Context("when the client fails to compare the elements", func() {
			It("should return an error", func() {
				element.IsEqualToCall.Err = errors.New("some error")
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).To(MatchError("failed to compare 'CSS: #selector' to 'CSS: #other_selector': some error"))
			})
		})
	})

	Describe("#SwitchToFrame", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		It("should successfully switch to the frame indicated by the selection", func() {
			Expect(selection.SwitchToFrame()).To(Succeed())
			Expect(client.FrameCall.Frame).To(Equal(element))
		})

		Context("when multiple elements are selected", func() {
			It("should return an error with the number of elements", func() {
				client.GetElementsCall.ReturnElements = []types.Element{element, element}
				err := selection.SwitchToFrame()
				Expect(err).To(MatchError("failed to select 'CSS: #selector': method does not support multiple elements (2)"))
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
