package mobile_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/internal/mocks"
	"github.com/sclevine/agouti/api/mobile"
)

// TODO: meta-warp those tests into something similar to api/element_test.api .. which
// are now a lot more readable :)

var _ = Describe("Element", func() {
	var (
		element *mobile.Element
		session *mobile.Session
		bus     *mocks.Bus
		err     error
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		session = &mobile.Session{&api.Session{bus}}
		element = &mobile.Element{&api.Element{"some-id", session.Session}, session}
	})

	ItShouldMakeAnElementRequest := func(endpoint, method string, body ...string) {
		It("should hit the desired endpoint", func() {
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/" + endpoint))
		})

		It("should make a "+method+" request", func() {
			Expect(bus.SendCall.Method).To(Equal(method))
		})

		It("should not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		if len(body) > 0 {
			It("should set the request body", func() {
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(body[0]))
			})
		}
	}

	Describe("#GetElement", func() {
		var singleElement *mobile.Element

		BeforeEach(func() {
			bus.SendCall.Result = `{"ELEMENT": "some-id"}`
			singleElement, err = element.GetElement(api.Selector{"css selector", "#selector"})
		})

		ItShouldMakeAnElementRequest("element", "POST", `{"using": "css selector", "value": "#selector"}`)

		It("should return an element with an ID and session", func() {
			Expect(singleElement.ID).To(Equal("some-id"))
			Expect(singleElement.Session).To(Equal(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetElement(api.Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

})
