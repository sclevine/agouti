package webdriver_test

import (
	. "github.com/sclevine/agouti/webdriver"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webdriver", func() {
	var (
		element *Element
		session *mockSession
		err     error
	)

	BeforeEach(func() {
		session = &mockSession{}
		element = &Element{"some-id", session}
	})

	Describe("#GetText", func() {
		var text string

		BeforeEach(func() {
			session.result = `"some text"`
			text, err = element.GetText()
		})

		It("makes a GET request", func() {
			Expect(session.method).To(Equal("GET"))
		})

		It("hits the /element/:id/text endpoint", func() {
			Expect(session.endpoint).To(Equal("element/some-id/text"))
		})

		Context("when the session indicates a success", func() {
			It("returns the visible text on the element", func() {
				Expect(text).To(Equal("some text"))
			})

			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the text", func() {
				session.err = errors.New("some error")
				_, err = element.GetText()
				Expect(err).To(MatchError("failed to retrieve text: some error"))
			})
		})
	})
})
