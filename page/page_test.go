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
		window  *mocks.Window
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		failer = &mocks.Failer{}
		window = &mocks.Window{}
		element = &mocks.Element{}
		page = NewPage(driver, failer)
	})

	Describe("#Navigate", func() {
		Context("when the navigate succeeds", func() {
			It("directs the driver to navigate to the provided URL", func() {
				page.Navigate("http://example.com")
				Expect(driver.NavigateCall.URL).To(Equal("http://example.com"))
			})

			It("returns the page", func() {
				Expect(page.Navigate("http://example.com")).To(Equal(page))
			})

			It("ends with a net-zero caller skip", func() {
				page.Navigate("http://example.com")
				Expect(failer.DownCount).To(Equal(1))
				Expect(failer.UpCount).To(Equal(1))
			})
		})

		Context("when the navigate fails", func() {
			BeforeEach(func() {
				driver.NavigateCall.Err = errors.New("some error")
			})

			It("should fail the test", func() {
				Expect(func() { page.Navigate("http://example.com") }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to navigate: some error"))
			})

			It("fails the test with an caller skip of 1", func() {
				Expect(func() { page.Navigate("http://example.com") }).To(Panic())
				Expect(failer.DownCount).To(Equal(1))
				Expect(failer.UpCount).To(Equal(0))
			})
		})
	})

	Describe("#Size", func() {
		BeforeEach(func() {
			driver.GetWindowCall.ReturnWindow = window
		})

		Context("when the size setting succeeds", func () {
			It("sizes the window correctly", func() {
				page.Size(640,480)
				Expect(window.SizeCall.Width).To(Equal(480))
				Expect(window.SizeCall.Height).To(Equal(640))
			})
		})

		Context("when you fail to get a window", func() {
			BeforeEach(func() {
				driver.GetWindowCall.Err = errors.New("some error")
			})

			It("should fail the test", func() {
				Expect(func() { page.Size(640,480) }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to get a window: some error"))
			})
		})

		Context("when you fail to size the window returned", func() {
			BeforeEach(func() {
				window.SizeCall.Err = errors.New("some error")
			})

			It("should fail the test", func() {
				Expect(func() { page.Size(640,480) }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to re-size the window: some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		Context("when setting the cookie succeeds", func() {
			var cookie webdriver.Cookie

			BeforeEach(func() {
				cookie = webdriver.Cookie{
					Name:     "theName",
					Value:    42,
					Path:     "/my-path",
					Domain:   "example.com",
					Secure:   false,
					HTTPOnly: false,
					Expiry:   1412358590,
				}
			})

			It("instructs the driver to add the cookie to the session", func() {
				page.SetCookie(cookie)
				Expect(driver.SetCookieCall.Cookie.Name).To(Equal("theName"))
				Expect(driver.SetCookieCall.Cookie.Value).To(Equal(42))
			})

			It("returns the page", func() {
				Expect(page.SetCookie(cookie)).To(Equal(page))
			})

			It("ends with a net-zero caller skip", func() {
				page.SetCookie(cookie)
				Expect(failer.DownCount).To(Equal(1))
				Expect(failer.UpCount).To(Equal(1))
			})
		})

		Context("when the webdriver fails to set the cookie", func() {
			BeforeEach(func() {
				driver.SetCookieCall.Err = errors.New("some error")
			})

			It("fails the test with the propagated URL", func() {
				Expect(func() { page.SetCookie(webdriver.Cookie{}) }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to set cookie: some error"))
			})

			It("fails the test with a net-one caller skip", func() {
				Expect(func() { page.SetCookie(webdriver.Cookie{}) }).To(Panic())
				Expect(failer.DownCount).To(Equal(1))
				Expect(failer.UpCount).To(Equal(0))
			})
		})
	})

	Describe("#URL", func() {
		Context("when retrieving the URL is successful", func() {
			It("returns the URL of the current page", func() {
				driver.GetURLCall.ReturnURL = "http://example.com"
				url := page.URL()
				Expect(url).To(Equal("http://example.com"))
			})

			It("ends with a net-zero caller skip", func() {
				page.URL()
				Expect(failer.DownCount).To(Equal(1))
				Expect(failer.UpCount).To(Equal(1))
			})
		})

		Context("when the driver fails to retrieve the URL", func() {
			BeforeEach(func() {
				driver.GetURLCall.Err = errors.New("some error")
			})

			It("fails the test with the propagated URL", func() {
				Expect(func() { page.URL() }).To(Panic())
				Expect(failer.Message).To(Equal("Failed to retrieve URL: some error"))
			})

			It("fails the test with a net-one caller skip", func() {
				Expect(func() { page.URL() }).To(Panic())
				Expect(failer.DownCount).To(Equal(1))
				Expect(failer.UpCount).To(Equal(0))
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

		It("ends with a net-zero caller skip", func() {
			page.Click()
			Expect(failer.DownCount).To(Equal(3))
			Expect(failer.UpCount).To(Equal(3))
		})
	})
})
