package page_test

import (
	. "github.com/sclevine/agouti/page"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/mocks"
	"github.com/sclevine/agouti/webdriver"
)

var _ = Describe("Page", func() {
	var (
		page 	  *Page
		failer    *mocks.Failer
		driver    *mocks.Driver
		element	  *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		failer = &mocks.Failer{}
		element = &mocks.Element{}
		page = &Page{driver, failer.Fail}
	})

	Describe("#Within", func() {
		It("returns a selection", func() {
			selection := page.Within("#selector")
			Expect(selection.Selector()).To(Equal("#selector"))
		})

		It("provides a subselection to any specified callable body functions", func() {
			page.Within("#selector",
				Do(func(selection Selection) {
					Expect(selection.Selector()).To(Equal("#selector"))
				}),
				Do(func(selection Selection) {
					Expect(selection.Selector()).To(Equal("#selector"))
				}),
			)
		})
	})

	Describe("#Selector", func() {
		It("returns body as the selector", func() {
			Expect(page.Selector()).To(Equal("body"))
		})
	})

	Describe("#ShouldContainText", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
			element.GetTextCall.ReturnText = "element text"
		})

		It("calls selection#ShouldContainText on the body of the page", func() {
			page.ShouldContainText("ment tex")
			Expect(driver.GetElementsCall.Selector).To(Equal("body"))
		})

		It("passes on contained text", func() {
			page.ShouldContainText("ment tex")
			Expect(failer.Failed).To(BeFalse())
		})

		It("fails on non-contained text", func() {
			page.ShouldContainText("banana")
			Expect(failer.Failed).To(BeTrue())
		})
	})
})
