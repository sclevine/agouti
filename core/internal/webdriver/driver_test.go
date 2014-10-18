package webdriver_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	"github.com/sclevine/agouti/core/internal/types"
	. "github.com/sclevine/agouti/core/internal/webdriver"
	"github.com/sclevine/agouti/core/internal/webdriver/element"
	"github.com/sclevine/agouti/core/internal/webdriver/window"
)

var _ = Describe("Webdriver", func() {
	var (
		driver  *Driver
		session *mocks.Session
		err     error
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		driver = &Driver{session}
	})

	Describe("#DeleteSession", func() {
		BeforeEach(func() {
			err = driver.DeleteSession()
		})

		It("makes a DELETE request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("DELETE"))
		})

		It("hits the / endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal(""))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to delete the cookies", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.DeleteSession()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		var elements []types.Element

		BeforeEach(func() {
			session.ExecuteCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = driver.GetElements(types.Selector{Using: "css selector", Value: "#selector"})
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /elements endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("elements"))
		})

		It("includes the selection in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		Context("when the session indicates a success", func() {
			It("returns a slice of elements with IDs and sessions", func() {
				Expect(elements[0].(*element.Element).ID).To(Equal("some-id"))
				Expect(elements[0].(*element.Element).Session).To(Equal(session))
				Expect(elements[1].(*element.Element).ID).To(Equal("some-other-id"))
				Expect(elements[1].(*element.Element).Session).To(Equal(session))
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = driver.GetElements(types.Selector{Using: "css selector", Value: "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetWindow", func() {
		var driverWindow types.Window

		BeforeEach(func() {
			session.ExecuteCall.Result = `"some-id"`
			driverWindow, err = driver.GetWindow()
		})

		It("makes a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("hits the /window_handle endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("window_handle"))
		})

		Context("when the session indicates a success", func() {
			It("returns the window with the retrieved ID and session", func() {
				Expect(driverWindow.(*window.Window).ID).To(Equal("some-id"))
				Expect(driverWindow.(*window.Window).Session).To(Equal(session))
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = driver.GetWindow()
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

			err = driver.SetCookie(cookie)
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /cookie endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("cookie"))
		})

		It("includes the cookie to add in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"cookie":{"name":"some-name","value":42,"path":"/my-path","domain":"example.com","secure":false,"httpOnly":false,"expiry":1412358590}}`))
		})

		Context("when the sesssion indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to add the cookie", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.SetCookie(cookie)
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		BeforeEach(func() {
			err = driver.DeleteCookie("some-cookie")
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("DELETE"))
		})

		It("hits the /cookie/:name endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("cookie/some-cookie"))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to delete the cookie", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.DeleteCookie("some-cookie")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookies", func() {
		BeforeEach(func() {
			err = driver.DeleteCookies()
		})

		It("makes a DELETE request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("DELETE"))
		})

		It("hits the /cookie endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("cookie"))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to delete the cookies", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.DeleteCookies()
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
			image, err = driver.GetScreenshot()
		})

		It("makes a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("hits the /screenshot endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("screenshot"))
		})

		Context("when the session indicates a success", func() {
			Context("when the image is valid base64", func() {
				It("returns the decoded image", func() {
					Expect(string(image)).To(Equal("some-png"))
				})

				It("does not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("and the image is not valid base64", func() {
				BeforeEach(func() {
					session.ExecuteCall.Result = `"..."`
					image, err = driver.GetScreenshot()
				})

				It("returns an error", func() {
					Expect(err).To(MatchError("illegal base64 data at input byte 0"))
				})
			})
		})

		Context("when the session indicates a failure", func() {
			BeforeEach(func() {
				session.ExecuteCall.Err = errors.New("some error")
				image, err = driver.GetScreenshot()
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetURL", func() {
		var url string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"http://example.com"`
			url, err = driver.GetURL()
		})

		It("makes a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("url"))
		})

		Context("when the sesssion indicates a success", func() {
			It("returns the page URL", func() {
				Expect(url).To(Equal("http://example.com"))
			})

			It("doesn't return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to retrieve the URL", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = driver.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetURL", func() {
		BeforeEach(func() {
			err = driver.SetURL("http://example.com")
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("url"))
		})

		It("includes the new URL in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"url": "http://example.com"}`))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to change URL", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.SetURL("http://example.com")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetTitle", func() {
		var title string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"Some Title"`
			title, err = driver.GetTitle()
		})

		It("makes a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("hits the /title endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("title"))
		})

		Context("when the sesssion indicates a success", func() {
			It("returns the page title", func() {
				Expect(title).To(Equal("Some Title"))
			})

			It("doesn't return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to retrieve the title", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = driver.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetSource", func() {
		var source string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"some source"`
			source, err = driver.GetSource()
		})

		It("makes a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("hits the /source endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("source"))
		})

		Context("when the sesssion indicates a success", func() {
			It("returns the page source", func() {
				Expect(source).To(Equal("some source"))
			})

			It("doesn't return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to retrieve the source", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = driver.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DoubleClick", func() {
		BeforeEach(func() {
			err = driver.DoubleClick()
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /doubleclick endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("doubleclick"))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to double-click", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.DoubleClick()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#MoveTo", func() {
		BeforeEach(func() {
			err = driver.MoveTo(nil, nil)
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /moveto endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("moveto"))
		})

		It("encodes no element or point if not provided", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{}`))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to move the mouse", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.MoveTo(nil, nil)
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when an element is provided", func() {
			It("encodes the element into the request JSON", func() {
				element := &mocks.Element{}
				element.GetIDCall.ReturnID = "some-id"
				driver.MoveTo(element, nil)
				Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"element": "some-id"}`))
			})
		})

		Context("when a X point is provided", func() {
			It("encodes the element into the request JSON", func() {
				driver.MoveTo(nil, types.XPoint(100))
				Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"xoffset": 100}`))
			})
		})

		Context("when a Y point is provided", func() {
			It("encodes the element into the request JSON", func() {
				driver.MoveTo(nil, types.YPoint(200))
				Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"yoffset": 200}`))
			})
		})

		Context("when an XY point is provided", func() {
			It("encodes the element into the request JSON", func() {
				driver.MoveTo(nil, types.XYPoint{XPos: 300, YPos: 400})
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
			err = driver.Execute("some javascript code", []interface{}{1, "two"}, &result)
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /execute endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("execute"))
		})

		It("includes the javascript and arguments in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"script": "some javascript code", "args": [1, "two"]}`))
		})

		Context("when the session indicates a success", func() {
			It("fills the provided results interface", func() {
				Expect(result.Some).To(Equal("result"))
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.Execute("", nil, &result)
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Forward", func() {
		BeforeEach(func() {
			err = driver.Forward()
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /forward endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("forward"))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to go forward in history", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.Forward()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Back", func() {
		BeforeEach(func() {
			err = driver.Back()
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /back endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("back"))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to go back in history", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.Back()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Refresh", func() {
		BeforeEach(func() {
			err = driver.Refresh()
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /refresh endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("refresh"))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to refresh the page", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = driver.Refresh()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

})
