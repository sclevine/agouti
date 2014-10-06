package page_test

import (
	. "github.com/sclevine/agouti/page"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/mocks"
	"github.com/sclevine/agouti/webdriver"
	"errors"
)

var _ = Describe("Page", func() {
	var (
		page    Page
		failer  *mocks.Failer
		driver  *mocks.Driver
		element *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		failer = &mocks.Failer{}
		element = &mocks.Element{}
		page = NewPage(driver, failer.Fail)
	})

	Describe("#Navigate", func() {
		It("directs the driver to navigate to the provided URL", func() {
			page.Navigate("http://example.com")
			Expect(driver.NavigateCall.URL).To(Equal("http://example.com"))
		})

		It("returns the page", func() {
			Expect(page.Navigate("http://example.com")).To(Equal(page))
		})

		Context("when the navigate fails", func() {
			BeforeEach(func() {
				driver.NavigateCall.Err = errors.New("some error")
				page.Navigate("http://example.com")
			})

			It("should fail the test", func() {
				Expect(failer.Message).To(Equal("some error"))
			})

			It("fails the test with an offset of 1", func() {
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})
	})

	Describe("#SetCookie", func() {
		It("instructs the driver to add the cookie to the session", func() {
			cookie := webdriver.Cookie{
				Name:     "theName",
				Value:    42,
				Path:     "/my-path",
				Domain:   "example.com",
				Secure:   false,
				HTTPOnly: false,
				Expiry:   1412358590,
			}

			page.SetCookie(cookie)
			Expect(driver.SetCookieCall.Cookie.Name).To(Equal("theName"))
			Expect(driver.SetCookieCall.Cookie.Value).To(Equal(42))
		})

		Context("when the webdriver fails", func() {
			BeforeEach(func() {
				driver.SetCookieCall.Err = errors.New("some error")
				page.SetCookie(webdriver.Cookie{})
			})

			It("should fail the test", func() {
				Expect(failer.Message).To(Equal("some error"))
			})

			It("fails the test with an offset of 1", func() {
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})
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
