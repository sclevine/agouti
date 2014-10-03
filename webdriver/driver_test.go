package webdriver_test

import (
	. "github.com/sclevine/agouti/webdriver"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webdriver", func() {
	var (
		driver  *Driver
		session *mockSession
		err     error
	)

	BeforeEach(func() {
		session = &mockSession{}
		driver = &Driver{session}
	})

	Describe("#Navigate", func() {
		BeforeEach(func() {
			err = driver.Navigate("http://example.com")
		})

		It("makes a POST request", func() {
			Expect(session.method).To(Equal("POST"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.endpoint).To(Equal("url"))
		})

		It("includes the new URL in the request body", func() {
			Expect(session.bodyJSON).To(MatchJSON(`{"url": "http://example.com"}`))
		})

		Context("when the sesssion indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the page failed to navigate", func() {
				session.err = errors.New("some error")
				err = driver.Navigate("http://example.com")
				Expect(err).To(MatchError("failed to navigate: some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		var elements []*Element

		BeforeEach(func() {
			session.result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = driver.GetElements("#selector")
		})

		It("makes a POST request", func() {
			Expect(session.method).To(Equal("POST"))
		})

		It("hits the /url endpoint", func() {
			Expect(session.endpoint).To(Equal("elements"))
		})

		It("includes the new URL in the request body", func() {
			Expect(session.bodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		Context("when the session indicates a success", func() {
			It("returns a slice of elements with IDs and sessions", func() {
				Expect(elements[0].ID).To(Equal("some-id"))
				Expect(elements[0].Session).To(Equal(session))
				Expect(elements[1].ID).To(Equal("some-other-id"))
				Expect(elements[1].Session).To(Equal(session))
			})

			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the elements", func() {
				session.err = errors.New("some error")
				_, err = driver.GetElements("#selector")
				Expect(err).To(MatchError("failed to get elements with selector '#selector': some error"))
			})
		})
	})
})
