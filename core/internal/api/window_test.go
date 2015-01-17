package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/mocks"
)

var _ = Describe("Window", func() {
	var (
		window  *Window
		session *mocks.Session
		err     error
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		window = &Window{"some-id", session}
	})

	ItShouldMakeAWindowRequest := func(method, endpoint string, body ...string) {
		It("should make a "+method+" request", func() {
			Expect(session.ExecuteCall.Method).To(Equal(method))
		})

		It("should hit the desired endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("window/some-id/" + endpoint))
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

	Describe("#SetSize", func() {
		BeforeEach(func() {
			err = window.SetSize(640, 480)
		})

		ItShouldMakeAWindowRequest("POST", "size", `{"width":640,"height":480}`)

		Context("when the session indicates a failure", func() {
			It("should return an error", func() {
				session.ExecuteCall.Err = errors.New("some error")
				Expect(window.SetSize(640, 480)).To(MatchError("some error"))
			})
		})
	})
})
