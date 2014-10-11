package phantom_test

import (
	. "github.com/sclevine/agouti/page/internal/phantom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"
)

var _ = Describe("Phantom service", func() {
	var service *Service

	BeforeEach(func() {
		service = &Service{Address: "127.0.0.1:42344", Timeout: 3 * time.Second}
	})

	Describe("#Start", func() {
		var err error

		Context("when PhantomJS is started multiple times", func() {
			It("returns an error indicating that PhantomJS is already running", func() {
				defer service.Stop()
				service.Start()
				err = service.Start()
				Expect(err).To(MatchError("PhantomJS is already running"))
			})
		})

		Context("when the PhantomJS binary is available in PATH", func() {
			BeforeEach(func() {
				err = service.Start()
			})

			AfterEach(func() {
				service.Stop()
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("starts a PhantomJS webdriver server on the provided port", func() {
				response, _ := http.Get("http://127.0.0.1:42344/status")
				body, _ := ioutil.ReadAll(response.Body)
				Expect(string(body)).To(ContainSubstring(`"status":0`))
			})
		})

		Context("when the PhantomJS binary is not available in PATH", func() {
			It("returns an error indicating the PhantomJS needs to be installed", func() {
				oldPATH := os.Getenv("PATH")
				os.Setenv("PATH", "")
				err := service.Start()
				Expect(err).To(MatchError("PhantomJS binary not found"))
				os.Setenv("PATH", oldPATH)
			})
		})

		Context("when the PhantomJS server fails to start after the provided timeout", func() {
			It("returns an error indicating that it failed to start", func() {
				service.Timeout = 0
				Expect(service.Start()).To(MatchError("PhantomJS webdriver failed to start"))
			})
		})
	})

	Describe("#Stop", func() {
		It("stops a running server", func() {
			service.Start()
			service.Stop()
			_, err := http.Get("http://127.0.0.1:42344/status")
			Expect(err).To(HaveOccurred())
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
					Expect(err).NotTo(HaveOccurred())
					Expect(session.URL).To(MatchRegexp(`http://127\.0\.0\.1:42344/session/([0-9a-f]+-)+[0-9a-f]+`))
				})
			})

			Context("if the request fails", func() {
				It("returns the request error", func() {
					service.Address = "potato"
					_, err := service.CreateSession()
					Expect(err.Error()).To(ContainSubstring("Post http://potato/session: dial tcp"))
				})
			})

			Context("if the request does not contain a session ID", func() {
				It("returns an error indicating that it failed to receive a session ID", func() {
					fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
						response.Write([]byte("{}"))
					}))
					service.Address = strings.Split(fakeServer.URL, "/")[2]
					_, err := service.CreateSession()
					Expect(err).To(MatchError("PhantomJS webdriver failed to return a session ID"))
					fakeServer.Close()
				})
			})
		})

		Context("without a running server", func() {
			It("returns an error", func() {
				service.Start()
				service.Stop()
				_, err := service.CreateSession()
				Expect(err).To(MatchError("PhantomJS not running"))
			})
		})
	})
})
