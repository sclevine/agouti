package window_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/api/window"
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

	Describe("#SetSize", func() {
		BeforeEach(func() {
			err = window.SetSize(640, 480)
		})

		It("should make a POST request", func() {
			Expect(session.ExecuteCall.Method).To(Equal("POST"))
		})

		It("should hit the /window/:id/size endpoint", func() {
			Expect(session.ExecuteCall.Endpoint).To(Equal("window/some-id/size"))
		})

		It("should send the width and height as the post body", func() {
			Expect(session.ExecuteCall.BodyJSON).To(MatchJSON(`{"width":640,"height":480}`))
		})

		Context("when the session indicates a success", func() {
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the session indicates a failure", func() {
			It("should return an error indicating the session failed to retrieve the text", func() {
				session.ExecuteCall.Err = errors.New("some error")
				err = window.SetSize(640, 480)
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
