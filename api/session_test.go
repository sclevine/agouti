package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/internal/mocks"
)

var _ = Describe("Bus", func() {
	var (
		session *Session
		bus     *mocks.Bus
		err     error
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		session = &Session{bus}
	})

	ItShouldMakeARequest := func(method, endpoint string, body ...string) {
		It("should make a "+method+" request", func() {
			Expect(bus.SendCall.Method).To(Equal(method))
		})

		It("should hit the desired endpoint", func() {
			Expect(bus.SendCall.Endpoint).To(Equal(endpoint))
		})

		It("should not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		if len(body) > 0 {
			It("should set the request body", func() {
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(body[0]))
			})
		}
	}

	Describe("#Delete", func() {
		BeforeEach(func() {
			err = session.Delete()
		})

		ItShouldMakeARequest("DELETE", "")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Delete()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElement", func() {
		var element *Element

		BeforeEach(func() {
			bus.SendCall.Result = `{"ELEMENT": "some-id"}`
			element, err = session.GetElement(Selector{"css selector", "#selector"})
		})

		ItShouldMakeARequest("POST", "element", `{"using": "css selector", "value": "#selector"}`)

		It("should return an element with an ID and session", func() {
			Expect(element.ID).To(Equal("some-id"))
			Expect(element.Session).To(Equal(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetElement(Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		var elements []*Element

		BeforeEach(func() {
			bus.SendCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = session.GetElements(Selector{"css selector", "#selector"})
		})

		ItShouldMakeARequest("POST", "elements", `{"using": "css selector", "value": "#selector"}`)

		It("should return a slice of elements with IDs and sessions", func() {
			Expect(elements[0].ID).To(Equal("some-id"))
			Expect(elements[0].Session).To(Equal(session))
			Expect(elements[1].ID).To(Equal("some-other-id"))
			Expect(elements[1].Session).To(Equal(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetElements(Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetActiveElement", func() {
		var element *Element

		BeforeEach(func() {
			bus.SendCall.Result = `{"ELEMENT": "some-id"}`
			element, err = session.GetActiveElement()
		})

		ItShouldMakeARequest("POST", "element/active")

		It("should return the active element with an ID and session", func() {
			Expect(element.ID).To(Equal("some-id"))
			Expect(element.Session).To(Equal(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetActiveElement()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetWindow", func() {
		var window *Window

		BeforeEach(func() {
			bus.SendCall.Result = `"some-id"`
			window, err = session.GetWindow()
		})

		ItShouldMakeARequest("GET", "window_handle")

		It("should return the current window with the retrieved ID and session", func() {
			Expect(window.ID).To(Equal("some-id"))
			Expect(window.Session).To(Equal(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetWindow()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetWindows", func() {
		var windows []*Window

		BeforeEach(func() {
			bus.SendCall.Result = `["some-id", "some-other-id"]`
			windows, err = session.GetWindows()
		})

		ItShouldMakeARequest("GET", "window_handles")

		It("should return all windows with their retrieved IDs and sessions", func() {
			Expect(windows[0].ID).To(Equal("some-id"))
			Expect(windows[0].Session).To(Equal(session))
			Expect(windows[1].ID).To(Equal("some-other-id"))
			Expect(windows[1].Session).To(Equal(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetWindows()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetWindow", func() {
		var window *Window

		BeforeEach(func() {
			window = &Window{ID: "some-id"}
			err = session.SetWindow(window)
		})

		ItShouldMakeARequest("POST", "window", `{"name": "some-id"}`)

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetWindow(window)).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetWindowByName", func() {
		BeforeEach(func() {
			err = session.SetWindowByName("some name")
		})

		ItShouldMakeARequest("POST", "window", `{"name": "some name"}`)

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetWindowByName("some name")).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteWindow", func() {
		BeforeEach(func() {
			err = session.DeleteWindow()
		})

		ItShouldMakeARequest("DELETE", "window")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				err := session.DeleteWindow()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		var cookie map[string]interface{}

		BeforeEach(func() {
			cookie = map[string]interface{}{"name": "some-cookie"}
			err = session.SetCookie(cookie)
		})

		ItShouldMakeARequest("POST", "cookie", `{"cookie": {"name": "some-cookie"}}`)

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetCookie(cookie)).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		BeforeEach(func() {
			err = session.DeleteCookie("some-cookie")
		})

		ItShouldMakeARequest("DELETE", "cookie/some-cookie")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DeleteCookie("some-cookie")).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookies", func() {
		BeforeEach(func() {
			err = session.DeleteCookies()
		})

		ItShouldMakeARequest("DELETE", "cookie")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DeleteCookies()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetScreenshot", func() {
		var image []byte

		BeforeEach(func() {
			bus.SendCall.Result = `"c29tZS1wbmc="`
			image, err = session.GetScreenshot()
		})

		ItShouldMakeARequest("GET", "screenshot")

		Context("when the image is valid base64", func() {
			It("should return the decoded image", func() {
				Expect(string(image)).To(Equal("some-png"))
			})
		})

		Context("when the image is not valid base64", func() {
			It("should return an error", func() {
				bus.SendCall.Result = `"..."`
				_, err := session.GetScreenshot()
				Expect(err).To(MatchError("illegal base64 data at input byte 0"))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetScreenshot()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetURL", func() {
		var url string

		BeforeEach(func() {
			bus.SendCall.Result = `"http://example.com"`
			url, err = session.GetURL()
		})

		ItShouldMakeARequest("GET", "url")

		It("should return the page URL", func() {
			Expect(url).To(Equal("http://example.com"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetURL", func() {
		BeforeEach(func() {
			err = session.SetURL("http://example.com")
		})

		ItShouldMakeARequest("POST", "url", `{"url": "http://example.com"}`)

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetURL("http://example.com")).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetTitle", func() {
		var title string

		BeforeEach(func() {
			bus.SendCall.Result = `"Some Title"`
			title, err = session.GetTitle()
		})

		ItShouldMakeARequest("GET", "title")

		It("should return the page title", func() {
			Expect(title).To(Equal("Some Title"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err = session.GetTitle()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetSource", func() {
		var source string

		BeforeEach(func() {
			bus.SendCall.Result = `"some source"`
			source, err = session.GetSource()
		})

		ItShouldMakeARequest("GET", "source")

		It("should return the page source", func() {
			Expect(source).To(Equal("some source"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetSource()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DoubleClick", func() {
		BeforeEach(func() {
			err = session.DoubleClick()
		})

		ItShouldMakeARequest("POST", "doubleclick")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DoubleClick()).To(MatchError("some error"))
			})
		})
	})

	Describe("#MoveTo", func() {
		BeforeEach(func() {
			err = session.MoveTo(nil, nil)
		})

		ItShouldMakeARequest("POST", "moveto", "{}")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.MoveTo(nil, nil)).To(MatchError("some error"))
			})
		})

		Context("when an element is provided", func() {
			It("should encode the element into the request JSON", func() {
				element := &Element{ID: "some-id"}
				session.MoveTo(element, nil)
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"element": "some-id"}`))
			})
		})

		Context("when a X offset is provided", func() {
			It("should encode the element into the request JSON", func() {
				session.MoveTo(nil, XOffset(100))
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"xoffset": 100}`))
			})
		})

		Context("when a Y offset is provided", func() {
			It("should encode the element into the request JSON", func() {
				session.MoveTo(nil, YOffset(200))
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"yoffset": 200}`))
			})
		})

		Context("when an XY offset is provided", func() {
			It("should encode the element into the request JSON", func() {
				session.MoveTo(nil, XYOffset{X: 300, Y: 400})
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"xoffset": 300, "yoffset": 400}`))
			})
		})
	})

	Describe("#Frame", func() {
		BeforeEach(func() {
			err = session.Frame(&Element{ID: "some-id"})
		})

		ItShouldMakeARequest("POST", "frame", `{"id": {"ELEMENT": "some-id"}}`)

		Context("When the provided frame in nil", func() {
			BeforeEach(func() {
				err = session.Frame(nil)
			})

			ItShouldMakeARequest("POST", "frame", `{"id": null}`)
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Frame(&Element{ID: "some-id"})).To(MatchError("some error"))
			})
		})
	})

	Describe("#FrameParent", func() {
		BeforeEach(func() {
			err = session.FrameParent()
		})

		ItShouldMakeARequest("POST", "frame/parent")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.FrameParent()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Execute", func() {
		var (
			result struct{ Some string }
			err    error
		)

		BeforeEach(func() {
			bus.SendCall.Result = `{"some": "result"}`
			err = session.Execute("some javascript code", []interface{}{1, "two"}, &result)
		})

		ItShouldMakeARequest("POST", "execute", `{"script": "some javascript code", "args": [1, "two"]}`)

		It("should fill the provided results interface", func() {
			Expect(result.Some).To(Equal("result"))
		})

		Context("when called with nil arguments", func() {
			It("should send an empty list for args", func() {
				session.Execute("some javascript code", nil, nil)
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"script": "some javascript code", "args": []}`))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Execute("", nil, &result)).To(MatchError("some error"))
			})
		})
	})

	Describe("#Forward", func() {
		BeforeEach(func() {
			err = session.Forward()
		})

		ItShouldMakeARequest("POST", "forward")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Forward()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Back", func() {
		BeforeEach(func() {
			err = session.Back()
		})

		ItShouldMakeARequest("POST", "back")

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to go back in history", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Back()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Refresh", func() {
		BeforeEach(func() {
			err = session.Refresh()
		})

		ItShouldMakeARequest("POST", "refresh")

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to refresh the page", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Refresh()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetAlertText", func() {
		var text string

		BeforeEach(func() {
			bus.SendCall.Result = `"some text"`
			text, err = session.GetAlertText()
		})

		ItShouldMakeARequest("GET", "alert_text")

		It("should return the current alert text", func() {
			Expect(text).To(Equal("some text"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetAlertText()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetAlertText", func() {
		BeforeEach(func() {
			err = session.SetAlertText("some text")
		})

		ItShouldMakeARequest("POST", "alert_text", `{"text": "some text"}`)

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetAlertText("some text")).To(MatchError("some error"))
			})
		})
	})

	Describe("#AcceptAlert", func() {
		BeforeEach(func() {
			err = session.AcceptAlert()
		})

		ItShouldMakeARequest("POST", "accept_alert")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.AcceptAlert()).To(MatchError("some error"))
			})
		})
	})

	Describe("#DismissAlert", func() {
		BeforeEach(func() {
			err = session.DismissAlert()
		})

		ItShouldMakeARequest("POST", "dismiss_alert")

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DismissAlert()).To(MatchError("some error"))
			})
		})
	})

	Describe("#NewLogs", func() {
		var logs []Log

		BeforeEach(func() {
			bus.SendCall.Result = `[
				{"message": "some message", "level": "INFO", "timestamp": 1417988844498},
				{"message": "some other message", "level": "WARNING", "timestamp": 1417982864598}
			]`
			logs, err = session.NewLogs("browser")
		})

		ItShouldMakeARequest("POST", "log", `{"type": "browser"}`)

		It("should return all logs", func() {
			Expect(logs[0].Message).To(Equal("some message"))
			Expect(logs[0].Level).To(Equal("INFO"))
			Expect(logs[0].Timestamp).To(BeEquivalentTo(1417988844498))
			Expect(logs[1].Message).To(Equal("some other message"))
			Expect(logs[1].Level).To(Equal("WARNING"))
			Expect(logs[1].Timestamp).To(BeEquivalentTo(1417982864598))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.NewLogs("browser")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetLogTypes", func() {
		var types []string

		BeforeEach(func() {
			bus.SendCall.Result = `["first type", "second type"]`
			types, err = session.GetLogTypes()
		})

		ItShouldMakeARequest("GET", "log/types")

		It("should return the current alert text", func() {
			Expect(types).To(Equal([]string{"first type", "second type"}))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetLogTypes()
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
