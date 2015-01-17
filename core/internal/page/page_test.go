package page_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/page"
)

var _ = Describe("Page", func() {
	var (
		page   *Page
		client *mocks.Client
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		page = &Page{Client: client}
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
			Expect(page.SetCookie("some cookie")).To(Succeed())
			Expect(client.SetCookieCall.Cookie).To(Equal("some cookie"))
		})

		Context("when the client fails to set the cookie", func() {
			It("should return an error", func() {
				client.SetCookieCall.Err = errors.New("some error")
				err := page.SetCookie("some cookie")
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
		var (
			windowSession *mocks.Session
			window        *api.Window
		)

		BeforeEach(func() {
			windowSession = &mocks.Session{}
			window = &api.Window{Session: windowSession}
		})

		It("should set the window width and height to the provided dimensions", func() {
			client.GetWindowCall.ReturnWindow = window
			Expect(page.Size(640, 480)).To(Succeed())
			Expect(windowSession.ExecuteCall.BodyJSON).To(MatchJSON(`{"width": 640, "height": 480}`))
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
				windowSession.ExecuteCall.Err = errors.New("some error")
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

	Describe("#SwitchToParentFrame", func() {
		It("should successfully instruct the client to change focus to the parent frame", func() {
			Expect(page.SwitchToParentFrame()).To(Succeed())
			Expect(client.FrameParentCall.Called).To(BeTrue())
		})

		Context("when switching to the parent frame fails", func() {
			It("should return an error", func() {
				client.FrameParentCall.Err = errors.New("some error")
				Expect(page.SwitchToParentFrame()).To(MatchError("failed to switch to parent frame: some error"))
			})
		})
	})

	Describe("#SwitchToRootFrame", func() {
		It("should successfully instruct the client to change focus to the root frame", func() {
			client.FrameCall.Frame = &api.Element{}
			Expect(page.SwitchToRootFrame()).To(Succeed())
			Expect(client.FrameCall.Frame).To(BeNil())
		})

		Context("when switching to the root frame fails", func() {
			It("should return an error", func() {
				client.FrameCall.Err = errors.New("some error")
				Expect(page.SwitchToRootFrame()).To(MatchError("failed to switch to original page frame: some error"))
			})
		})
	})

	Describe("#SwitchToWindow", func() {
		It("should successfully instruct the client to switch to the named window", func() {
			Expect(page.SwitchToWindow("some name")).To(Succeed())
			Expect(client.SetWindowByNameCall.Name).To(Equal("some name"))
		})

		Context("when switching to the root frame fails", func() {
			It("should return an error", func() {
				client.SetWindowByNameCall.Err = errors.New("some error")
				Expect(page.SwitchToWindow("some name")).To(MatchError("failed to switch to named window: some error"))
			})
		})
	})

	Describe("#NextWindow", func() {
		BeforeEach(func() {
			firstWindow := &api.Window{ID: "first window"}
			secondWindow := &api.Window{ID: "second window"}
			thirdWindow := &api.Window{ID: "third window"}
			client.GetWindowsCall.ReturnWindows = []*api.Window{secondWindow, firstWindow, thirdWindow}
			client.GetWindowCall.ReturnWindow = firstWindow
		})

		It("should successfully instruct the client to switch to the next window in sorted order", func() {
			Expect(page.NextWindow()).To(Succeed())
			Expect(client.SetWindowCall.Window.ID).To(Equal("second window"))
		})

		Context("when retrieving the available windows fails", func() {
			It("should return an error", func() {
				client.GetWindowsCall.Err = errors.New("some error")
				Expect(page.NextWindow()).To(MatchError("failed to find available windows: some error"))
			})
		})

		Context("when retrieving the active window fails", func() {
			It("should return an error", func() {
				client.GetWindowCall.Err = errors.New("some error")
				Expect(page.NextWindow()).To(MatchError("failed to find active window: some error"))
			})
		})

		Context("when setting the active window fails", func() {
			It("should return an error", func() {
				client.SetWindowCall.Err = errors.New("some error")
				Expect(page.NextWindow()).To(MatchError("failed to change active window: some error"))
			})
		})
	})

	Describe("#CloseWindow", func() {
		It("should successfully instruct the client to close the active window", func() {
			Expect(page.CloseWindow()).To(Succeed())
			Expect(client.DeleteWindowCall.Called).To(BeTrue())
		})

		Context("when closing the active window fails", func() {
			It("should return an error", func() {
				client.DeleteWindowCall.Err = errors.New("some error")
				Expect(page.CloseWindow()).To(MatchError("failed to close active window: some error"))
			})
		})
	})

	Describe("#WindowCount", func() {
		It("should successfully return the number of windows from the client", func() {
			client.GetWindowsCall.ReturnWindows = []*api.Window{&api.Window{}, &api.Window{}}
			count, err := page.WindowCount()
			Expect(count).To(Equal(2))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when retrieving the available windows fails", func() {
			It("should return an error", func() {
				client.GetWindowsCall.Err = errors.New("some error")
				_, err := page.WindowCount()
				Expect(err).To(MatchError("failed to find available windows: some error"))
			})
		})
	})

	Describe("#ReadLogs", func() {
		It("should request logs of the provided log type from the client", func() {
			_, err := page.ReadLogs("some type")
			Expect(err).To(Succeed())
			Expect(client.NewLogsCall.LogType).To(Equal("some type"))
		})

		Context("when the client fails to retrieve logs", func() {
			It("should return an error", func() {
				client.NewLogsCall.Err = errors.New("some error")
				_, err := page.ReadLogs("some type")
				Expect(err).To(MatchError("failed to retrieve logs: some error"))
			})
		})

		Describe("returned logs", func() {
			BeforeEach(func() {
				client.NewLogsCall.ReturnLogs = []api.Log{api.Log{"old log", "old level", 1418196096123}}
				page.ReadLogs("some type")
				client.NewLogsCall.ReturnLogs = []api.Log{api.Log{"new log (1:22)", "new level", 1418196097543}}
			})

			Context("when only new logs are requested", func() {
				It("should return new logs with the correct time and code location", func() {
					logs, err := page.ReadLogs("some type")
					Expect(err).NotTo(HaveOccurred())
					Expect(logs).To(HaveLen(1))
					Expect(logs[0].Message).To(Equal("new log"))
					Expect(logs[0].Location).To(Equal("1:22"))
					Expect(logs[0].Level).To(Equal("new level"))
					Expect(logs[0].Time.Unix()).To(BeEquivalentTo(1418196097))
				})
			})

			Context("when all logs are requested", func() {
				It("should return all logs", func() {
					logs, err := page.ReadLogs("some type", true)
					Expect(err).NotTo(HaveOccurred())
					Expect(logs).To(HaveLen(2))
					Expect(logs[0].Message).To(Equal("old log"))
					Expect(logs[1].Message).To(Equal("new log"))
				})
			})
		})
	})

	Describe("#LogTypes", func() {
		It("should successfully return the log types", func() {
			client.GetLogTypesCall.ReturnTypes = []string{"first type", "second type"}
			types, err := page.LogTypes()
			Expect(types).To(Equal([]string{"first type", "second type"}))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the client fails to retrieve the log types", func() {
			It("should return an error", func() {
				client.GetLogTypesCall.Err = errors.New("some error")
				_, err := page.LogTypes()
				Expect(err).To(MatchError("failed to retrieve log types: some error"))
			})
		})
	})
})
