package page_test

import (
	. "github.com/sclevine/agouti/page"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/internal/mocks"
	"github.com/sclevine/agouti/page/internal/webdriver"
)

var _ = Describe("Selection", func() {
	var (
		selection Selection
		driver    *mocks.Driver
		element   *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		element = &mocks.Element{}
		selection = (&Page{driver}).Find("#selector")
	})

	ItShouldEnsureASingleElement := func(matcher func() error) {
		Context("when the driver fails to retrieve any elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.Err = errors.New("some error")
			})

			It("returns error from the driver", func() {
				Expect(matcher()).To(MatchError("failed to retrieve element with selector '#selector': some error"))
			})
		})

		Context("when the driver retrieves more than one element", func() {
			BeforeEach(func() {
				driver.GetElementsCall.ReturnElements = []webdriver.Element{element, element}
			})

			It("returns an error with the number of elements", func() {
				Expect(matcher()).To(MatchError("failed to retrieve element with selector '#selector': mutiple elements (2) were selected"))
			})
		})

		Context("when the driver retrieves zero elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.ReturnElements = []webdriver.Element{}
			})

			It("fails with an error indicating there were no elements", func() {
				Expect(matcher()).To(MatchError("failed to retrieve element with selector '#selector': no element found"))
			})
		})
	}

	Describe("#Find", func() {
		It("returns a subselection", func() {
			Expect(selection.Find("#subselector").Selector()).To(Equal("#selector #subselector"))
		})
	})

	Describe("#Selector", func() {
		It("returns the selector", func() {
			Expect(selection.Selector()).To(Equal("#selector"))
		})
	})

	Describe("#Click", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.Click()
		})

		Context("if the click fails", func() {
			BeforeEach(func() {
				element.ClickCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Click()).To(MatchError("failed to click on selector '#selector': some error"))
			})
		})

		Context("if the click succeeds", func() {
			It("clicks on an element", func() {
				selection.Click()
				Expect(element.ClickCall.Called).To(BeTrue())
			})

			It("returns nil", func() {
				Expect(selection.Click()).To(BeNil())
			})
		})
	})

	Describe("#Fill", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.Fill("some text")
		})

		Context("if clearing the element fails", func() {
			BeforeEach(func() {
				element.ClearCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Fill("some text")).To(MatchError("failed to clear selector '#selector': some error"))
			})
		})

		Context("if entering text into the element fails", func() {
			BeforeEach(func() {
				element.ValueCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Fill("some text")).To(MatchError("failed to enter text into selector '#selector': some error"))
			})
		})

		Context("if the fill succeeds", func() {
			It("clears the element", func() {
				selection.Fill("some text")
				Expect(element.ClearCall.Called).To(BeTrue())
			})

			It("fills the element with the provided text", func() {
				selection.Fill("some text")
				Expect(element.ValueCall.Text).To(Equal("some text"))
			})

			It("returns nil", func() {
				Expect(selection.Fill("some text")).To(BeNil())
			})
		})
	})

	Describe("#Check", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.Check()
		})

		It("checks the type of the checkbox", func() {
			selection.Check()
			Expect(element.GetAttributeCall.Attribute).To(Equal("type"))
		})

		Context("when the the driver fails to retrieve the 'type' attribute", func() {
			BeforeEach(func() {
				element.GetAttributeCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Check()).To(MatchError("failed to retrieve type of selector '#selector': some error"))
			})
		})

		Context("when the selection is not a checkbox", func() {
			BeforeEach(func() {
				element.GetAttributeCall.ReturnValue = "banana"
			})

			It("returns an error", func() {
				Expect(selection.Check()).To(MatchError("selector '#selector' does not refer to a checkbox"))
			})
		})

		Context("when the selection is a checkbox", func() {
			BeforeEach(func() {
				element.GetAttributeCall.ReturnValue = "checkbox"
			})

			Context("when the determining the selected status of the element fails", func() {
				BeforeEach(func() {
					element.SelectedCall.Err = errors.New("some error")
				})

				It("returns an error", func() {
					Expect(selection.Check()).To(MatchError("failed to retrieve selected state of selector '#selector': some error"))
				})
			})

			Context("when the box is already checked", func() {
				BeforeEach(func() {
					element.SelectedCall.ReturnSelected = true
				})

				It("does not click on the checkbox", func() {
					selection.Check()
					Expect(element.ClickCall.Called).To(BeFalse())
				})
			})

			Context("when the box is not checked", func() {
				BeforeEach(func() {
					element.SelectedCall.ReturnSelected = false
				})

				It("clicks on the checkbox", func() {
					selection.Check()
					Expect(element.ClickCall.Called).To(BeTrue())
				})

				Context("when clicking on the checkbox fails", func() {
					BeforeEach(func() {
						element.ClickCall.Err = errors.New("some error")
					})

					It("returns an error", func() {
						Expect(selection.Check()).To(MatchError("failed to check selector '#selector': some error"))
					})
				})
			})
		})
	})

	Describe("#Text", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
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
				Expect(err).To(MatchError("failed to retrieve text for selector '#selector': some error"))
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
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
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
				Expect(err).To(MatchError("failed to retrieve attribute value for selector '#selector': some error"))
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
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
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
				Expect(err).To(MatchError("failed to retrieve CSS property for selector '#selector': some error"))
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
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Selected()
			return err
		})

		Context("if the the driver fails to retrieve the element's selected status", func() {
			It("returns an error", func() {
				element.SelectedCall.Err = errors.New("some error")
				_, err := selection.Selected()
				Expect(err).To(MatchError("failed to determine whether selector '#selector' is selected: some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element's selected status", func() {
			It("returns the selected status when selected", func() {
				element.SelectedCall.ReturnSelected = true
				value, _ := selection.Selected()
				Expect(value).To(BeTrue())
			})

			It("returns the selected status when not selected", func() {
				element.SelectedCall.ReturnSelected = false
				value, _ := selection.Selected()
				Expect(value).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := selection.Selected()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
