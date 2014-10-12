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
				Expect(err).NotTo(HaveOccurred())
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
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the attribute", func() {
				session.Err = errors.New("some error")
				_, err = element.GetAttribute("some-name")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetCSS", func() {
		var value string

		BeforeEach(func() {
			session.Result = `"some value"`
			value, err = element.GetCSS("some-property")
		})

		It("makes a GET request", func() {
			Expect(session.Method).To(Equal("GET"))
		})

		It("hits the /element/:id/css/:name endpoint", func() {
			Expect(session.Endpoint).To(Equal("element/some-id/css/some-property"))
		})

		Context("when the session indicates a success", func() {
			It("returns the value of the CSS property", func() {
				Expect(value).To(Equal("some value"))
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the CSS property", func() {
				session.Err = errors.New("some error")
				_, err = element.GetCSS("some-property")
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
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to click", func() {
				session.Err = errors.New("some error")
				err = element.Click()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Clear", func() {
		BeforeEach(func() {
			err = element.Clear()
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("POST"))
		})

		It("hits the /element/:id/clear endpoint", func() {
			Expect(session.Endpoint).To(Equal("element/some-id/clear"))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to clear the field", func() {
				session.Err = errors.New("some error")
				err = element.Clear()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Value", func() {
		BeforeEach(func() {
			err = element.Value("text")
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("POST"))
		})

		It("hits the /element/:id/click endpoint", func() {
			Expect(session.Endpoint).To(Equal("element/some-id/value"))
		})

		It("includes the text to enter in the request body", func() {
			Expect(session.BodyJSON).To(MatchJSON(`{"value": ["t", "e", "x", "t"]}`))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to enter the text", func() {
				session.Err = errors.New("some error")
				err = element.Value("text")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsSelected", func() {
		var value bool

		BeforeEach(func() {
			session.Result = `true`
			value, err = element.IsSelected()
		})

		It("makes a GET request", func() {
			Expect(session.Method).To(Equal("GET"))
		})

		It("hits the /element/:id/selected endpoint", func() {
			Expect(session.Endpoint).To(Equal("element/some-id/selected"))
		})

		Context("when the session indicates a success", func() {
			It("returns the selected status", func() {
				Expect(value).To(BeTrue())
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the selected status", func() {
				session.Err = errors.New("some error")
				_, err = element.IsSelected()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Submit", func() {
		BeforeEach(func() {
			err = element.Submit()
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("POST"))
		})

		It("hits the /element/:id/submit endpoint", func() {
			Expect(session.Endpoint).To(Equal("element/some-id/submit"))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to submit", func() {
				session.Err = errors.New("some error")
				err = element.Submit()
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
