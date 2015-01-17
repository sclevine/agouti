package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/mocks"
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

	ItShouldMakeAnElementRequest := func(method, endpoint string, body ...string) {
		It("should make a "+method+" request", func() {
			Expect(session.ExecuteCall.Method).To(Equal(method))
		})

		It("should hit the desired endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("element/some-id/" + endpoint))
		})

		It("should not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		if len(body) > 0 {
			It("should set the request body", func() {
				Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(body[0]))
			})
		}
	}

	Describe("#GetElement", func() {
		var singleElement *Element

		BeforeEach(func() {
			session.ExecuteCall.Result = `{"ELEMENT": "some-id"}`
			singleElement, err = element.GetElement(Selector{"css selector", "#selector"})
		})

		ItShouldMakeAnElementRequest("POST", "element", `{"using": "css selector", "value": "#selector"}`)

		It("should return an element with an ID and session", func() {
			Expect(singleElement.ID).To(Equal("some-id"))
			Expect(singleElement.Session).To(Equal(session))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.GetElement(Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		var elements []*Element

		BeforeEach(func() {
			session.ExecuteCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err = element.GetElements(Selector{"css selector", "#selector"})
		})

		ItShouldMakeAnElementRequest("POST", "elements", `{"using": "css selector", "value": "#selector"}`)

		It("should return a slice of elements with IDs and sessions", func() {
			Expect(elements[0].ID).To(Equal("some-id"))
			Expect(elements[0].Session).To(Equal(session))
			Expect(elements[1].ID).To(Equal("some-other-id"))
			Expect(elements[1].Session).To(Equal(session))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.GetElements(Selector{"css selector", "#selector"})
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

		ItShouldMakeAnElementRequest("GET", "text")

		It("should return the visible text on the element", func() {
			Expect(text).To(Equal("some text"))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the text", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.GetText()
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

		ItShouldMakeAnElementRequest("GET", "attribute/some-name")

		It("should return the value of the attribute", func() {
			Expect(value).To(Equal("some value"))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the attribute", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.GetAttribute("some-name")
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

		ItShouldMakeAnElementRequest("GET", "css/some-property")

		It("should return the value of the CSS property", func() {
			Expect(value).To(Equal("some value"))
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the CSS property", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.GetCSS("some-property")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Click", func() {
		BeforeEach(func() {
			err = element.Click()
		})

		ItShouldMakeAnElementRequest("POST", "click")

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to click", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(element.Click()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Clear", func() {
		BeforeEach(func() {
			err = element.Clear()
		})

		ItShouldMakeAnElementRequest("POST", "clear")

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to clear the field", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(element.Clear()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Value", func() {
		BeforeEach(func() {
			err = element.Value("text")
		})

		ItShouldMakeAnElementRequest("POST", "value", `{"value": ["t", "e", "x", "t"]}`)

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to enter the text", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(element.Value("text")).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsSelected", func() {
		var value bool

		BeforeEach(func() {
			session.ExecuteCall.Result = "true"
			value, err = element.IsSelected()
		})

		ItShouldMakeAnElementRequest("GET", "selected")

		It("should return the selected status", func() {
			Expect(value).To(BeTrue())
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the selected status", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.IsSelected()
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

		ItShouldMakeAnElementRequest("GET", "displayed")

		It("should return the displayed status", func() {
			Expect(value).To(BeTrue())
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the displayed status", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.IsDisplayed()
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

		ItShouldMakeAnElementRequest("GET", "enabled")

		It("should return the enabled status", func() {
			Expect(value).To(BeTrue())
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the enabled status", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.IsEnabled()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Submit", func() {
		BeforeEach(func() {
			err = element.Submit()
		})

		ItShouldMakeAnElementRequest("POST", "submit")

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to submit", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(element.Submit()).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsEqualTo", func() {
		var (
			equal        bool
			otherElement *Element
		)

		BeforeEach(func() {
			otherElement = &Element{"other-id", session}
			session.ExecuteCall.Result = "true"
			equal, err = element.IsEqualTo(otherElement)
		})

		ItShouldMakeAnElementRequest("GET", "equals/other-id")

		It("should return whether the elements are equal", func() {
			Expect(equal).To(BeTrue())
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to compare the elements", func() {
				session.ExecuteCall.Err = errors.New("some error")
				_, err := element.IsEqualTo(otherElement)
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
