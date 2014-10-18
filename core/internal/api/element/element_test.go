package element_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/api/element"
	"github.com/sclevine/agouti/core/internal/mocks"
	"github.com/sclevine/agouti/core/internal/types"
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

	Describe("#GetID", func() {
		It("should return the stored element ID", func() {
			Expect(element.GetID()).To(Equal("some-id"))
		})
	})

	Describe("#GetElements", func() {
		var elements []types.Element

		BeforeEach(func() {
			session.ExecuteCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = element.GetElements(types.Selector{Using: "css selector", Value: "#selector"})
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /element/:id/elements endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/elements"))
		})

		It("should include the selection in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		Context("when the session indicates a success", func() {
			It("should return a slice of elements with IDs and sessions", func() {
				Expect(elements[0].(*Element).ID).To(Equal("some-id"))
				Expect(elements[0].(*Element).Session).To(Equal(session))
				Expect(elements[1].(*Element).ID).To(Equal("some-other-id"))
				Expect(elements[1].(*Element).Session).To(Equal(session))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = element.GetElements(types.Selector{Using: "css selector", Value: "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetText", func() {
		var text string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"some text"`
			text, err = element.GetText()
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /element/:id/text endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/text"))
		})

		Context("when the session indicates a success", func() {
			It("should return the visible text on the element", func() {
				Expect(text).To(Equal("some text"))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the text", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = element.GetText()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetAttribute", func() {
		var value string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"some value"`
			value, err = element.GetAttribute("some-name")
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /element/:id/attribute/:name endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/attribute/some-name"))
		})

		Context("when the session indicates a success", func() {
			It("should return the value of the attribute", func() {
				Expect(value).To(Equal("some value"))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the attribute", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = element.GetAttribute("some-name")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetCSS", func() {
		var value string

		BeforeEach(func() {
			session.ExecuteCall.Result = `"some value"`
			value, err = element.GetCSS("some-property")
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /element/:id/css/:name endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/css/some-property"))
		})

		Context("when the session indicates a success", func() {
			It("should return the value of the CSS property", func() {
				Expect(value).To(Equal("some value"))
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the CSS property", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = element.GetCSS("some-property")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Click", func() {
		BeforeEach(func() {
			err = element.Click()
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /element/:id/click endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/click"))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to click", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = element.Click()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Clear", func() {
		BeforeEach(func() {
			err = element.Clear()
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /element/:id/clear endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/clear"))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to clear the field", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = element.Clear()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Value", func() {
		BeforeEach(func() {
			err = element.Value("text")
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /element/:id/click endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/value"))
		})

		It("should include the text to enter in the request body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"value": ["t", "e", "x", "t"]}`))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to enter the text", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = element.Value("text")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsSelected", func() {
		var value bool

		BeforeEach(func() {
			session.ExecuteCall.Result = "true"
			value, err = element.IsSelected()
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /element/:id/selected endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/selected"))
		})

		Context("when the session indicates a success", func() {
			It("should return the selected status", func() {
				Expect(value).To(BeTrue())
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the selected status", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = element.IsSelected()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsDisplayed", func() {
		var value bool

		BeforeEach(func() {
			session.ExecuteCall.Result = "true"
			value, err = element.IsDisplayed()
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /element/:id/displayed endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/displayed"))
		})

		Context("when the session indicates a success", func() {
			It("should return the displayed status", func() {
				Expect(value).To(BeTrue())
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the displayed status", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = element.IsDisplayed()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsEnabled", func() {
		var value bool

		BeforeEach(func() {
			session.ExecuteCall.Result = "true"
			value, err = element.IsEnabled()
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /element/:id/enabled endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/enabled"))
		})

		Context("when the session indicates a success", func() {
			It("should return the enabled status", func() {
				Expect(value).To(BeTrue())
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the enabled status", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = element.IsEnabled()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Submit", func() {
		BeforeEach(func() {
			err = element.Submit()
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /element/:id/submit endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/submit"))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to submit", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = element.Submit()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Equals", func() {
		var (
			equal        bool
			otherElement *Element
		)

		BeforeEach(func() {
			otherElement = &Element{"other-id", session}
			equal, err = element.IsEqualTo(otherElement)
		})

		It("should make a GET request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("GET"))
		})

		It("should hit the /element/:id/equals/:other endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/equals/other-id"))
		})

		Context("when the session indicates a success", func() {
			It("should return true when the comparison returns true", func() {
				session.ExecuteCall.Result = "true"
				equal, _ = element.IsEqualTo(otherElement)
				Expect(equal).To(BeTrue())
			})

			It("should return false when the comparison returns false", func() {
				session.ExecuteCall.Result = "false"
				equal, _ = element.IsEqualTo(otherElement)
				Expect(equal).To(BeFalse())
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to compare the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err = element.IsEqualTo(otherElement)
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
