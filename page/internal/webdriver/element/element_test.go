package element_test

import (
	. "github.com/sclevine/agouti/page/internal/webdriver/element"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/internal/mocks"
)

var _ = Describe("Element", func() {
	var (
		element *Element
		session *mocks.Session
		err     error
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		element = &Element{"some-id", session}
	})

	Describe("#GetText", func() {
		var text string

		BeforeEach(func() {
			session.Result = `"some text"`
			text, err = element.GetText()
		})

		It("makes a GET request", func() {
			Expect(session.Method).To(Equal("GET"))
		})

		It("hits the /element/:id/text endpoint", func() {
			Expect(session.Endpoint).To(Equal("element/some-id/text"))
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
				session.Err = errors.New("some error")
				_, err = element.GetText()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetAttribute", func() {
		var value string

		BeforeEach(func() {
			session.Result = `"some value"`
			value, err = element.GetAttribute("some-name")
		})

		It("makes a GET request", func() {
			Expect(session.Method).To(Equal("GET"))
		})

		It("hits the /element/:id/attribute/:name endpoint", func() {
			Expect(session.Endpoint).To(Equal("element/some-id/attribute/some-name"))
		})

		Context("when the session indicates a success", func() {
			It("returns the value of the attribute", func() {
				Expect(value).To(Equal("some value"))
			})

			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the text", func() {
				session.Err = errors.New("some error")
				_, err = element.GetAttribute("some-name")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Click", func() {
		BeforeEach(func() {
			err = element.Click()
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("POST"))
		})

		It("hits the /element/:id/click endpoint", func() {
			Expect(session.Endpoint).To(Equal("element/some-id/click"))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the text", func() {
				session.Err = errors.New("some error")
				err = element.Click()
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
