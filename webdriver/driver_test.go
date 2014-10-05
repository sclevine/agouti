package webdriver_test

import (
	. "github.com/sclevine/agouti/webdriver"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/mocks"
	"github.com/sclevine/agouti/webdriver/element"
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
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to navigate", func() {
				session.Err = errors.New("some error")
				err = driver.Navigate("http://example.com")
				Expect(err).To(MatchError("failed to navigate: some error"))
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
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the elements", func() {
				session.Err = errors.New("some error")
				_, err = driver.GetElements("#selector")
				Expect(err).To(MatchError("failed to get elements with selector '#selector': some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		var cookie *Cookie

		BeforeEach(func() {
			cookie = &Cookie{
				Name: "theName",
				Value: 42,
				Path: "/my-path",
				Domain: "example.com",
				Secure: false,
				HTTPOnly: false,
				Expiry: 1412358590,
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
			Expect(session.BodyJSON).To(MatchJSON(`{"cookie":{"name":"theName","value":42,"path":"/my-path","domain":"example.com","secure":false,"httpOnly":false,"expiry":1412358590}}`))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to add the cookie", func() {
				session.Err = errors.New("some error")
				err = driver.SetCookie(cookie)
				Expect(err).To(MatchError("failed to add cookie: some error"))
			})
		})
	})
})
