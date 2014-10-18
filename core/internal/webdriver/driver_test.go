package webdriver_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	"github.com/sclevine/agouti/core/internal/session"
	. "github.com/sclevine/agouti/core/internal/webdriver"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Driver", func() {
	var (
		driver  *Driver
		service *mocks.Service
	)

	BeforeEach(func() {
		service = &mocks.Service{}
		driver = &Driver{Service: service}
	})

	Describe("#Start", func() {
		It("should start the service", func() {
			driver.Start()
			Expect(service.StartCall.Called).To(BeTrue())
		})

		Context("when starting the service fails", func() {
			It("should return an error", func() {
				service.StartCall.Err = errors.New("some error")
				Expect(driver.Start()).To(MatchError("failed to start service: some error"))
			})
		})

		Context("when starting the service succeeds", func() {
			It("should return nil", func() {
				err := driver.Start()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Stop", func() {
		var (
			fakeServer      *httptest.Server
			deletedSessions int
		)

		BeforeEach(func() {
			fakeServer = httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
				if request.Method == "DELETE" && request.URL.Path == "/" {
					deletedSessions += 1
				}
			}))
			service.CreateSessionCall.ReturnSession = &session.Session{URL: fakeServer.URL}
			driver.Page()
			driver.Page()
		})

		AfterEach(func() {
			fakeServer.Close()
		})

		It("should attempt to destroy all sessions", func() {
			driver.Stop()
			Expect(deletedSessions).To(Equal(2))
		})

		It("should stop the service", func() {
			driver.Stop()
			Expect(service.StopCall.Called).To(BeTrue())
		})
	})

	Describe("#Page", func() {
		Context("with zero arguments", func() {
			It("should create a session with no browser name", func() {
				_, err := driver.Page()
				Expect(err).NotTo(HaveOccurred())
				Expect(service.CreateSessionCall.Capabilities["browserName"]).To(BeNil())
			})
		})

		Context("with one argument", func() {
			It("should create a session with the provided browser name", func() {
				_, err := driver.Page("some-name")
				Expect(err).NotTo(HaveOccurred())
				Expect(service.CreateSessionCall.Capabilities["browserName"]).To(Equal("some-name"))
			})
		})

		Context("with more than one argument", func() {
			It("should return an error", func() {
				_, err := driver.Page("one", "two")
				Expect(err).To(MatchError("too many arguments"))
			})
		})

		Context("when creating the session fails", func() {
			It("should return an error", func() {
				service.CreateSessionCall.Err = errors.New("some error")
				_, err := driver.Page()
				Expect(err).To(MatchError("failed to generate page: some error"))
			})
		})

		It("should return a page with a client with the created session", func() {
			var sessionInPage bool
			fakeServer := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				sessionInPage = true
			}))
			defer fakeServer.Close()
			service.CreateSessionCall.ReturnSession = &session.Session{URL: fakeServer.URL}
			page, _ := driver.Page()
			page.URL()
			Expect(sessionInPage).To(BeTrue())
		})
	})
})
