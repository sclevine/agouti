package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/internal/mocks"
)

var _ = Describe("Window", func() {
	var (
		window *Window
		bus    *mocks.Bus
		err    error
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		window = &Window{"some-id", &Session{bus}}
	})

	ItShouldMakeAWindowRequest := func(method, endpoint string, body ...string) {
		It("should make a "+method+" request", func() {
			Expect(bus.SendCall.Method).To(Equal(method))
		})

		It("should hit the desired endpoint", func() {
			Expect(bus.SendCall.Endpoint).To(Equal("window/some-id/" + endpoint))
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

	Describe("#SetSize", func() {
		BeforeEach(func() {
			err = window.SetSize(640, 480)
		})

		ItShouldMakeAWindowRequest("POST", "size", `{"width":640,"height":480}`)

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(window.SetSize(640, 480)).To(MatchError("some error"))
			})
		})
	})
})
