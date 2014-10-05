package page_test

import (
	. "github.com/sclevine/agouti/page"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/mocks"
	"github.com/sclevine/agouti/webdriver"
)

var _ = Describe("Selection", func() {
	var (
		selection Selection
		failer    *mocks.Failer
		driver    *mocks.Driver
		element   *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		failer = &mocks.Failer{}
		element = &mocks.Element{}
		selection = (&Page{driver, failer.Fail}).Within("#selector")
	})

	Describe("#Within", func() {
		It("returns a subselection", func() {
			subselection := selection.Within("#subselector")
			Expect(subselection.Selector()).To(Equal("#selector #subselector"))
		})

		It("provides a subselection to any specified callable body functions", func() {
			selection.Within("#subselector",
				Do(func(subselection Selection) {
					Expect(subselection.Selector()).To(Equal("#selector #subselector"))
				}),
				Do(func(subselection Selection) {
					Expect(subselection.Selector()).To(Equal("#selector #subselector"))
				}),
			)
		})
	})

	Describe("#Selector", func() {
		It("returns the selector", func() {
			Expect(selection.Selector()).To(Equal("#selector"))
		})
	})

	Describe("#ShouldContainText", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
			element.GetTextCall.ReturnText = "element text"
		})

		Context("when the driver fails to retrieve any elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.Err = errors.New("some error")
			})

			It("fails with an error", func() {
				selection.ShouldContainText("text")
				Expect(failer.Message).To(Equal("Failed to retrieve elements: some error"))
			})

			It("fails with an offset of one", func() {
				selection.ShouldContainText("text")
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})

		Context("when the driver retrieves more than one element", func() {
			BeforeEach(func() {
				driver.GetElementsCall.ReturnElements = []webdriver.Element{element, element}
			})

			It("fails with the number of elements", func() {
				selection.ShouldContainText("text")
				Expect(failer.Message).To(Equal("Mutiple elements (2) were selected."))
			})

			It("fails with an offset of one", func() {
				selection.ShouldContainText("text")
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})

		Context("when the driver retrieves zero elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.ReturnElements = []webdriver.Element{}
			})

			It("fails with an error indicating there were no elements", func() {
				selection.ShouldContainText("text")
				Expect(failer.Message).To(Equal("No elements found."))
			})

			It("fails with an offset of one", func() {
				selection.ShouldContainText("text")
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})

		Context("when the driver cannot retrieve an element's text", func() {
			BeforeEach(func() {
				element.GetTextCall.Err = errors.New("some error")
			})

			It("fails with the selector and an error", func() {
				selection.ShouldContainText("text")
				Expect(failer.Message).To(Equal("Failed to retrieve text for selector '#selector': some error"))
			})

			It("fails with an offset of one", func() {
				selection.ShouldContainText("text")
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})

		Context("when the a single element text is found", func() {
			Context("if the provided text is a substring of the element text", func() {
				It("it does not fail the test", func() {
					selection.ShouldContainText("ment tex")
					Expect(failer.Failed).To(BeFalse())
				})
			})

			Context("if the provided text is not a substring of the element text", func() {
				It("fails with information about the failure", func() {
					selection.ShouldContainText("banana")
					Expect(failer.Message).To(Equal("Failed to find text 'banana' for selector '#selector'.\nFound: 'element text'"))
				})

				It("fails with an offset of 1", func() {
					selection.ShouldContainText("banana")
					Expect(failer.CallerSkip).To(Equal(1))
				})
			})
		})
	})
})
