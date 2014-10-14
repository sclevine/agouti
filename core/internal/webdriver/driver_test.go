package webdriver_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/webdriver"
	"github.com/sclevine/agouti/core/internal/webdriver/element"
	"github.com/sclevine/agouti/core/internal/webdriver/window"
	"github.com/sclevine/agouti/internal/mocks"
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

	Describe("#GetElements", func() {
		var elements []Element

		BeforeEach(func() {
			session.ExecuteCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = driver.GetElements("#selector")
		})

		It("makes a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("elements"))
		})

		It("includes the new URL in the request body", func() {
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
				_, err = driver.GetElements("#selector")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetWindow", func() {
		var driverWindow Window

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
		var cookie *Cookie

		BeforeEach(func() {
			cookie = &Cookie{
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
})
