package browser_test

import (
	. "github.com/sclevine/agouti/page/internal/browser"
	"github.com/sclevine/agouti/page/internal/session"
	"github.com/sclevine/agouti/page/internal/driver"
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"errors"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Browser", func() {
	var (
		browser *Browser
		service *mocks.Service
	)

	BeforeEach(func() {
		service = &mocks.Service{}
		browser = &Browser{Service: service}
	})

	Describe("#Start", func() {
		It("starts the service", func() {
			browser.Start()
			Expect(service.StartCall.Called).To(BeTrue())
		})

		Context("when starting the service fails", func() {
			It("returns an error", func() {
				service.StartCall.Err = errors.New("some error")
				Expect(browser.Start()).To(MatchError("failed to start service: some error"))
			})
		})

		Context("when starting the service succeeds", func() {
			It("returns nil", func() {
				err := browser.Start()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Stop", func() {
		var deletedSessions int
		var destroyStatus int

		BeforeEach(func() {
			destroyStatus = http.StatusOK
			fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				deletedSessions += 1
				response.WriteHeader(destroyStatus)
			}))
			service.CreateSessionCall.ReturnSession = &session.Session{fakeServer.URL}
			browser.Page()
			browser.Page()
		})

		It("attempts to destroy all sessions", func() {
			browser.Stop()
			Expect(deletedSessions).To(Equal(2))
		})

		Context("when all session destroys succeed", func() {
			It("stops the service", func() {
				browser.Stop()
				Expect(service.StopCall.Called).To(BeTrue())
			})
		})

		Context("when any session destroys fail", func() {
			BeforeEach(func() {
				destroyStatus = http.StatusBadRequest
			})

			It("returns a non-fatal error", func() {
				err := browser.Stop()
				Expect(err).To(MatchError("failed to destroy all running sessions"))
			})

			It("stops the service regardless", func() {
				browser.Stop()
				Expect(service.StopCall.Called).To(BeTrue())
			})
		})
	})

	Describe("#Page", func() {
		Context("with zero arguments", func() {
			It("creates a session with no browser name", func() {
				_, err := browser.Page()
				Expect(err).NotTo(HaveOccurred())
				Expect(service.CreateSessionCall.Capabilities.BrowserName).To(Equal(""))
			})
		})

		Context("with one argument", func() {
			It("creates a session with the provided browser name", func() {
				_, err := browser.Page("some-name")
				Expect(err).NotTo(HaveOccurred())
				Expect(service.CreateSessionCall.Capabilities.BrowserName).To(Equal("some-name"))
			})
		})

		Context("with more than one argument", func() {
			It("returns an error", func() {
				_, err := browser.Page("one", "two")
				Expect(err).To(MatchError("too many arguments"))
			})
		})

		It("returns a page with a driver with the created session", func() {
			createdSession := &session.Session{}
			service.CreateSessionCall.ReturnSession = createdSession
			page, _ := browser.Page()
			sameSession := page.Driver.(*driver.Driver).Session.(*session.Session)
			Expect(createdSession).To(Equal(sameSession))
		})
	})
})
