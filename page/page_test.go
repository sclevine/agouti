package page_test

import (
	. "github.com/sclevine/agouti/page"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/internal/mocks"
	"github.com/sclevine/agouti/page/internal/webdriver"
)

var _ = Describe("Page", func() {
	var (
		page    Page
		driver  *mocks.Driver
		element *mocks.Element
		window  *mocks.Window
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		window = &mocks.Window{}
		element = &mocks.Element{}
		page = Page{driver}
	})

	Describe("#Navigate", func() {
		Context("when the navigate succeeds", func() {
			It("directs the driver to navigate to the provided URL", func() {
				page.Navigate("http://example.com")
				Expect(driver.NavigateCall.URL).To(Equal("http://example.com"))
			})

			It("returns nil", func() {
				Expect(page.Navigate("http://example.com")).To(BeNil())
			})
		})

		Context("when the navigate fails", func() {
			BeforeEach(func() {
				driver.NavigateCall.Err = errors.New("some error")
			})

			It("returns the driver error", func() {
				Expect(page.Navigate("http://example.com")).To(MatchError("failed to navigate: some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		It("instructs the driver to add the cookie to the session", func() {
			page.SetCookie("theName", 42, "/my-path", "example.com", false, false, 1412358590)
			Expect(driver.SetCookieCall.Cookie.Name).To(Equal("theName"))
			Expect(driver.SetCookieCall.Cookie.Value).To(Equal(42))
		})

		Context("when setting the cookie succeeds", func() {
			It("returns nil", func() {
				err := page.SetCookie("theName", 42, "/my-path", "example.com", false, false, 1412358590)
				Expect(err).To(BeNil())
			})
		})

		Context("when the driver fails to set the cookie", func() {
			It("returns an error", func() {
				driver.SetCookieCall.Err = errors.New("some error")
				err := page.SetCookie("theName", 42, "/my-path", "example.com", false, false, 1412358590)
				Expect(err).To(MatchError("failed to set cookie: some error"))
			})
		})
	})

	Describe("#URL", func() {
		Context("when retrieving the URL is successful", func() {
			var (
				url string
				err error
			)

			BeforeEach(func() {
				driver.GetURLCall.ReturnURL = "http://example.com"
				url, err = page.URL()
			})

			It("returns the URL of the current page", func() {
				Expect(url).To(Equal("http://example.com"))
			})

			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the driver fails to retrieve the URL", func() {
			It("returns an error", func() {
				driver.GetURLCall.Err = errors.New("some error")
				_, err := page.URL()
				Expect(err).To(MatchError("failed to retrieve URL: some error"))
			})
		})
	})

	Describe("#Size", func() {
		BeforeEach(func() {
			driver.GetWindowCall.ReturnWindow = window
		})

		Context("when the size setting succeeds", func() {
			It("sets the window width and height accordingly", func() {
				page.Size(640, 480)
				Expect(window.SizeCall.Width).To(Equal(480))
				Expect(window.SizeCall.Height).To(Equal(640))
			})

			It("does not return an error", func() {
				Expect(page.Size(640, 480)).To(BeNil())
			})
		})

		Context("when the driver fails to retrieve a window", func() {
			BeforeEach(func() {
				driver.GetWindowCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(page.Size(640, 480)).To(MatchError("failed to retrieve window: some error"))
			})
		})

		Context("when the window fails to retrieve its size", func() {
			BeforeEach(func() {
				window.SizeCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(page.Size(640, 480)).To(MatchError("failed to set window size: some error"))
			})
		})
	})

	Describe("#Find", func() {
		It("returns a selection", func() {
			Expect(page.Find("#selector").Selector()).To(Equal("#selector"))
		})
	})

	Describe("#Selector", func() {
		It("returns body as the selector", func() {
			Expect(page.Selector()).To(Equal("body"))
		})
	})

	Describe("methods that defer to a selection on the page body", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
		})

		Describe("#Click", func() {
			It("calls selection#Click() on the body of the page", func() {
				element.ClickCall.Err = errors.New("some error")
				Expect(page.Click()).To(MatchError("failed to click on selector 'body': some error"))
			})
		})

		Describe("#Text", func() {
			It("calls selection#Text() on the body of the page", func() {
				element.GetTextCall.Err = errors.New("some error")
				_, err := page.Text()
				Expect(err).To(MatchError("failed to retrieve text for selector 'body': some error"))
			})
		})

		Describe("#Attribute", func() {
			It("calls selection#Click() on the body of the page", func() {
				element.GetAttributeCall.Err = errors.New("some error")
				_, err := page.Attribute("some-attribute")
				Expect(err).To(MatchError("failed to retrieve attribute value for selector 'body': some error"))
			})
		})
	})
})
