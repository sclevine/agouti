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
		selection types.Selection
		driver    *mocks.Driver
		element   *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		element = &mocks.Element{}
		selection = &Selection{Driver: driver}
		selection = selection.Find("#selector")
	})

	ItShouldEnsureASingleElement := func(matcher func() error) {
		Context("ensures a single element is returned", func() {
			It("returns an error with the number of elements", func() {
				driver.GetElementsCall.ReturnElements = []types.Element{element, element}
				Expect(matcher()).To(MatchError("failed to retrieve element with 'CSS: #selector': mutiple elements (2) were selected"))
			})
		})
	}

	Describe("#Text", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
		})

		Context("if the the driver fails to retrieve the element text", func() {
			BeforeEach(func() {
				element.GetTextCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				_, err := selection.Text()
				Expect(err).To(MatchError("failed to retrieve text for 'CSS: #selector': some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element text", func() {
			BeforeEach(func() {
				element.GetTextCall.ReturnText = "some text"
			})

			It("returns the text", func() {
				text, _ := selection.Text()
				Expect(text).To(Equal("some text"))
			})

			It("does not return an error", func() {
				_, err := selection.Text()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Attribute", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
		})

		It("requests the attribute value using the attribute name", func() {
			selection.Attribute("some-attribute")
			Expect(element.GetAttributeCall.Attribute).To(Equal("some-attribute"))
		})

		Context("if the the driver fails to retrieve the requested element attribute", func() {
			It("returns an error", func() {
				element.GetAttributeCall.Err = errors.New("some error")
				_, err := selection.Attribute("some-attribute")
				Expect(err).To(MatchError("failed to retrieve attribute value for 'CSS: #selector': some error"))
			})
		})

		Context("if the driver succeeds in retrieving the requested element attribute", func() {
			BeforeEach(func() {
				element.GetAttributeCall.ReturnValue = "some value"
			})

			It("returns the attribute value", func() {
				value, _ := selection.Attribute("some-attribute")
				Expect(value).To(Equal("some value"))
			})

			It("does not return an error", func() {
				_, err := selection.Attribute("some-attribute")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#CSS", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
		})

		It("requests the CSS property value using the property name", func() {
			selection.CSS("some-property")
			Expect(element.GetCSSCall.Property).To(Equal("some-property"))
		})

		Context("if the the driver fails to retrieve the requested element CSS property", func() {
			It("returns an error", func() {
				element.GetCSSCall.Err = errors.New("some error")
				_, err := selection.CSS("some-property")
				Expect(err).To(MatchError("failed to retrieve CSS property for 'CSS: #selector': some error"))
			})
		})

		Context("if the driver succeeds in retrieving the requested element CSS property", func() {
			BeforeEach(func() {
				element.GetCSSCall.ReturnValue = "some value"
			})

			It("returns the property value", func() {
				value, _ := selection.CSS("some-property")
				Expect(value).To(Equal("some value"))
			})

			It("does not return an error", func() {
				_, err := selection.CSS("some-property")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Selected", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Selected()
			return err
		})

		Context("if the the driver fails to retrieve the element's selected status", func() {
			It("returns an error", func() {
				element.IsSelectedCall.Err = errors.New("some error")
				_, err := selection.Selected()
				Expect(err).To(MatchError("failed to determine whether 'CSS: #selector' is selected: some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element's selected status", func() {
			It("returns the selected status when selected", func() {
				element.IsSelectedCall.ReturnSelected = true
				value, _ := selection.Selected()
				Expect(value).To(BeTrue())
			})

			It("returns the selected status when not selected", func() {
				element.IsSelectedCall.ReturnSelected = false
				value, _ := selection.Selected()
				Expect(value).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := selection.Selected()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Visible", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Visible()
			return err
		})

		Context("if the the driver fails to retrieve the element's visible status", func() {
			It("returns an error", func() {
				element.IsDisplayedCall.Err = errors.New("some error")
				_, err := selection.Visible()
				Expect(err).To(MatchError("failed to determine whether 'CSS: #selector' is visible: some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element's visible status", func() {
			It("returns the visible status when visible", func() {
				element.IsDisplayedCall.ReturnDisplayed = true
				value, _ := selection.Visible()
				Expect(value).To(BeTrue())
			})

			It("returns the visible status when not visible", func() {
				element.IsDisplayedCall.ReturnDisplayed = false
				value, _ := selection.Visible()
				Expect(value).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := selection.Visible()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Enabled", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Enabled()
			return err
		})

		Context("if the the driver fails to retrieve the element's enabled status", func() {
			It("returns an error", func() {
				element.IsEnabledCall.Err = errors.New("some error")
				_, err := selection.Enabled()
				Expect(err).To(MatchError("failed to determine whether 'CSS: #selector' is enabled: some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element's enabled status", func() {
			It("returns the enabled status when enabled", func() {
				element.IsEnabledCall.ReturnEnabled = true
				value, _ := selection.Enabled()
				Expect(value).To(BeTrue())
			})

			It("returns the enabled status when not enabled", func() {
				element.IsEnabledCall.ReturnEnabled = false
				value, _ := selection.Enabled()
				Expect(value).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := selection.Enabled()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
