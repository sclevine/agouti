package page_test

import (
	. "github.com/sclevine/agouti/page"

	"bytes"
	"encoding/base64"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/internal/mocks"
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
				Expect(driver.NavigateCall.URL).To(Equal("http://example.com"))
			})

			It("returns nil", func() {
				Expect(page.Navigate("http://example.com")).ToNot(HaveOccurred())
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
			Expect(driver.DeleteAllCookiesCall.WasCalled).To(BeTrue())
		})

		Context("when deleteing all cookies succeeds", func() {
			It("returns nil", func() {
				err := page.ClearCookies()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deleting all cookies fails", func() {
			It("returns an error", func() {
				driver.DeleteAllCookiesCall.Err = errors.New("some error")
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
				Expect(window.SizeCall.Width).To(Equal(480))
				Expect(window.SizeCall.Height).To(Equal(640))
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

	Describe("#TakeScreenshot", func() {
		var (
			screenshotPath string
			directory      string
		)

		AfterEach(func() {
			os.Remove(screenshotPath)
		})

		Context("when the driver sucessfully retrieves a screenshot", func() {
			BeforeEach(func() {
				directory, _ = os.Getwd()
				screenshotPath = filepath.Join(directory, ".test.screenshot.png")
				base64Image := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVQYV2P4DwABAQEAWk1v8QAAAABJRU5ErkJggg=="
				imageBytes, _ := base64.StdEncoding.DecodeString(base64Image)
				driver.ScreenshotCall.ReturnImage = bytes.NewBuffer(imageBytes)
			})

			It("successfully saves the screenshot to the filepath provided", func() {
				_, err := os.Stat(screenshotPath)
				Expect(err).To(HaveOccurred())

				page.Screenshot(directory, ".test.screenshot")
				_, err = os.Stat(screenshotPath)
				Expect(err).ToNot(HaveOccurred())
			})

			It("is a PNG", func() {
				page.Screenshot(directory, ".test.screenshot")
				fileBytes, _ := ioutil.ReadFile(screenshotPath)
				Expect(bytes.HasPrefix(fileBytes, []byte("\x89PNG\r\n\x1a\n"))).To(BeTrue())
			})

			It("does not return an error", func() {
				err := page.Screenshot(directory, ".test.screenshot")
				Expect(err).ToNot(HaveOccurred())
			})

			Context("but the screenshot cannot be decoded", func() {
				BeforeEach(func() {
					garbage := "YXNoZGZrbGphc2hkZmtqaGFzbGtkamZoYWxrc2pkaGZha2xqc2RoZmxrYWhzZGZrbGpoYWRzZg=="
					garbageBytes, _ := base64.StdEncoding.DecodeString(garbage)
					driver.ScreenshotCall.ReturnImage = bytes.NewBuffer(garbageBytes)
				})

				It("returns an error", func() {
					driver.ScreenshotCall.Err = errors.New("some error")
					err := page.Screenshot(directory, ".test.screenshot")
					Expect(err).To(MatchError("failed to retrieve screenshot: some error"))
				})
			})
		})

		Context("when the driver fails to retrieve a screenshot", func() {
			BeforeEach(func() {
				driver.ScreenshotCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(page.Screenshot(directory, ".test.screenshot.png")).To(MatchError("failed to retrieve screenshot: some error"))
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
	})
})
