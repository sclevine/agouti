package window_test

import (
	. "github.com/sclevine/agouti/webdriver/window"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/mocks"
	"errors"
)

var _ = Describe("Window", func() {
	var (
		window *Window
		session *mocks.Session
		err     error
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		window = &Window{"some-id", session}
	})

	Describe("#SetSize", func() {
		BeforeEach(func() {
			session.Result = `"text is unimportant"`
			err = window.SetSize(640,480)
		})

		It("makes a POST request", func() {
			Expect(session.Method).To(Equal("POST"))
		})

		It("hits the /window/:id/size", func() {
			Expect(session.Endpoint).To(Equal("window/some-id/size"))
		})

		It("sends the width and height as the post body", func() {
			Expect(session.BodyJSON).To(MatchJSON(`{"width":480,"height":640}`))
		})

		Context("when the session indicates a success", func() {
			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the session indicates a failure", func() {
			It("returns an error indicating the session failed to retrieve the text", func() {
				session.Err = errors.New("some error")
				err = window.SetSize(640,480)
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
