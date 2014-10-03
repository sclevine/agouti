package phantom_test

import (
	. "github.com/sclevine/agouti/phantom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"time"
)

var _ = Describe("Phantom service", func() {
	var service *Service

	BeforeEach(func() {
		service = &Service{Host: "127.0.0.1", Port: 42344, Timeout: time.Second * 5}
	})

	Describe("#Start", func() {
		var err error

		Context("when the phantomjs binary is available in PATH", func() {
			BeforeEach(func() {
				err = service.Start()
			})

			AfterEach(func() {
				service.Stop()
			})

			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("starts a phantom webdriver server on the provided port", func() {
				response, _ := http.Get("http://127.0.0.1:42344/status")
				body, _ := ioutil.ReadAll(response.Body)
				Expect(string(body)).To(ContainSubstring(`"status":0`))
			})
		})

		Context("when the phantomjs binary is not available in PATH", func() {
			It("returns an error indicating the phantomjs needs to be installed", func() {
				oldPATH := os.Getenv("PATH")
				os.Setenv("PATH", "")
				err := service.Start()
				Expect(err).To(MatchError("phantomjs not found"))
				os.Setenv("PATH", oldPATH)
			})
		})

		Context("when the phantomjs server fails to start after the provided timeout", func() {
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
			_, err := http.Get("http://127.0.0.1:42344/status")
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("#CreateSession", func() {
		Context("with a running server", func() {
			BeforeEach(func() {
				service.Start()
			})

			AfterEach(func() {
				service.Stop()
			})

			Context("if the request succeeds", func() {
				It("returns a session with session URL", func() {
					session, err := service.CreateSession()
					Expect(err).To(BeNil())
					Expect(session.URL).To(MatchRegexp(`http://127\.0\.0\.1:42344/session/([0-9a-f]+-)+[0-9a-f]+`))
				})
			})

			Context("if the request fails", func() {
				It("returns the request error", func() {
					service.Port = 0
					_, err := service.CreateSession()
					Expect(err.Error()).To(ContainSubstring("Post http://127.0.0.1:0/session: dial tcp 127.0.0.1:0:"))
				})
			})

			Context("if the request does not contain a session ID", func() {
				It("returns an error indicating that it failed to receive a session ID", func() {
					fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
						response.Write([]byte("{}"))
					}))
					service.Port, _ = strconv.Atoi(strings.Split(fakeServer.URL, ":")[2])
					_, err := service.CreateSession()
					Expect(err).To(MatchError("phantomjs webdriver failed to return a session ID"))
					fakeServer.Close()
				})
			})
		})

		Context("without a running server", func() {
			It("returns an error", func() {
				service.Start()
				service.Stop()
				_, err := service.CreateSession()
				Expect(err).To(MatchError("phantomjs not running"))
			})
		})
	})
})
