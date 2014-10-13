package page_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/sclevine/agouti/page"
	"github.com/sclevine/agouti/page/internal/webdriver"
	"io/ioutil"
	"os"
	"path/filepath"
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
				Expect(driver.SetURLCall.URL).To(Equal("http://example.com"))
			})

			It("returns nil", func() {
				Expect(page.Navigate("http://example.com")).ToNot(HaveOccurred())
			})
		})

		Context("when the navigate fails", func() {
			BeforeEach(func() {
				driver.SetURLCall.Err = errors.New("some error")
			})

			It("returns the driver error", func() {
				Expect(page.Navigate("http://example.com")).To(MatchError("failed to navigate: some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		It("instructs the driver to add the cookie to the session", func() {
			page.SetCookie("some-name", 42, "/my-path", "example.com", false, false, 1412358590)
			Expect(driver.SetCookieCall.Cookie.Name).To(Equal("some-name"))
			Expect(driver.SetCookieCall.Cookie.Value).To(Equal(42))
		})

		Context("when setting the cookie succeeds", func() {
			It("returns nil", func() {
				err := page.SetCookie("some-name", 42, "/my-path", "example.com", false, false, 1412358590)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the driver fails to set the cookie", func() {
			It("returns an error", func() {
				driver.SetCookieCall.Err = errors.New("some error")
				err := page.SetCookie("some-name", 42, "/my-path", "example.com", false, false, 1412358590)
				Expect(err).To(MatchError("failed to set cookie: some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		It("instructs the driver to delete a named cookie", func() {
			page.DeleteCookie("some-name")
			Expect(driver.DeleteCookieCall.Name).To(Equal("some-name"))
		})

		Context("when deleteing the named cookie succeeds", func() {
			It("returns nil", func() {
				err := page.DeleteCookie("some-name")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deleting the named cookie fails", func() {
			It("returns an error", func() {
				driver.DeleteCookieCall.Err = errors.New("some error")
				err := page.DeleteCookie("some-name")
				Expect(err).To(MatchError("failed to delete cookie some-name: some error"))
			})
		})
	})

	Describe("#ClearCookies", func() {
		It("instructs the driver to delete all cookies", func() {
			page.ClearCookies()
			Expect(driver.DeleteCookiesCall.WasCalled).To(BeTrue())
		})

		Context("when deleteing all cookies succeeds", func() {
			It("returns nil", func() {
				err := page.ClearCookies()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deleting all cookies fails", func() {
			It("returns an error", func() {
				driver.DeleteCookiesCall.Err = errors.New("some error")
				err := page.ClearCookies()
				Expect(err).To(MatchError("failed to clear cookies: some error"))
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
				Expect(err).ToNot(HaveOccurred())
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
				Expect(window.SizeCall.Width).To(Equal(640))
				Expect(window.SizeCall.Height).To(Equal(480))
			})

			It("does not return an error", func() {
				Expect(page.Size(640, 480)).ToNot(HaveOccurred())
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

	Describe("#Screenshot", func() {
		var filename string

		BeforeEach(func() {
			directory, _ := os.Getwd()
			filename = filepath.Join(directory, ".test.screenshot.png")
		})

		Context("when the file path cannot be constructed", func() {
			It("returns an error indicating that it could not create a directory", func() {
				err := page.Screenshot("\000/a") // try NUL
				Expect(err).To(MatchError("failed to create directory for screenshot: mkdir \x00: invalid argument"))
			})
		})

		Context("when a new screenshot file cannot be created", func() {
			It("returns an error indicating so", func() {
				err := page.Screenshot("")
				Expect(err).To(MatchError("failed to create file for screenshot: open : no such file or directory"))
			})
		})

		Context("when the driver fails to retrieve a screenshot", func() {
			BeforeEach(func() {
				driver.GetScreenshotCall.Err = errors.New("some error")
			})

			It("returns an error indicating so", func() {
				err := page.Screenshot(filename)
				Expect(err).To(MatchError("failed to retrieve screenshot: some error"))
			})

			It("removes the newly-created file", func() {
				page.Screenshot(filename)
				_, err := os.Stat(filename)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the screenshot cannot be written to a file", func() {
			// NOTE: would need to cause write-error to test
		})

		Context("when a screenshot is successfully written to a file", func() {
			var err error

			BeforeEach(func() {
				driver.GetScreenshotCall.ReturnImage = []byte("some-image")
				err = page.Screenshot(filename)
			})

			AfterEach(func() {
				os.Remove(filename)
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("successfully saves the screenshot", func() {
				result, _ := ioutil.ReadFile(filename)
				Expect(string(result)).To(Equal("some-image"))
			})
		})
	})

	Describe("#Title", func() {
		Context("when retrieving the page title is successful", func() {
			var (
				title string
				err   error
			)

			BeforeEach(func() {
				driver.GetTitleCall.ReturnTitle = "Some Title"
				title, err = page.Title()
			})

			It("returns the title of the current page", func() {
				Expect(title).To(Equal("Some Title"))
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the driver fails to retrieve the page title", func() {
			It("returns an error", func() {
				driver.GetTitleCall.Err = errors.New("some error")
				_, err := page.Title()
				Expect(err).To(MatchError("failed to retrieve page title: some error"))
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

		Describe("#Fill", func() {
			It("calls selection#Fill() with text on the body of the page", func() {
				element.ValueCall.Err = errors.New("some error")
				Expect(page.Fill("some text")).To(MatchError("failed to enter text into selector 'body': some error"))
			})
		})

		Describe("#Check", func() {
			It("calls selection#Check() on the body of the page", func() {
				element.ClickCall.Err = errors.New("some error")
				Expect(page.Check()).To(MatchError("selector 'body' does not refer to a checkbox"))
			})
		})

		Describe("#Uncheck", func() {
			It("calls selection#Uncheck() on the body of the page", func() {
				element.ClickCall.Err = errors.New("some error")
				Expect(page.Uncheck()).To(MatchError("selector 'body' does not refer to a checkbox"))
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
			It("calls selection#Attribute() on the body of the page", func() {
				element.GetAttributeCall.Err = errors.New("some error")
				_, err := page.Attribute("some-attribute")
				Expect(err).To(MatchError("failed to retrieve attribute value for selector 'body': some error"))
			})
		})

		Describe("#CSS", func() {
			It("calls selection#CSS() on the body of the page", func() {
				element.GetCSSCall.Err = errors.New("some error")
				_, err := page.CSS("some-property")
				Expect(err).To(MatchError("failed to retrieve CSS property for selector 'body': some error"))
			})
		})

		Describe("#Selected", func() {
			It("calls selection#Selected() on the body of the page", func() {
				element.IsSelectedCall.Err = errors.New("some error")
				_, err := page.Selected()
				Expect(err).To(MatchError("failed to determine whether selector 'body' is selected: some error"))
			})
		})

		Describe("#Select", func() {
			It("calls selection#Select() on the body of the page", func() {
				driver.GetElementsCall.Err = errors.New("some error")
				err := page.Select("some text")
				Expect(err).To(MatchError("failed to retrieve options for selector 'body': some error"))
			})
		})

		Describe("#Submit", func() {
			It("calls selection#Submit() on the body of the page", func() {
				element.SubmitCall.Err = errors.New("some error")
				Expect(page.Submit()).To(MatchError("failed to submit selector 'body': some error"))
			})
		})
	})
})
