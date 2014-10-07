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
		selection = NewPage(driver, failer.Fail).Within("#selector")
	})

	ItShouldRetrieveASingleElement := func(matcherCall func()) {
		Context("when the driver fails to retrieve any elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.Err = errors.New("some error")
			})

			It("fails with an error", func() {
				Expect(matcherCall).To(Panic())
				Expect(failer.Message).To(Equal("Failed to retrieve element with selector '#selector': some error"))
			})

			It("fails with an offset of one", func() {
				Expect(matcherCall).To(Panic())
				Expect(failer.CallerSkip).To(Equal(2))
			})
		})

		Context("when the driver retrieves more than one element", func() {
			BeforeEach(func() {
				driver.GetElementsCall.ReturnElements = []webdriver.Element{element, element}
			})

			It("fails with the number of elements", func() {
				Expect(matcherCall).To(Panic())
				Expect(failer.Message).To(Equal("Mutiple elements (2) with selector '#selector' were selected."))
			})

			It("fails with an offset of one", func() {
				Expect(matcherCall).To(Panic())
				Expect(failer.CallerSkip).To(Equal(2))
			})
		})

		Context("when the driver retrieves zero elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.ReturnElements = []webdriver.Element{}
			})

			It("fails with an error indicating there were no elements", func() {
				Expect(matcherCall).To(Panic())
				Expect(failer.Message).To(Equal("No element with selector '#selector' found."))
			})

			It("fails with an offset of one", func() {
				Expect(matcherCall).To(Panic())
				Expect(failer.CallerSkip).To(Equal(2))
			})
		})
	}

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

		ItShouldRetrieveASingleElement(func() {
			selection.ShouldContainText("text")
		})

		Context("when the driver cannot retrieve an element's text", func() {
			BeforeEach(func() {
				element.GetTextCall.Err = errors.New("some error")
			})

			It("fails with the selector and an error", func() {
				Expect(func() { selection.ShouldContainText("text") }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to retrieve text for selector '#selector': some error"))
			})

			It("fails with an offset of one", func() {
				Expect(func() { selection.ShouldContainText("text") }).To(Panic())
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})

		Context("when the a single element text is found", func() {
			Context("if the provided text is a substring of the element text", func() {
				It("it does not fail the test", func() {
					Expect(func() { selection.ShouldContainText("ment tex") }).NotTo(Panic())
				})
			})

			Context("if the provided text is not a substring of the element text", func() {
				It("fails with information about the failure", func() {
					Expect(func() { selection.ShouldContainText("banana") }).To(Panic())
					Expect(failer.Message).To(Equal("Failed to find text 'banana' for selector '#selector'.\nFound: 'element text'"))
				})

				It("fails with an offset of 1", func() {
					Expect(func() { selection.ShouldContainText("banana") }).To(Panic())
					Expect(failer.CallerSkip).To(Equal(1))
				})
			})
		})
	})

	Describe("#Click", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
		})

		ItShouldRetrieveASingleElement(func() {
			selection.Click()
		})

		Context("if the click fails", func() {
			BeforeEach(func() {
				element.ClickCall.Err = errors.New("some error")
			})

			It("fails with information about the failure", func() {
				Expect(func() { selection.Click() }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to click on selector '#selector': some error"))
			})

			It("fails with an offset of 1", func() {
				Expect(func() { selection.Click() }).To(Panic())
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})

		Context("if the click succeeds", func() {
			It("clicks on an element", func() {
				element.Click()
				Expect(element.ClickCall.Called).To(BeTrue())
			})

			It("does not fail the test", func() {
				Expect(func() { selection.Click() }).NotTo(Panic())
			})
		})
	})
})
