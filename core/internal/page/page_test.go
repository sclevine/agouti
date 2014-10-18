package page_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/page"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Page", func() {
	var (
		page    *Page
		client  *mocks.Client
		element *mocks.Element
		window  *mocks.Window
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		window = &mocks.Window{}
		element = &mocks.Element{}
		page = &Page{client}
	})

	Describe("#Destroy", func() {
		Context("when deleting the session succeeds", func() {
			It("directs the client to delete the session", func() {
				page.Destroy()
				Expect(client.DeleteSessionCall.Called).To(BeTrue())
			})

			It("does not return an error", func() {
				err := page.Destroy()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deleting the session fails", func() {
			It("returns the client error", func() {
				client.DeleteSessionCall.Err = errors.New("some error")
				Expect(page.Destroy()).To(MatchError("failed to destroy session: some error"))
			})
		})
	})

	Describe("#Navigate", func() {
		Context("when the navigate succeeds", func() {
			It("directs the client to navigate to the provided URL", func() {
				page.Navigate("http://example.com")
				Expect(client.SetURLCall.URL).To(Equal("http://example.com"))
			})

			It("returns nil", func() {
				err := page.Navigate("http://example.com")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the navigate fails", func() {
			It("returns the client error", func() {
				client.SetURLCall.Err = errors.New("some error")
				Expect(page.Navigate("http://example.com")).To(MatchError("failed to navigate: some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		It("instructs the client to add the cookie to the session", func() {
			page.SetCookie("some-name", 42, "/my-path", "example.com", false, false, 1412358590)
			Expect(client.SetCookieCall.Cookie.Name).To(Equal("some-name"))
			Expect(client.SetCookieCall.Cookie.Value).To(Equal(42))
		})

		Context("when setting the cookie succeeds", func() {
			It("returns nil", func() {
				err := page.SetCookie("some-name", 42, "/my-path", "example.com", false, false, 1412358590)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the client fails to set the cookie", func() {
			It("returns an error", func() {
				client.SetCookieCall.Err = errors.New("some error")
				err := page.SetCookie("some-name", 42, "/my-path", "example.com", false, false, 1412358590)
				Expect(err).To(MatchError("failed to set cookie: some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		It("instructs the client to delete a named cookie", func() {
			page.DeleteCookie("some-name")
			Expect(client.DeleteCookieCall.Name).To(Equal("some-name"))
		})

		Context("when deleteing the named cookie succeeds", func() {
			It("returns nil", func() {
				err := page.DeleteCookie("some-name")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deleting the named cookie fails", func() {
			It("returns an error", func() {
				client.DeleteCookieCall.Err = errors.New("some error")
				err := page.DeleteCookie("some-name")
				Expect(err).To(MatchError("failed to delete cookie some-name: some error"))
			})
		})
	})

	Describe("#ClearCookies", func() {
		It("instructs the client to delete all cookies", func() {
			page.ClearCookies()
			Expect(client.DeleteCookiesCall.Called).To(BeTrue())
		})

		Context("when deleteing all cookies succeeds", func() {
			It("returns nil", func() {
				err := page.ClearCookies()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deleting all cookies fails", func() {
			It("returns an error", func() {
				client.DeleteCookiesCall.Err = errors.New("some error")
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
				client.GetURLCall.ReturnURL = "http://example.com"
				url, err = page.URL()
			})

			It("returns the URL of the current page", func() {
				Expect(url).To(Equal("http://example.com"))
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the client fails to retrieve the URL", func() {
			It("returns an error", func() {
				client.GetURLCall.Err = errors.New("some error")
				_, err := page.URL()
				Expect(err).To(MatchError("failed to retrieve URL: some error"))
			})
		})
	})

	Describe("#Size", func() {
		BeforeEach(func() {
			client.GetWindowCall.ReturnWindow = window
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

		Context("when the client fails to retrieve a window", func() {
			It("returns an error", func() {
				client.GetWindowCall.Err = errors.New("some error")
				Expect(page.Size(640, 480)).To(MatchError("failed to retrieve window: some error"))
			})
		})

		Context("when the window fails to retrieve its size", func() {
			It("returns an error", func() {
				window.SizeCall.Err = errors.New("some error")
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

		Context("when the client fails to retrieve a screenshot", func() {
			BeforeEach(func() {
				client.GetScreenshotCall.Err = errors.New("some error")
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
				client.GetScreenshotCall.ReturnImage = []byte("some-image")
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
				client.GetTitleCall.ReturnTitle = "Some Title"
				title, err = page.Title()
			})

			It("returns the title of the current page", func() {
				Expect(title).To(Equal("Some Title"))
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the client fails to retrieve the page title", func() {
			It("returns an error", func() {
				client.GetTitleCall.Err = errors.New("some error")
				_, err := page.Title()
				Expect(err).To(MatchError("failed to retrieve page title: some error"))
			})
		})
	})

	Describe("#HTML", func() {
		Context("when retrieving the page HTML is successful", func() {
			var (
				html string
				err  error
			)

			BeforeEach(func() {
				client.GetSourceCall.ReturnSource = "Some HTML"
				html, err = page.HTML()
			})

			It("returns the HTML of the current page", func() {
				Expect(html).To(Equal("Some HTML"))
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the client fails to retrieve the page HTML", func() {
			It("returns an error", func() {
				client.GetSourceCall.Err = errors.New("some error")
				_, err := page.HTML()
				Expect(err).To(MatchError("failed to retrieve page HTML: some error"))
			})
		})
	})

	Describe("#RunScript", func() {
		var (
			result struct{ Some string }
			err    error
		)

		BeforeEach(func() {
			client.ExecuteCall.Result = `{"some": "result"}`
			err = page.RunScript("some javascript code", map[string]interface{}{"argument": "value"}, &result)
		})

		It("provides the client with an argument-provided javascript function", func() {
			Expect(client.ExecuteCall.Body).To(Equal("return (function(argument) { some javascript code; }).apply(this, arguments);"))
		})

		It("provides the client with arguments to call the provided function with", func() {
			Expect(client.ExecuteCall.Arguments).To(Equal([]interface{}{"value"}))
		})

		It("unmarshalls the returned result into the provided result interface", func() {
			Expect(result.Some).To(Equal("result"))
		})

		Context("when executing the script succeeds", func() {
			It("returns nil", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when running the script fails", func() {
			It("returns the client error", func() {
				client.ExecuteCall.Err = errors.New("some error")
				err = page.RunScript("", map[string]interface{}{}, &result)
				Expect(err).To(MatchError("failed to run script: some error"))
			})
		})
	})

	Describe("#Forward", func() {
		It("instructs the client to move forward in history", func() {
			page.Forward()
			Expect(client.ForwardCall.Called).To(BeTrue())
		})

		Context("when navigating forward succeeds", func() {
			It("does not return an error", func() {
				err := page.Forward()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when navigating forward fails", func() {
			It("returns an error", func() {
				client.ForwardCall.Err = errors.New("some error")
				err := page.Forward()
				Expect(err).To(MatchError("failed to navigate forward in history: some error"))
			})
		})
	})

	Describe("#Back", func() {
		It("instructs the client to move back in history", func() {
			page.Back()
			Expect(client.BackCall.Called).To(BeTrue())
		})

		Context("when navigating back succeeds", func() {
			It("does not return an error", func() {
				err := page.Back()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when navigating back fails", func() {
			It("returns an error", func() {
				client.BackCall.Err = errors.New("some error")
				err := page.Back()
				Expect(err).To(MatchError("failed to navigate backwards in history: some error"))
			})
		})
	})

	Describe("#Refresh", func() {
		It("instructs the client to refresh", func() {
			page.Refresh()
			Expect(client.RefreshCall.Called).To(BeTrue())
		})

		Context("when navigating refresh succeeds", func() {
			It("does not return an error", func() {
				err := page.Refresh()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when refreshing the page fails", func() {
			It("returns an error", func() {
				client.RefreshCall.Err = errors.New("some error")
				err := page.Refresh()
				Expect(err).To(MatchError("failed to refresh page: some error"))
			})
		})
	})

	Describe("#Find", func() {
		It("defers to selection#Find", func() {
			Expect(page.Find("#selector").String()).To(Equal("CSS: #selector"))
		})
	})

	Describe("#FindXPath", func() {
		It("defers to selection#FindXPath", func() {
			Expect(page.FindXPath("//selector").String()).To(Equal("XPath: //selector"))
		})
	})

	Describe("#FindLink", func() {
		It("defers to selection#FindLink", func() {
			Expect(page.FindLink("some text").String()).To(Equal(`Link: "some text"`))
		})
	})

	Describe("#FindByLabel", func() {
		It("defers to selection#FindByLabel", func() {
			Expect(page.FindByLabel("label name").String()).To(ContainSubstring(`XPath: //input`))
		})
	})
})
