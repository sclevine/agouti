package page_test

import (
	. "github.com/sclevine/agouti/page"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/mocks"
	"github.com/sclevine/agouti/webdriver"
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
		page = NewPage(driver, failer)
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
			})

			It("should fail the test", func() {
				Expect(func() { page.Navigate("http://example.com") }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to navigate: some error"))
			})

			It("fails the test with an offset of 1", func() {
				Expect(func() { page.Navigate("http://example.com") }).To(Panic())
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
			})

			It("fails the test with the propagated URL", func() {
				Expect(func() { page.SetCookie(webdriver.Cookie{}) }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to set cookie: some error"))
			})

			It("fails the test with an offset of 1", func() {
				Expect(func() { page.SetCookie(webdriver.Cookie{}) }).To(Panic())
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})
	})

	Describe("#URL", func() {
		Context("when the driver fails to retrieve the URL", func() {
			BeforeEach(func() {
				driver.GetURLCall.Err = errors.New("some error")
			})

			It("fails the test with the propagated URL", func() {
				Expect(func() { page.URL() }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to retrieve URL: some error"))
			})

			It("fails the ", func() {
				Expect(func() { page.URL() }).To(Panic())
				Expect(failer.CallerSkip).To(Equal(1))
			})
		})

		Context("when the driver successfully retrieves the URL", func() {
			It("returns the URL of the current page", func() {
				driver.GetURLCall.ReturnURL = "http://example.com"
				url := page.URL()
				Expect(url).To(Equal("http://example.com"))
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

	Describe("#Should", func() {
		It("returns a final selector for the body of the page", func() {
			Expect(page.Should().Selector()).To(Equal("body"))
		})
	})

	Describe("#ShouldNot", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
			element.GetTextCall.ReturnText = "element text"
		})

		It("inverts the selector matcher", func() {
			Expect(func() { page.ShouldNot().ContainText("ment tex") }).To(Panic())
			Expect(func() { page.ShouldNot().ContainText("banana") }).NotTo(Panic())
		})
	})

	Describe("#Click", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
		})

		It("calls selection#Click on the body of the page", func() {
			page.Click()
			Expect(element.ClickCall.Called).To(BeTrue())
			Expect(driver.GetElementsCall.Selector).To(Equal("body"))
		})

		It("increments the caller skip", func() {
			page.Click()
			Expect(failer.CallerSkip).To(Equal(2))
		})
	})
})
