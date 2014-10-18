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

	Describe("#DeleteSession", func() {
		BeforeEach(func() {
			err = client.DeleteSession()
		})

		It("should make a DELETE request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("DELETE"))
		})

		It("should hit the / endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal(""))
		})

		Context("when the sesssion indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the page failed to delete the cookies", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.DeleteSession()
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

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /elements endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("elements"))
		})

		It("should include the selection in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		Context("when the session indicates a success", func() {
			It("should return a slice of elements with IDs and sessions", func() {
				Expect(elements[0].(*element.Element).ID).To(Equal("some-id"))
				Expect(elements[0].(*element.Element).Session).To(Equal(session))
				Expect(elements[1].(*element.Element).ID).To(Equal("some-other-id"))
				Expect(elements[1].(*element.Element).Session).To(Equal(session))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = client.GetElements(types.Selector{Using: "css selector", Value: "#selector"})
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

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /window_handle endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("window_handle"))
		})

		Context("when the session indicates a success", func() {
			It("should return the window with the retrieved ID and session", func() {
				Expect(clientWindow.(*window.Window).ID).To(Equal("some-id"))
				Expect(clientWindow.(*window.Window).Session).To(Equal(session))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = client.GetWindow()
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

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /cookie endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("cookie"))
		})

		It("should include the cookie to add in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"cookie":{"name":"some-name","value":42,"path":"/my-path","domain":"example.com","secure":false,"httpOnly":false,"expiry":1412358590}}`))
		})

		Context("when the sesssion indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the page failed to add the cookie", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.SetCookie(cookie)
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		BeforeEach(func() {
			err = client.DeleteCookie("some-cookie")
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("DELETE"))
		})

		It("should hit the /cookie/:name endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("cookie/some-cookie"))
		})

		Context("when the sesssion indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the page failed to delete the cookie", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.DeleteCookie("some-cookie")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookies", func() {
		BeforeEach(func() {
			err = client.DeleteCookies()
		})

		It("should make a DELETE request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("DELETE"))
		})

		It("should hit the /cookie endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("cookie"))
		})

		Context("when the sesssion indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the page failed to delete the cookies", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.DeleteCookies()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetScreenshot", func() {
		var (
			image []byte
			err   error
		)

		BeforeEach(func() {
			session.ExecuteCall.Result = `"c29tZS1wbmc="`
			image, err = client.GetScreenshot()
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /screenshot endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("screenshot"))
		})

		Context("when the session indicates a success", func() {
			Context("when the image is valid base64", func() {
				It("should return the decoded image", func() {
					Expect(string(image)).To(Equal("some-png"))
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("and the image is not valid base64", func() {
				BeforeEach(func() {
					session.ExecuteCall.Result = `"..."`
					image, err = client.GetScreenshot()
				})

				It("should return an error", func() {
					Expect(err).To(MatchError("illegal base64 data at input byte 0"))
				})
			})
		})

		Context("when the session indicates a failure", func() {
			BeforeEach(func() {
				session.ExecuteCall.Err = errors.New("some error")
				image, err = client.GetScreenshot()
			})

			It("should return an error", func() {
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

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /url endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("url"))
		})

		Context("when the sesssion indicates a success", func() {
			It("should return the page URL", func() {
				Expect(url).To(Equal("http://example.com"))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the page failed to retrieve the URL", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = client.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetURL", func() {
		BeforeEach(func() {
			err = client.SetURL("http://example.com")
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /url endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("url"))
		})

		It("should include the new URL in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"url": "http://example.com"}`))
		})

		Context("when the sesssion indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the page failed to change URL", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.SetURL("http://example.com")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetTitle", func() {
		var title string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"Some Title"`
			title, err = client.GetTitle()
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /title endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("title"))
		})

		Context("when the sesssion indicates a success", func() {
			It("should return the page title", func() {
				Expect(title).To(Equal("Some Title"))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the page failed to retrieve the title", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = client.GetURL()
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

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /source endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("source"))
		})

		Context("when the sesssion indicates a success", func() {
			It("should return the page source", func() {
				Expect(source).To(Equal("some source"))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the page failed to retrieve the source", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = client.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DoubleClick", func() {
		BeforeEach(func() {
			err = client.DoubleClick()
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /doubleclick endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("doubleclick"))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to double-click", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.DoubleClick()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#MoveTo", func() {
		BeforeEach(func() {
			err = client.MoveTo(nil, nil)
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /moveto endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("moveto"))
		})

		It("should encode no element or point if not provided", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{}`))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to move the mouse", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.MoveTo(nil, nil)
				Expect(err).To(MatchError("some error"))
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

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /execute endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("execute"))
		})

		It("should include the javascript and arguments in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"script": "some javascript code", "args": [1, "two"]}`))
		})

		Context("when the session indicates a success", func() {
			It("should fill the provided results interface", func() {
				Expect(result.Some).To(Equal("result"))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.Execute("", nil, &result)
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Forward", func() {
		BeforeEach(func() {
			err = client.Forward()
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /forward endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("forward"))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to go forward in history", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.Forward()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Back", func() {
		BeforeEach(func() {
			err = client.Back()
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /back endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("back"))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to go back in history", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.Back()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Refresh", func() {
		BeforeEach(func() {
			err = client.Refresh()
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /refresh endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("refresh"))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to refresh the page", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = client.Refresh()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

})
