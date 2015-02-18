package agouti_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/mocks"
)

var _ = Describe("Page", func() {
	var (
		page    *Page
		session *mocks.Session
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		page = NewTestPage(session)
	})

	Describe("#Destroy", func() {
		It("should successfully delete the session", func() {
			Expect(page.Destroy()).To(Succeed())
			Expect(session.DeleteCall.Called).To(BeTrue())
		})

		Context("when deleting the session fails", func() {
			It("should return an error", func() {
				session.DeleteCall.Err = errors.New("some error")
				Expect(page.Destroy()).To(MatchError("failed to destroy session: some error"))
			})
		})
	})

	Describe("#Navigate", func() {
		It("should successfully instruct the session to navigate to the provided URL", func() {
			Expect(page.Navigate("http://example.com")).To(Succeed())
			Expect(session.SetURLCall.URL).To(Equal("http://example.com"))
		})

		Context("when the navigate fails", func() {
			It("should return an error", func() {
				session.SetURLCall.Err = errors.New("some error")
				Expect(page.Navigate("http://example.com")).To(MatchError("failed to navigate: some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		It("should successfully instruct the session to add the cookie to the session", func() {
			Expect(page.SetCookie(Cookie{"name": "some cookie"})).To(Succeed())
			Expect(session.SetCookieCall.Cookie).To(BeEquivalentTo(Cookie{"name": "some cookie"}))
		})

		Context("when the session fails to set the cookie", func() {
			It("should return an error", func() {
				session.SetCookieCall.Err = errors.New("some error")
				err := page.SetCookie(Cookie{})
				Expect(err).To(MatchError("failed to set cookie: some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		It("should successfully instruct the session to delete a named cookie", func() {
			Expect(page.DeleteCookie("some-name")).To(Succeed())
			Expect(session.DeleteCookieCall.Name).To(Equal("some-name"))
		})

		Context("when deleting the named cookie fails", func() {
			It("should return an error", func() {
				session.DeleteCookieCall.Err = errors.New("some error")
				Expect(page.DeleteCookie("some-name")).To(MatchError("failed to delete cookie some-name: some error"))
			})
		})
	})

	Describe("#ClearCookies", func() {
		It("should successfully instruct the session to delete all cookies", func() {
			Expect(page.ClearCookies()).To(Succeed())
			Expect(session.DeleteCookiesCall.Called).To(BeTrue())
		})

		Context("when deleting all cookies fails", func() {
			It("should return an error", func() {
				session.DeleteCookiesCall.Err = errors.New("some error")
				Expect(page.ClearCookies()).To(MatchError("failed to clear cookies: some error"))
			})
		})
	})

	Describe("#URL", func() {
		It("should successfully return the URL of the current page", func() {
			session.GetURLCall.ReturnURL = "http://example.com"
			url, err := page.URL()
			Expect(url).To(Equal("http://example.com"))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the session fails to retrieve the URL", func() {
			It("should return an error", func() {
				session.GetURLCall.Err = errors.New("some error")
				_, err := page.URL()
				Expect(err).To(MatchError("failed to retrieve URL: some error"))
			})
		})
	})

	Describe("#Size", func() {
		var (
			bus    *mocks.Bus
			window *api.Window
		)

		BeforeEach(func() {
			bus = &mocks.Bus{}
			window = &api.Window{Session: &api.Session{Bus: bus}}
		})

		It("should set the window width and height to the provided dimensions", func() {
			session.GetWindowCall.ReturnWindow = window
			Expect(page.Size(640, 480)).To(Succeed())
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"width": 640, "height": 480}`))
		})

		Context("when the session fails to retrieve a window", func() {
			It("should return an error", func() {
				session.GetWindowCall.Err = errors.New("some error")
				Expect(page.Size(640, 480)).To(MatchError("failed to retrieve window: some error"))
			})
		})

		Context("when the window fails to retrieve its size", func() {
			It("should return an error", func() {
				session.GetWindowCall.ReturnWindow = window
				bus.SendCall.Err = errors.New("some error")
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

		Context("when the session fails to retrieve a screenshot", func() {
			BeforeEach(func() {
				session.GetScreenshotCall.Err = errors.New("some error")
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
				session.GetScreenshotCall.ReturnImage = []byte("some-image")
				Expect(page.Screenshot(filename)).To(Succeed())
				defer os.Remove(filename)
				result, _ := ioutil.ReadFile(filename)
				Expect(string(result)).To(Equal("some-image"))
			})
		})
	})

	Describe("#Title", func() {
		It("should successfully return the title of the current page", func() {
			session.GetTitleCall.ReturnTitle = "Some Title"
			title, err := page.Title()
			Expect(title).To(Equal("Some Title"))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the session fails to retrieve the page title", func() {
			It("should return an error", func() {
				session.GetTitleCall.Err = errors.New("some error")
				_, err := page.Title()
				Expect(err).To(MatchError("failed to retrieve page title: some error"))
			})
		})
	})

	Describe("#HTML", func() {
		It("should return the HTML of the current page", func() {
			session.GetSourceCall.ReturnSource = "Some HTML"
			html, err := page.HTML()
			Expect(html).To(Equal("Some HTML"))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the session fails to retrieve the page HTML", func() {
			It("should return an error", func() {
				session.GetSourceCall.Err = errors.New("some error")
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
			session.ExecuteCall.Result = `{"some": "result"}`
			err = page.RunScript("some javascript code", map[string]interface{}{"argument": "value"}, &result)
		})

		It("should provide the session with an argument-provided javascript function", func() {
			Expect(session.ExecuteCall.Body).To(Equal("return (function(argument) { some javascript code; }).apply(this, arguments);"))
		})

		It("should provide the session with arguments to call the provided function with", func() {
			Expect(session.ExecuteCall.Arguments).To(Equal([]interface{}{"value"}))
		})

		It("should unmarshall the returned result into the provided result interface", func() {
			Expect(result.Some).To(Equal("result"))
		})

		It("should be successful", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when running the script fails", func() {
			It("should return the session error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = page.RunScript("", map[string]interface{}{}, &result)
				Expect(err).To(MatchError("failed to run script: some error"))
			})
		})
	})

	Describe("#PopupText", func() {
		It("should return the popup text of the popup and succeed", func() {
			session.GetAlertTextCall.ReturnText = "some popup text"
			text, err := page.PopupText()
			Expect(text).To(Equal("some popup text"))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the session fails to retrieve the page popup text", func() {
			It("should return an error", func() {
				session.GetAlertTextCall.Err = errors.New("some error")
				_, err := page.PopupText()
				Expect(err).To(MatchError("failed to retrieve popup text: some error"))
			})
		})
	})

	Describe("#EnterPopupText", func() {
		It("should provide the session with the text to enter and succeed", func() {
			Expect(page.EnterPopupText("some text")).To(Succeed())
			Expect(session.SetAlertTextCall.Text).To(Equal("some text"))
		})

		Context("when the session fails to enter the page popup text", func() {
			It("should return an error", func() {
				session.SetAlertTextCall.Err = errors.New("some error")
				Expect(page.EnterPopupText("some text")).To(MatchError("failed to enter popup text: some error"))
			})
		})
	})

	Describe("#ConfirmPopup", func() {
		It("should instruct the session to confirm an alert", func() {
			Expect(page.ConfirmPopup()).To(Succeed())
			Expect(session.AcceptAlertCall.Called).To(BeTrue())
		})

		Context("when the session fails to confirm an alert", func() {
			It("should return an error", func() {
				session.AcceptAlertCall.Err = errors.New("some error")
				Expect(page.ConfirmPopup()).To(MatchError("failed to confirm popup: some error"))
			})
		})
	})

	Describe("#CancelPopup", func() {
		It("should instruct the session to cancel an alert", func() {
			Expect(page.CancelPopup()).To(Succeed())
			Expect(session.DismissAlertCall.Called).To(BeTrue())
		})

		Context("when the session fails to cancel an alert", func() {
			It("should return an error", func() {
				session.DismissAlertCall.Err = errors.New("some error")
				Expect(page.CancelPopup()).To(MatchError("failed to cancel popup: some error"))
			})
		})
	})

	Describe("#Forward", func() {
		It("should successfully instruct the session to move forward in history", func() {
			Expect(page.Forward()).To(Succeed())
			Expect(session.ForwardCall.Called).To(BeTrue())
		})

		Context("when navigating forward fails", func() {
			It("should return an error", func() {
				session.ForwardCall.Err = errors.New("some error")
				Expect(page.Forward()).To(MatchError("failed to navigate forward in history: some error"))
			})
		})
	})

	Describe("#Back", func() {
		It("should successfully instruct the session to move back in history", func() {
			Expect(page.Back()).To(Succeed())
			Expect(session.BackCall.Called).To(BeTrue())
		})

		Context("when navigating back fails", func() {
			It("should return an error", func() {
				session.BackCall.Err = errors.New("some error")
				Expect(page.Back()).To(MatchError("failed to navigate backwards in history: some error"))
			})
		})
	})

	Describe("#Refresh", func() {
		It("should successfully instruct the session to refresh", func() {
			Expect(page.Refresh()).To(Succeed())
			Expect(session.RefreshCall.Called).To(BeTrue())
		})

		Context("when refreshing the page fails", func() {
			It("should return an error", func() {
				session.RefreshCall.Err = errors.New("some error")
				Expect(page.Refresh()).To(MatchError("failed to refresh page: some error"))
			})
		})
	})

	Describe("#SwitchToParentFrame", func() {
		It("should successfully instruct the session to change focus to the parent frame", func() {
			Expect(page.SwitchToParentFrame()).To(Succeed())
			Expect(session.FrameParentCall.Called).To(BeTrue())
		})

		Context("when switching to the parent frame fails", func() {
			It("should return an error", func() {
				session.FrameParentCall.Err = errors.New("some error")
				Expect(page.SwitchToParentFrame()).To(MatchError("failed to switch to parent frame: some error"))
			})
		})
	})

	Describe("#SwitchToRootFrame", func() {
		It("should successfully instruct the session to change focus to the root frame", func() {
			session.FrameCall.Frame = &api.Element{}
			Expect(page.SwitchToRootFrame()).To(Succeed())
			Expect(session.FrameCall.Frame).To(BeNil())
		})

		Context("when switching to the root frame fails", func() {
			It("should return an error", func() {
				session.FrameCall.Err = errors.New("some error")
				Expect(page.SwitchToRootFrame()).To(MatchError("failed to switch to original page frame: some error"))
			})
		})
	})

	Describe("#SwitchToWindow", func() {
		It("should successfully instruct the session to switch to the named window", func() {
			Expect(page.SwitchToWindow("some name")).To(Succeed())
			Expect(session.SetWindowByNameCall.Name).To(Equal("some name"))
		})

		Context("when switching to the root frame fails", func() {
			It("should return an error", func() {
				session.SetWindowByNameCall.Err = errors.New("some error")
				Expect(page.SwitchToWindow("some name")).To(MatchError("failed to switch to named window: some error"))
			})
		})
	})

	Describe("#NextWindow", func() {
		BeforeEach(func() {
			firstWindow := &api.Window{ID: "first window"}
			secondWindow := &api.Window{ID: "second window"}
			thirdWindow := &api.Window{ID: "third window"}
			session.GetWindowsCall.ReturnWindows = []*api.Window{secondWindow, firstWindow, thirdWindow}
			session.GetWindowCall.ReturnWindow = firstWindow
		})

		It("should successfully instruct the session to switch to the next window in sorted order", func() {
			Expect(page.NextWindow()).To(Succeed())
			Expect(session.SetWindowCall.Window.ID).To(Equal("second window"))
		})

		Context("when retrieving the available windows fails", func() {
			It("should return an error", func() {
				session.GetWindowsCall.Err = errors.New("some error")
				Expect(page.NextWindow()).To(MatchError("failed to find available windows: some error"))
			})
		})

		Context("when retrieving the active window fails", func() {
			It("should return an error", func() {
				session.GetWindowCall.Err = errors.New("some error")
				Expect(page.NextWindow()).To(MatchError("failed to find active window: some error"))
			})
		})

		Context("when setting the active window fails", func() {
			It("should return an error", func() {
				session.SetWindowCall.Err = errors.New("some error")
				Expect(page.NextWindow()).To(MatchError("failed to change active window: some error"))
			})
		})
	})

	Describe("#CloseWindow", func() {
		It("should successfully instruct the session to close the active window", func() {
			Expect(page.CloseWindow()).To(Succeed())
			Expect(session.DeleteWindowCall.Called).To(BeTrue())
		})

		Context("when closing the active window fails", func() {
			It("should return an error", func() {
				session.DeleteWindowCall.Err = errors.New("some error")
				Expect(page.CloseWindow()).To(MatchError("failed to close active window: some error"))
			})
		})
	})

	Describe("#WindowCount", func() {
		It("should successfully return the number of windows from the session", func() {
			session.GetWindowsCall.ReturnWindows = []*api.Window{&api.Window{}, &api.Window{}}
			Expect(page.WindowCount()).To(Equal(2))
		})

		Context("when retrieving the available windows fails", func() {
			It("should return an error", func() {
				session.GetWindowsCall.Err = errors.New("some error")
				_, err := page.WindowCount()
				Expect(err).To(MatchError("failed to find available windows: some error"))
			})
		})
	})

	Describe("#ReadNewLogs", func() {
		It("should request new logs of the provided log type from the session", func() {
			_, err := page.ReadNewLogs("some type")
			Expect(err).To(Succeed())
			Expect(session.NewLogsCall.LogType).To(Equal("some type"))
		})

		Context("when the session fails to retrieve logs", func() {
			It("should return an error", func() {
				session.NewLogsCall.Err = errors.New("some error")
				_, err := page.ReadNewLogs("some type")
				Expect(err).To(MatchError("failed to retrieve logs: some error"))
			})
		})

		Context("when only new logs are requested", func() {
			It("should return only new logs with the correct time and code location", func() {
				session.NewLogsCall.ReturnLogs = []api.Log{api.Log{"old log", "old level", 1418196096123}}
				page.ReadNewLogs("some type")
				session.NewLogsCall.ReturnLogs = []api.Log{
					api.Log{"new log (1:22)", "new level", 1418196097543},
					api.Log{"newer log (:)", "newer level", 1418196098376},
				}

				logs, err := page.ReadNewLogs("some type")
				Expect(err).NotTo(HaveOccurred())
				Expect(logs).To(HaveLen(2))
				Expect(logs[0].Message).To(Equal("new log"))
				Expect(logs[0].Location).To(Equal("1:22"))
				Expect(logs[0].Level).To(Equal("new level"))
				Expect(logs[0].Time.Unix()).To(BeEquivalentTo(1418196097))
				Expect(logs[1].Message).To(Equal("newer log"))
				Expect(logs[1].Location).To(Equal(":"))
				Expect(logs[1].Level).To(Equal("newer level"))
				Expect(logs[1].Time.Unix()).To(BeEquivalentTo(1418196098))
			})
		})
	})

	Describe("#ReadAllLogs", func() {
		It("should call ReadNewLogs and return previously read logs", func() {
			session.NewLogsCall.ReturnLogs = []api.Log{api.Log{"old log", "old level", 1418196096123}}
			page.ReadNewLogs("some type")
			session.NewLogsCall.ReturnLogs = []api.Log{
				api.Log{"new log (1:22)", "new level", 1418196097543},
				api.Log{"newer log (:)", "newer level", 1418196098376},
			}

			logs, err := page.ReadAllLogs("some type")
			Expect(err).NotTo(HaveOccurred())
			Expect(logs).To(HaveLen(3))
			Expect(logs[0].Message).To(Equal("old log"))
			Expect(logs[1].Message).To(Equal("new log"))
			Expect(logs[2].Message).To(Equal("newer log"))
		})

		It("should return a copy of the stored logs", func() {
			session.NewLogsCall.ReturnLogs = []api.Log{api.Log{"some log", "some level", 1418196096123}}
			logs, _ := page.ReadAllLogs("some type")
			logs[0].Message = "some changed log"
			logs, _ = page.ReadAllLogs("some type")
			Expect(logs[0].Message).To(Equal("some log"))
		})

		Context("when Page#ReadNewLogs fails", func() {
			It("should return an error", func() {
				session.NewLogsCall.Err = errors.New("some error")
				_, err := page.ReadNewLogs("some type")
				Expect(err).To(MatchError("failed to retrieve logs: some error"))
			})
		})
	})

	Describe("#LogTypes", func() {
		It("should successfully return the log types", func() {
			session.GetLogTypesCall.ReturnTypes = []string{"first type", "second type"}
			Expect(page.LogTypes()).To(Equal([]string{"first type", "second type"}))
		})

		Context("when the session fails to retrieve the log types", func() {
			It("should return an error", func() {
				session.GetLogTypesCall.Err = errors.New("some error")
				_, err := page.LogTypes()
				Expect(err).To(MatchError("failed to retrieve log types: some error"))
			})
		})
	})
})
