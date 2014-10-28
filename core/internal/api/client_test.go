package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/api/element"
	"github.com/sclevine/agouti/core/internal/api/window"
	"github.com/sclevine/agouti/core/internal/mocks"
	"github.com/sclevine/agouti/core/internal/types"
)

var _ = Describe("API Client", func() {
	var (
		client  *Client
		session *mocks.Session
		err     error
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		client = &Client{session}
	})

	ItShouldMakeARequest := func(method, endpoint string, body ...string) {
		It("should make a "+method+" request", func() {
			ExpectWithOffset(1, session.ExecuteCall.Method).To(Equal(method))
		})

		It("should hit the desired endpoint", func() {
			ExpectWithOffset(1, session.ExecuteCall.Endpoint).To(Equal(endpoint))
		})

		It("should not return an error", func() {
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
		})

		if len(body) > 0 {
			It("should set the request body", func() {
				ExpectWithOffset(1, session.ExecuteCall.BodyJSON).To(MatchJSON(body[0]))
			})
		}
	}

	Describe("#DeleteSession", func() {
		BeforeEach(func() {
			err = client.DeleteSession()
		})

		ItShouldMakeARequest("DELETE", "")

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.DeleteSession()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElement", func() {
		var singleElement types.Element

		BeforeEach(func() {
			session.ExecuteCall.Result = `{"ELEMENT": "some-id"}`
			singleElement, err = client.GetElement(types.Selector{Using: "css selector", Value: "#selector"})
		})

		ItShouldMakeARequest("POST", "element", `{"using": "css selector", "value": "#selector"}`)

		It("should return an element with an ID and session", func() {
			Expect(singleElement.(*element.Element).ID).To(Equal("some-id"))
			Expect(singleElement.(*element.Element).Session).To(Equal(session))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := client.GetElement(types.Selector{Using: "css selector", Value: "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		var elements []types.Element

		BeforeEach(func() {
			session.ExecuteCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = client.GetElements(types.Selector{Using: "css selector", Value: "#selector"})
		})

		ItShouldMakeARequest("POST", "elements", `{"using": "css selector", "value": "#selector"}`)

		It("should return a slice of elements with IDs and sessions", func() {
			Expect(elements[0].(*element.Element).ID).To(Equal("some-id"))
			Expect(elements[0].(*element.Element).Session).To(Equal(session))
			Expect(elements[1].(*element.Element).ID).To(Equal("some-other-id"))
			Expect(elements[1].(*element.Element).Session).To(Equal(session))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := client.GetElements(types.Selector{Using: "css selector", Value: "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetWindow", func() {
		var clientWindow types.Window

		BeforeEach(func() {
			session.ExecuteCall.Result = `"some-id"`
			clientWindow, err = client.GetWindow()
		})

		ItShouldMakeARequest("GET", "window_handle")

		It("should return the window with the retrieved ID and session", func() {
			Expect(clientWindow.(*window.Window).ID).To(Equal("some-id"))
			Expect(clientWindow.(*window.Window).Session).To(Equal(session))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := client.GetWindow()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		var cookie *types.Cookie

		BeforeEach(func() {
			cookie = &types.Cookie{
				Name:     "some-name",
				Value:    42,
				Path:     "/my-path",
				Domain:   "example.com",
				Secure:   false,
				HTTPOnly: false,
				Expiry:   1412358590,
			}
			err = client.SetCookie(cookie)
		})

		ItShouldMakeARequest("POST", "cookie", `{"cookie":{"name":"some-name","value":42,"path":"/my-path","domain":"example.com","secure":false,"httpOnly":false,"expiry":1412358590}}`)

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.SetCookie(cookie)).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		BeforeEach(func() {
			err = client.DeleteCookie("some-cookie")
		})

		ItShouldMakeARequest("DELETE", "cookie/some-cookie")

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.DeleteCookie("some-cookie")).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookies", func() {
		BeforeEach(func() {
			err = client.DeleteCookies()
		})

		ItShouldMakeARequest("DELETE", "cookie")

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.DeleteCookies()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetScreenshot", func() {
		var image []byte

		BeforeEach(func() {
			session.ExecuteCall.Result = `"c29tZS1wbmc="`
			image, err = client.GetScreenshot()
		})

		ItShouldMakeARequest("GET", "screenshot")

		Context("when the image is valid base64", func() {
			It("should return the decoded image", func() {
				Expect(string(image)).To(Equal("some-png"))
			})
		})

		Context("when the image is not valid base64", func() {
			It("should return an error", func() {
				session.ExecuteCall.Result = `"..."`
				_, err := client.GetScreenshot()
				Expect(err).To(MatchError("illegal base64 data at input byte 0"))
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := client.GetScreenshot()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetURL", func() {
		var url string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"http://example.com"`
			url, err = client.GetURL()
		})

		ItShouldMakeARequest("GET", "url")

		It("should return the page URL", func() {
			Expect(url).To(Equal("http://example.com"))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := client.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetURL", func() {
		BeforeEach(func() {
			err = client.SetURL("http://example.com")
		})

		ItShouldMakeARequest("POST", "url", `{"url": "http://example.com"}`)

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.SetURL("http://example.com")).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetTitle", func() {
		var title string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"Some Title"`
			title, err = client.GetTitle()
		})

		ItShouldMakeARequest("GET", "title")

		It("should return the page title", func() {
			Expect(title).To(Equal("Some Title"))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = client.GetTitle()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetSource", func() {
		var source string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"some source"`
			source, err = client.GetSource()
		})

		ItShouldMakeARequest("GET", "source")

		It("should return the page source", func() {
			Expect(source).To(Equal("some source"))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := client.GetSource()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DoubleClick", func() {
		BeforeEach(func() {
			err = client.DoubleClick()
		})

		ItShouldMakeARequest("POST", "doubleclick")

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.DoubleClick()).To(MatchError("some error"))
			})
		})
	})

	Describe("#MoveTo", func() {
		BeforeEach(func() {
			err = client.MoveTo(nil, nil)
		})

		ItShouldMakeARequest("POST", "moveto", "{}")

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.MoveTo(nil, nil)).To(MatchError("some error"))
			})
		})

		Context("when an element is provided", func() {
			It("should encode the element into the request JSON", func() {
				element := &mocks.Element{}
				element.GetIDCall.ReturnID = "some-id"
				client.MoveTo(element, nil)
				Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"element": "some-id"}`))
			})
		})

		Context("when a X point is provided", func() {
			It("should encode the element into the request JSON", func() {
				client.MoveTo(nil, types.XPoint(100))
				Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"xoffset": 100}`))
			})
		})

		Context("when a Y point is provided", func() {
			It("should encode the element into the request JSON", func() {
				client.MoveTo(nil, types.YPoint(200))
				Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"yoffset": 200}`))
			})
		})

		Context("when an XY point is provided", func() {
			It("should encode the element into the request JSON", func() {
				client.MoveTo(nil, types.XYPoint{XPos: 300, YPos: 400})
				Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"xoffset": 300, "yoffset": 400}`))
			})
		})
	})

	Describe("#Execute", func() {
		var (
			result struct{ Some string }
			err    error
		)

		BeforeEach(func() {
			session.ExecuteCall.Result = `{"some": "result"}`
			err = client.Execute("some javascript code", []interface{}{1, "two"}, &result)
		})

		ItShouldMakeARequest("POST", "execute", `{"script": "some javascript code", "args": [1, "two"]}`)

		It("should fill the provided results interface", func() {
			Expect(result.Some).To(Equal("result"))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.Execute("", nil, &result)).To(MatchError("some error"))
			})
		})
	})

	Describe("#Forward", func() {
		BeforeEach(func() {
			err = client.Forward()
		})

		ItShouldMakeARequest("POST", "forward")

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.Forward()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Back", func() {
		BeforeEach(func() {
			err = client.Back()
		})

		ItShouldMakeARequest("POST", "back")

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to go back in history", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.Back()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Refresh", func() {
		BeforeEach(func() {
			err = client.Refresh()
		})

		ItShouldMakeARequest("POST", "refresh")

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to refresh the page", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.Refresh()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetAlertText", func() {
		var text string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"some text"`
			text, err = client.GetAlertText()
		})

		ItShouldMakeARequest("GET", "alert_text")

		It("should return the current alert text", func() {
			Expect(text).To(Equal("some text"))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := client.GetAlertText()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetAlertText", func() {
		BeforeEach(func() {
			err = client.SetAlertText("some text")
		})

		ItShouldMakeARequest("POST", "alert_text", `{"text": "some text"}`)

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(client.SetAlertText("some text")).To(MatchError("some error"))
			})
		})
	})
})
