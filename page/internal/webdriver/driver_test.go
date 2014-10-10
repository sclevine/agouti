package webdriver_test

import (
	. "github.com/sclevine/agouti/page/internal/webdriver"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/internal/mocks"
	"github.com/sclevine/agouti/page/internal/webdriver/element"
	"io"
	"image/png"
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

	Describe("#Navigate", func() {
		BeforeEach(func() {
			err = driver.Navigate("http://example.com")
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("POST"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.Endpoint).To(Equal("url"))
		})

		It("includes the new URL in the request body", func() {
			Expect(session.BodyJSON).To(MatchJSON(`{"url": "http://example.com"}`))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to navigate", func() {
				session.Err = errors.New("some error")
				err = driver.Navigate("http://example.com")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		var elements []Element

		BeforeEach(func() {
			session.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = driver.GetElements("#selector")
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("POST"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.Endpoint).To(Equal("elements"))
		})

		It("includes the new URL in the request body", func() {
			Expect(session.BodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
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
				session.Err = errors.New("some error")
				_, err = driver.GetElements("#selector")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetWindow", func() {
		var myWindow Window

		BeforeEach(func() {
			session.Result = `"a window"`
			myWindow, err = driver.GetWindow()
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("GET"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.Endpoint).To(Equal("window_handle"))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the elements", func() {
				session.Err = errors.New("some error")
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
			Expect(session.Method).To(Equal("POST"))
		})

		It("hits the /cookie endpoint", func() {
			Expect(session.Endpoint).To(Equal("cookie"))
		})

		It("includes the cookie to add in the request body", func() {
			Expect(session.BodyJSON).To(MatchJSON(`{"cookie":{"name":"some-name","value":42,"path":"/my-path","domain":"example.com","secure":false,"httpOnly":false,"expiry":1412358590}}`))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to add the cookie", func() {
				session.Err = errors.New("some error")
				err = driver.SetCookie(cookie)
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteAllCookies", func() {
		BeforeEach(func() {
			err = driver.DeleteAllCookies()
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("DELETE"))
		})

		It("hits the /cookie endpoint", func() {
			Expect(session.Endpoint).To(Equal("cookie"))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to add the cookie", func() {
				session.Err = errors.New("some error")
				err = driver.DeleteAllCookies()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		BeforeEach(func() {
			err = driver.DeleteCookie("myCookie")
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("DELETE"))
		})

		It("hits the /cookie/:name endpoint", func() {
			Expect(session.Endpoint).To(Equal("cookie/myCookie"))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to add the cookie", func() {
				session.Err = errors.New("some error")
				err = driver.DeleteCookie("myCookie")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Screenshot", func() {
		var reader io.Reader

		BeforeEach(func() {
			session.Result = `"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVQYV2P4DwABAQEAWk1v8QAAAABJRU5ErkJggg=="`
			reader, err = driver.Screenshot()
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("GET"))
		})

		It("hits the /screenshot endpoint", func() {
			Expect(session.Endpoint).To(Equal("screenshot"))
		})

		Context("when the session indicates a success", func() {
			Context("and the image is a real PNG", func() {
				It("returns a PNG", func() {
					_, err = png.Decode(reader)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("and the image is not a PNG", func() {
				BeforeEach(func() {
					session.Result = `"..."`
					reader, err = driver.Screenshot()

				})

				It("returns an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("when the session indicates a failure", func() {
			BeforeEach(func() {
				session.Err = errors.New("some error")
				reader, err = driver.Screenshot()
			})

			It("returns an error", func() {
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetURL", func() {
		var url string

		BeforeEach(func() {
			session.Result = `"http://example.com"`
			url, err = driver.GetURL()
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("GET"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.Endpoint).To(Equal("url"))
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
			It("returns an error indicating the page failed to add the cookie", func() {
				session.Err = errors.New("some error")
				_, err = driver.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
