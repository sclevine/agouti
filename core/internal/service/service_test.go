package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/service"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

func freeAddress() string {
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	defer listener.Close()
	return listener.Addr().String()
}

var _ = Describe("Service", func() {
	var (
		service *Service
		url     string
	)

	BeforeEach(func() {
		address := freeAddress()
		url = "http://" + address
		service = &Service{
			URL:     url,
			Timeout: 5 * time.Second,
			Command: []string{"phantomjs", "--webdriver=" + address},
		}
	})

	Describe("#Start", func() {
		Context("when the service is started multiple times", func() {
			It("returns an error indicating that service is already running", func() {
				defer service.Stop()
				Expect(service.Start()).To(Succeed())
				err := service.Start()
				Expect(err).To(MatchError("phantomjs is already running"))
			})
		})

		Context("when the binary is available in PATH", func() {
			var err error

			BeforeEach(func() {
				err = service.Start()
			})

			AfterEach(func() {
				service.Stop()
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("starts a webdriver server on the provided port", func() {
				response, _ := http.Get(url + "/status")
				body, _ := ioutil.ReadAll(response.Body)
				Expect(string(body)).To(ContainSubstring(`"status":0`))
			})
		})

		Context("when the binary is not available in PATH", func() {
			It("returns an error indicating the binary needs to be installed", func() {
				oldPATH := os.Getenv("PATH")
				os.Setenv("PATH", "")
				err := service.Start()
				Expect(err).To(MatchError("unable to run phantomjs: exec: \"phantomjs\": executable file not found in $PATH"))
				os.Setenv("PATH", oldPATH)
			})
		})

		Context("when the service fails to start after the provided timeout", func() {
			It("returns an error indicating that it failed to start", func() {
				service.Timeout = 0
				Expect(service.Start()).To(MatchError("phantomjs webdriver failed to start"))
			})
		})
	})

	Describe("#Stop", func() {
		It("stops a running server", func() {
			service.Start()
			service.Stop()
			_, err := http.Get(url + "/status")
			Expect(err).To(HaveOccurred())
		})

		It("does nothing if the service has not been started", func() {
			service.Stop()
		})
	})

	Describe("#CreateSession", func() {
		var capabilities *Capabilities

		BeforeEach(func() {
			capabilities = &Capabilities{}
		})

		Context("with a running server", func() {
			BeforeEach(func() {
				service.Start()
			})

			AfterEach(func() {
				service.Stop()
			})

			It("makes a POST request using the desired browser name", func() {
				var requestBody string

				fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
					requestBodyBytes, _ := ioutil.ReadAll(request.Body)
					requestBody = string(requestBodyBytes)
				}))
				defer fakeServer.Close()
				service.URL = fakeServer.URL
				capabilities.BrowserName = "some-browser"
				service.CreateSession(capabilities)
				Expect(requestBody).To(Equal(`{"desiredCapabilities": {"browserName":"some-browser"}}`))
			})

			Context("if the request is invalid", func() {
				It("returns the invalid request error", func() {
					service.URL = "%@#$%"
					_, err := service.CreateSession(capabilities)
					Expect(err.Error()).To(ContainSubstring(`parse %@: invalid URL escape "%@"`))
				})
			})

			Context("if the request fails", func() {
				It("returns the failed request error", func() {
					service.URL = "http://#"
					_, err := service.CreateSession(capabilities)
					Expect(err.Error()).To(ContainSubstring("Post http://#/session"))
				})
			})

			Context("if the request does not contain a session ID", func() {
				It("returns an error indicating that it failed to receive a session ID", func() {
					fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
						response.Write([]byte("{}"))
					}))
					defer fakeServer.Close()
					service.URL = fakeServer.URL
					_, err := service.CreateSession(capabilities)
					Expect(err).To(MatchError("phantomjs webdriver failed to return a session ID"))
				})
			})

			Context("if the request succeeds", func() {
				It("returns a session with session URL", func() {
					session, err := service.CreateSession(capabilities)
					Expect(err).NotTo(HaveOccurred())
					Expect(session.URL).To(MatchRegexp(`http://127\.0\.0\.1:[0-9]+/session/([0-9a-f]+-)+[0-9a-f]+`))
				})
			})
		})

		Context("without a running server", func() {
			It("returns an error", func() {
				service.Start()
				service.Stop()
				_, err := service.CreateSession(capabilities)
				Expect(err).To(MatchError("phantomjs not running"))
			})
		})
	})
})
