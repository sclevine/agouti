package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
)

var _ = Describe("Selection", func() {
	var (
		selection         *Selection
		client            *mocks.Client
		elementRepository *mocks.ElementRepository
		element           *mocks.Element
		secondElement     *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		element = &mocks.Element{}
		secondElement = &mocks.Element{}
		elementRepository = &mocks.ElementRepository{}
		emptySelection := &Selection{Client: client, Elements: elementRepository}
		selection = emptySelection.AppendCSS("#selector")
	})

	Describe("#Text", func() {
		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = element
		})

		Context("when the client fails to retrieve the element text", func() {
			It("should return an error", func() {
				element.GetTextCall.Err = errors.New("some error")
				_, err := selection.Text()
				Expect(err).To(MatchError("failed to retrieve text for 'CSS: #selector': some error"))
			})
		})

		Context("when the client succeeds in retrieving the element text", func() {
			It("should successfully return the text", func() {
				element.GetTextCall.ReturnText = "some text"
				text, err := selection.Text()
				Expect(text).To(Equal("some text"))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Active", func() {
		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = element
		})

		Context("when the client fails to retrieve the active element", func() {
			It("should return an error", func() {
				client.GetActiveElementCall.Err = errors.New("some error")
				_, err := selection.Active()
				Expect(err).To(MatchError("failed to retrieve active element: some error"))
			})
		})

		It("should compare the active and selected elements", func() {
			activeElement := &api.Element{}
			client.GetActiveElementCall.ReturnElement = activeElement
			selection.Active()
			Expect(element.IsEqualToCall.Element).To(Equal(activeElement))
		})

		Context("when the client fails to compare active element to the selected element", func() {
			It("should return an error", func() {
				element.IsEqualToCall.Err = errors.New("some error")
				_, err := selection.Active()
				Expect(err).To(MatchError("failed to compare selection to active element: some error"))
			})
		})

		Context("when the active element equals the selected element", func() {
			It("should successfully return true", func() {
				element.IsEqualToCall.ReturnEquals = true
				equal, err := selection.Active()
				Expect(equal).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the active element does not equal the selected element", func() {
			It("should successfully return false", func() {
				element.IsEqualToCall.ReturnEquals = false
				equal, err := selection.Active()
				Expect(equal).To(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Attribute", func() {
		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = element
		})

		It("should request the attribute value using the attribute name", func() {
			selection.Attribute("some-attribute")
			Expect(element.GetAttributeCall.Attribute).To(Equal("some-attribute"))
		})

		Context("when the client fails to retrieve the requested element attribute", func() {
			It("should return an error", func() {
				element.GetAttributeCall.Err = errors.New("some error")
				_, err := selection.Attribute("some-attribute")
				Expect(err).To(MatchError("failed to retrieve attribute value for 'CSS: #selector': some error"))
			})
		})

		Context("when the client succeeds in retrieving the requested element attribute", func() {
			It("should successfully return the attribute value", func() {
				element.GetAttributeCall.ReturnValue = "some value"
				value, err := selection.Attribute("some-attribute")
				Expect(value).To(Equal("some value"))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#CSS", func() {
		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = element
		})

		It("should successfully request the CSS property value using the property name", func() {
			selection.CSS("some-property")
			Expect(element.GetCSSCall.Property).To(Equal("some-property"))
		})

		Context("when the the client fails to retrieve the requested element CSS property", func() {
			It("should return an error", func() {
				element.GetCSSCall.Err = errors.New("some error")
				_, err := selection.CSS("some-property")
				Expect(err).To(MatchError("failed to retrieve CSS property value for 'CSS: #selector': some error"))
			})
		})

		Context("when the client succeeds in retrieving the requested element CSS property", func() {
			It("should successfully return the property value", func() {
				element.GetCSSCall.ReturnValue = "some value"
				value, err := selection.CSS("some-property")
				Expect(value).To(Equal("some value"))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Selected", func() {
		BeforeEach(func() {
			elementRepository.GetAtLeastOneCall.ReturnElements = []Element{element, secondElement}
		})

		Context("when the the client fails to retrieve any elements' selected status", func() {
			It("should return an error", func() {
				element.IsSelectedCall.ReturnSelected = true
				secondElement.IsSelectedCall.Err = errors.New("some error")
				_, err := selection.Selected()
				Expect(err).To(MatchError("failed to determine whether some 'CSS: #selector' is selected: some error"))
			})
		})

		Context("when the client succeeds in retrieving all elements' selected status", func() {
			It("should return true when all elements are selected", func() {
				element.IsSelectedCall.ReturnSelected = true
				secondElement.IsSelectedCall.ReturnSelected = true
				value, err := selection.Selected()
				Expect(value).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return false when any elements are not selected", func() {
				element.IsSelectedCall.ReturnSelected = true
				secondElement.IsSelectedCall.ReturnSelected = false
				value, err := selection.Selected()
				Expect(value).To(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Visible", func() {
		BeforeEach(func() {
			elementRepository.GetAtLeastOneCall.ReturnElements = []Element{element, secondElement}
		})

		Context("when the the client fails to retrieve any elements' visible status", func() {
			It("should return an error", func() {
				element.IsDisplayedCall.ReturnDisplayed = true
				secondElement.IsDisplayedCall.Err = errors.New("some error")
				_, err := selection.Visible()
				Expect(err).To(MatchError("failed to determine whether some 'CSS: #selector' is visible: some error"))
			})
		})

		Context("when the client succeeds in retrieving all elements' visible status", func() {
			It("should return true when all elements are visible", func() {
				element.IsDisplayedCall.ReturnDisplayed = true
				secondElement.IsDisplayedCall.ReturnDisplayed = true
				value, err := selection.Visible()
				Expect(value).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return false when any elements are not visible", func() {
				element.IsDisplayedCall.ReturnDisplayed = true
				secondElement.IsDisplayedCall.ReturnDisplayed = false
				value, err := selection.Visible()
				Expect(value).To(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Enabled", func() {
		BeforeEach(func() {
			elementRepository.GetAtLeastOneCall.ReturnElements = []Element{element, secondElement}
		})

		Context("when the the client fails to retrieve any elements' enabled status", func() {
			It("should return an error", func() {
				element.IsEnabledCall.ReturnEnabled = true
				secondElement.IsEnabledCall.Err = errors.New("some error")
				_, err := selection.Enabled()
				Expect(err).To(MatchError("failed to determine whether some 'CSS: #selector' is enabled: some error"))
			})
		})

		Context("when the client succeeds in retrieving all elements' enabled status", func() {
			It("should return true when all elements are enabled", func() {
				element.IsEnabledCall.ReturnEnabled = true
				secondElement.IsEnabledCall.ReturnEnabled = true
				value, err := selection.Enabled()
				Expect(value).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return false when any elements are not enabled", func() {
				element.IsEnabledCall.ReturnEnabled = true
				secondElement.IsEnabledCall.ReturnEnabled = false
				value, err := selection.Enabled()
				Expect(value).To(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
