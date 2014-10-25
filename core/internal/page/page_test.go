package page_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/page"
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
		It("should successfully direct the client to delete the session", func() {
			Expect(page.Destroy()).To(Succeed())
			Expect(client.DeleteSessionCall.Called).To(BeTrue())
		})

		Context("when deleting the session fails", func() {
			It("should return the client error", func() {
				client.DeleteSessionCall.Err = errors.New("some error")
				Expect(page.Destroy()).To(MatchError("failed to destroy session: some error"))
			})
		})
	})

	Describe("#Navigate", func() {
		It("should successfully direct the client to navigate to the provided URL", func() {
			Expect(page.Navigate("http://example.com")).To(Succeed())
			Expect(client.SetURLCall.URL).To(Equal("http://example.com"))
		})

		Context("when the navigate fails", func() {
			It("should return an error", func() {
				client.SetURLCall.Err = errors.New("some error")
				Expect(page.Navigate("http://example.com")).To(MatchError("failed to navigate: some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		It("should successfully instruct the client to add the cookie to the session", func() {
			Expect(page.SetCookie("some-name", 42, "/my-path", "example.com", false, false, 1412358590)).To(Succeed())
			Expect(client.SetCookieCall.Cookie.Name).To(Equal("some-name"))
			Expect(client.SetCookieCall.Cookie.Value).To(Equal(42))
		})

		Context("when the client fails to set the cookie", func() {
			It("should return an error", func() {
				client.SetCookieCall.Err = errors.New("some error")
				err := page.SetCookie("some-name", 42, "/my-path", "example.com", false, false, 1412358590)
				Expect(err).To(MatchError("failed to set cookie: some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		It("should successfully instruct the client to delete a named cookie", func() {
			Expect(page.DeleteCookie("some-name")).To(Succeed())
			Expect(client.DeleteCookieCall.Name).To(Equal("some-name"))
		})

		Context("when deleting the named cookie fails", func() {
			It("should return an error", func() {
				client.DeleteCookieCall.Err = errors.New("some error")
				Expect(page.DeleteCookie("some-name")).To(MatchError("failed to delete cookie some-name: some error"))
			})
		})
	})

	Describe("#ClearCookies", func() {
		It("should successfully instruct the client to delete all cookies", func() {
			Expect(page.ClearCookies()).To(Succeed())
			Expect(client.DeleteCookiesCall.Called).To(BeTrue())
		})

		Context("when deleting all cookies fails", func() {
			It("should return an error", func() {
				client.DeleteCookiesCall.Err = errors.New("some error")
				Expect(page.ClearCookies()).To(MatchError("failed to clear cookies: some error"))
			})
		})
	})

	Describe("#URL", func() {
		It("should successfully return the URL of the current page", func() {
			client.GetURLCall.ReturnURL = "http://example.com"
			url, err := page.URL()
			Expect(url).To(Equal("http://example.com"))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the client fails to retrieve the URL", func() {
			It("should return an error", func() {
				client.GetURLCall.Err = errors.New("some error")
				_, err := page.URL()
				Expect(err).To(MatchError("failed to retrieve URL: some error"))
			})
		})
	})

	Describe("#Size", func() {
		It("should set the window width and height to the provided dimensions", func() {
			client.GetWindowCall.ReturnWindow = window
			Expect(page.Size(640, 480)).To(Succeed())
			Expect(window.SizeCall.Width).To(Equal(640))
			Expect(window.SizeCall.Height).To(Equal(480))
		})

		Context("when the client fails to retrieve a window", func() {
			It("should return an error", func() {
				client.GetWindowCall.Err = errors.New("some error")
				Expect(page.Size(640, 480)).To(MatchError("failed to retrieve window: some error"))
			})
		})

		Context("when the window fails to retrieve its size", func() {
			It("should return an error", func() {
				client.GetWindowCall.ReturnWindow = window
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
			It("should return an error", func() {
				err := page.Screenshot("\000/a") // try NUL
				Expect(err).To(MatchError("failed to create directory for screenshot: mkdir \x00: invalid argument"))
			})
		})

		Context("when a new screenshot file cannot be created", func() {
			It("should return an error", func() {
				err := page.Screenshot("")
				Expect(err).To(MatchError("failed to create file for screenshot: open : no such file or directory"))
			})
		})

		Context("when the client fails to retrieve a screenshot", func() {
			BeforeEach(func() {
				client.GetScreenshotCall.Err = errors.New("some error")
			})

			It("should return an error indicating so", func() {
				Expect(page.Screenshot(filename)).To(MatchError("failed to retrieve screenshot: some error"))
			})

			It("should remove the newly-created file", func() {
				page.Screenshot(filename)
				_, err := os.Stat(filename)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the screenshot cannot be written to a file", func() {
			// NOTE: would need to cause write-error to test
		})

		Context("when a screenshot is successfully written to a file", func() {
			It("should successfully saves the screenshot", func() {
				client.GetScreenshotCall.ReturnImage = []byte("some-image")
				Expect(page.Screenshot(filename)).To(Succeed())
				defer os.Remove(filename)
				result, _ := ioutil.ReadFile(filename)
				Expect(string(result)).To(Equal("some-image"))
			})
		})
	})

	Describe("#Title", func() {
		It("should successfully return the title of the current page", func() {
			client.GetTitleCall.ReturnTitle = "Some Title"
			title, err := page.Title()
			Expect(title).To(Equal("Some Title"))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the client fails to retrieve the page title", func() {
			It("should return an error", func() {
				client.GetTitleCall.Err = errors.New("some error")
				_, err := page.Title()
				Expect(err).To(MatchError("failed to retrieve page title: some error"))
			})
		})
	})

	Describe("#HTML", func() {
		It("should return the HTML of the current page", func() {
			client.GetSourceCall.ReturnSource = "Some HTML"
			html, err := page.HTML()
			Expect(html).To(Equal("Some HTML"))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the client fails to retrieve the page HTML", func() {
			It("should return an error", func() {
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

		It("should provide the client with an argument-provided javascript function", func() {
			Expect(client.ExecuteCall.Body).To(Equal("return (function(argument) { some javascript code; }).apply(this, arguments);"))
		})

		It("should provide the client with arguments to call the provided function with", func() {
			Expect(client.ExecuteCall.Arguments).To(Equal([]interface{}{"value"}))
		})

		It("should unmarshall the returned result into the provided result interface", func() {
			Expect(result.Some).To(Equal("result"))
		})

		It("should be successful", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when running the script fails", func() {
			It("should return the client error", func() {
				client.ExecuteCall.Err = errors.New("some error")
				err = page.RunScript("", map[string]interface{}{}, &result)
				Expect(err).To(MatchError("failed to run script: some error"))
			})
		})
	})

	Describe("#PopupText", func() {
		It("should return the popup text of the popup and succeed", func() {
			client.GetAlertTextCall.ReturnText = "some popup text"
			text, err := page.PopupText()
			Expect(text).To(Equal("some popup text"))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the client fails to retrieve the page popup text", func() {
			It("should return an error", func() {
				client.GetAlertTextCall.Err = errors.New("some error")
				_, err := page.PopupText()
				Expect(err).To(MatchError("failed to retrieve popup text: some error"))
			})
		})
	})

	Describe("#EnterPopupText", func() {
		It("should provide the client with the text to enter and succeed", func() {
			Expect(page.EnterPopupText("some text")).To(Succeed())
			Expect(client.SetAlertTextCall.Text).To(Equal("some text"))
		})

		Context("when the client fails to enter the page popup text", func() {
			It("should return an error", func() {
				client.SetAlertTextCall.Err = errors.New("some error")
				Expect(page.EnterPopupText("some text")).To(MatchError("failed to enter popup text: some error"))
			})
		})
	})

	Describe("#ConfirmPopup", func() {
		It("should provide the client with a carriage return and succeed", func() {
			Expect(page.ConfirmPopup()).To(Succeed())
			Expect(client.SetAlertTextCall.Text).To(Equal("\u000d"))
		})

		Context("when the client fails to send the carriage return", func() {
			It("should return an error", func() {
				client.SetAlertTextCall.Err = errors.New("some error")
				Expect(page.ConfirmPopup()).To(MatchError("failed to confirm popup: some error"))
			})
		})
	})

	Describe("#CancelPopup", func() {
		It("should provide the client with an escape and succeed", func() {
			Expect(page.CancelPopup()).To(Succeed())
			Expect(client.SetAlertTextCall.Text).To(Equal("\u001b"))
		})

		Context("when the client fails to send the escape", func() {
			It("should return an error", func() {
				client.SetAlertTextCall.Err = errors.New("some error")
				Expect(page.CancelPopup()).To(MatchError("failed to cancel popup: some error"))
			})
		})
	})

	Describe("#Forward", func() {
		It("should successfully instruct the client to move forward in history", func() {
			Expect(page.Forward()).To(Succeed())
			Expect(client.ForwardCall.Called).To(BeTrue())
		})

		Context("when navigating forward fails", func() {
			It("should return an error", func() {
				client.ForwardCall.Err = errors.New("some error")
				Expect(page.Forward()).To(MatchError("failed to navigate forward in history: some error"))
			})
		})
	})

	Describe("#Back", func() {
		It("should successfully instruct the client to move back in history", func() {
			Expect(page.Back()).To(Succeed())
			Expect(client.BackCall.Called).To(BeTrue())
		})

		Context("when navigating back fails", func() {
			It("should return an error", func() {
				client.BackCall.Err = errors.New("some error")
				Expect(page.Back()).To(MatchError("failed to navigate backwards in history: some error"))
			})
		})
	})

	Describe("#Refresh", func() {
		It("should successfully instruct the client to refresh", func() {
			Expect(page.Refresh()).To(Succeed())
			Expect(client.RefreshCall.Called).To(BeTrue())
		})

		Context("when refreshing the page fails", func() {
			It("should return an error", func() {
				client.RefreshCall.Err = errors.New("some error")
				Expect(page.Refresh()).To(MatchError("failed to refresh page: some error"))
			})
		})
	})

	Describe("#Find", func() {
		It("should defer to selection#Find", func() {
			Expect(page.Find("#selector").String()).To(Equal("CSS: #selector [single]"))
		})
	})

	Describe("#FindByXPath", func() {
		It("should defer to selection#FindXByPath", func() {
			Expect(page.FindByXPath("//selector").String()).To(Equal("XPath: //selector [single]"))
		})
	})

	Describe("#FindByLink", func() {
		It("should defer to selection#FindByLink", func() {
			Expect(page.FindByLink("some text").String()).To(Equal(`Link: "some text" [single]`))
		})
	})

	Describe("#FindByLabel", func() {
		It("should defer to selection#FindByLabel", func() {
			Expect(page.FindByLabel("label name").String()).To(ContainSubstring("XPath: //input"))
			Expect(page.FindByLabel("label name").String()).To(ContainSubstring("[single]"))
		})
	})

	Describe("#All", func() {
		It("should defer to selection#All", func() {
			Expect(page.All("#selector").String()).To(Equal("CSS: #selector"))
		})
	})

	Describe("#AllByXPath", func() {
		It("should defer to selection#AllByXPath", func() {
			Expect(page.AllByXPath("//selector").String()).To(Equal("XPath: //selector"))
		})
	})

	Describe("#AllByLink", func() {
		It("should defer to selection#AllByLink", func() {
			Expect(page.AllByLink("some text").String()).To(Equal(`Link: "some text"`))
		})
	})

	Describe("#AllByLabel", func() {
		It("should defer to selection#AllByLabel", func() {
			Expect(page.AllByLabel("label name").String()).To(ContainSubstring("XPath: //input"))
		})
	})
})
