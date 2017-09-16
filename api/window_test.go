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
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		window = &Window{"some-id", &Session{bus}}
	})

	Describe("#Send", func() {
		It("should successfully send a request to the provided endpoint", func() {
			Expect(window.Send("method", "endpoint", "body", nil)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("method"))
			Expect(bus.SendCall.Endpoint).To(Equal("window/some-id/endpoint"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`"body"`))
		})

		It("should retrieve the result", func() {
			var result string
			bus.SendCall.Result = `"some result"`
			Expect(window.Send("method", "endpoint", "body", &result)).To(Succeed())
			Expect(result).To(Equal("some result"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				err := window.Send("method", "endpoint", "body", nil)
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetSize", func() {
		It("should successfully send a POST request to the size endpoint", func() {
			Expect(window.SetSize(640, 480)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("window/some-id/size"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"width":640,"height":480}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(window.SetSize(640, 480)).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetSize", func() {
		It("should successfully send a GET request to the size endpoint", func() {
			width, height, err := window.GetSize()
			Expect(err).To(Succeed())
			Expect(width).To(BeNumerically(">=", 0))
			Expect(height).To(BeNumerically(">=", 0))
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("window/some-id/size"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, _, err := window.GetSize()
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
