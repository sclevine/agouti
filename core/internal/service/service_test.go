package service_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/service"
)

var _ = Describe("Service", func() {
	var (
		service *Service
		started bool
	)

	BeforeEach(func() {
		started = false

		fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			if started && request.URL.Path == "/status" {
				response.WriteHeader(200)
			} else {
				response.WriteHeader(400)
			}
		}))

		service = &Service{
			URL:     fakeServer.URL,
			Timeout: 1500 * time.Millisecond,
			Command: []string{"cat"},
		}
	})

	Describe("#Start", func() {
		Context("when the service is started multiple times", func() {
			It("should return an error indicating that service is already running", func() {
				defer service.Stop()
				started = true
				Expect(service.Start()).To(Succeed())
				err := service.Start()
				Expect(err).To(MatchError("cat is already running"))
			})
		})

		Context("when the binary is not available in PATH", func() {
			It("should return an error indicating the binary needs to be installed", func() {
				service.Command = []string{"not-in-path"}
				err := service.Start()
				Expect(err).To(MatchError("unable to run not-in-path: exec: \"not-in-path\": executable file not found in $PATH"))
			})
		})

		Context("when the service starts before the provided timeout", func() {
			It("should not return an error", func() {
				defer service.Stop()
				go func() {
					time.Sleep(200 * time.Millisecond)
					started = true
				}()
				err := service.Start()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the service does not start before the provided timeout", func() {
			It("should return an error", func() {
				defer service.Stop()
				go func() {
					time.Sleep(3000 * time.Millisecond)
					started = true
				}()
				err := service.Start()
				Expect(err).To(MatchError("cat failed to start"))
			})
		})
	})

	Describe("#Stop", func() {
		It("should stop a running server", func() {
			defer service.Stop()
			started = true
			service.Start()
			service.Stop()
			err := service.Start()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("#CreateSession", func() {
		var capabilities map[string]interface{}

		BeforeEach(func() {
			capabilities = map[string]interface{}{"browserName": "some-browser"}
		})

		Context("when the server is not running", func() {
			It("should return an error", func() {
				_, err := service.CreateSession(capabilities)
				Expect(err).To(MatchError("cat not running"))
			})
		})

		Context("when the server is running", func() {
			It("should attempt to open a session using the desired capabilties", func() {
				defer service.Stop()
				started = true
				service.Start()
				var requestBody string
				fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
					requestBodyBytes, _ := ioutil.ReadAll(request.Body)
					requestBody = string(requestBodyBytes)
					response.Write([]byte(`{"sessionId": "some-id"}`))
				}))
				defer fakeServer.Close()
				service.URL = fakeServer.URL
				newSession, err := service.CreateSession(capabilities)
				Expect(err).NotTo(HaveOccurred())
				Expect(requestBody).To(MatchJSON(`{"desiredCapabilities": {"browserName": "some-browser"}}`))
				Expect(newSession.URL).To(ContainSubstring("/session/some-id"))
			})

			Context("when opening a new session fails", func() {
				It("should return the session error", func() {
					defer service.Stop()
					started = true
					service.Start()
					service.URL = "%@#$%"
					_, err := service.CreateSession(capabilities)
					Expect(err.Error()).To(ContainSubstring(`parse %@: invalid URL escape "%@"`))
				})
			})
		})
	})
})
