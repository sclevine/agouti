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
				Expect(err).To(BeNil())
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
				Expect(err).To(BeNil())
			})
		})
	})
})
