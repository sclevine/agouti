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
		selection     types.Selection
		client        *mocks.Client
		element       *mocks.Element
		secondElement *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		element = &mocks.Element{}
		secondElement = &mocks.Element{}
		selection = &Selection{Client: client}
		selection = selection.All("#selector")
	})

	ItShouldEnsureASingleElement := func(matcher func() error) {
		Context("when multiple elements are returned", func() {
			It("should return an error with the number of elements", func() {
				client.GetElementsCall.ReturnElements = []types.Element{element, secondElement}
				ExpectWithOffset(1, matcher()).To(MatchError("failed to select 'CSS: #selector': method does not support multiple elements (2)"))
			})
		})
	}

	ItShouldEnsureAtLeastOneElement := func(matcher func() error) {
		Context("when zero elements are returned", func() {
			It("should return an error with the number of elements", func() {
				client.GetElementsCall.ReturnElements = []types.Element{}
				ExpectWithOffset(1, matcher()).To(MatchError("failed to select 'CSS: #selector': no elements found"))
			})
		})
	}

	Describe("#Text", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
		})

		Context("when the client fails to retrieve the element text", func() {
			It("should return an error", func() {
				element.GetTextCall.Err = errors.New("some error")
				_, err := selection.Text()
				Expect(err).To(MatchError("failed to retrieve text for 'CSS: #selector': some error"))
			})
		})

		Context("when the client succeeds in retrieving the element text", func() {
			BeforeEach(func() {
				element.GetTextCall.ReturnText = "some text"
			})

			It("should return the text", func() {
				text, _ := selection.Text()
				Expect(text).To(Equal("some text"))
			})

			It("should not return an error", func() {
				_, err := selection.Text()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Attribute", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
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
			BeforeEach(func() {
				element.GetAttributeCall.ReturnValue = "some value"
			})

			It("should return the attribute value", func() {
				value, _ := selection.Attribute("some-attribute")
				Expect(value).To(Equal("some value"))
			})

			It("should not return an error", func() {
				_, err := selection.Attribute("some-attribute")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#CSS", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
		})

		It("should request the CSS property value using the property name", func() {
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
			BeforeEach(func() {
				element.GetCSSCall.ReturnValue = "some value"
			})

			It("should return the property value", func() {
				value, _ := selection.CSS("some-property")
				Expect(value).To(Equal("some value"))
			})

			It("should not return an error", func() {
				_, err := selection.CSS("some-property")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Selected", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			_, err := selection.Selected()
			return err
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
			client.GetElementsCall.ReturnElements = []types.Element{element, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			_, err := selection.Visible()
			return err
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
			client.GetElementsCall.ReturnElements = []types.Element{element, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			_, err := selection.Enabled()
			return err
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
