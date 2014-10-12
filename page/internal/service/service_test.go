package service_test

import (
	. "github.com/sclevine/agouti/page/internal/service"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"time"
)

var _ = Describe("Service", func() {
	Context("PhantomJS", func() {
		var service *Service

		BeforeEach(func() {
			desiredCapabilities := `{"desiredCapabilities": {}}`
			command := exec.Command("phantomjs", fmt.Sprintf("--webdriver=%s", "127.0.0.1:42344"))
			service = &Service{Address: "127.0.0.1:42344",
				Timeout:             3 * time.Second,
				ServiceName:         "phantomjs",
				Command:             command,
				DesiredCapabilities: desiredCapabilities}
		})

		Describe("#Start", func() {
			var err error

			Context("when PhantomJS is started multiple times", func() {
				It("returns an error indicating that PhantomJS is already running", func() {
					defer service.Stop()
					service.Start()
					err = service.Start()
					Expect(err).To(MatchError("phantomjs is already running"))
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
					Expect(err).To(MatchError("phantomjs binary not found"))
					os.Setenv("PATH", oldPATH)
				})
			})

			Context("when the PhantomJS server fails to start after the provided timeout", func() {
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
						service.Address = "#"
						_, err := service.CreateSession()
						Expect(err.Error()).To(ContainSubstring("Post http://#/session"))
					})
				})

				Context("if the request does not contain a session ID", func() {
					It("returns an error indicating that it failed to receive a session ID", func() {
						fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
							response.Write([]byte("{}"))
						}))
						service.Address = strings.Split(fakeServer.URL, "/")[2]
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

	Context("Selenium", func() {
		var service *Service

		BeforeEach(func() {
			command := exec.Command("selenium-server", "-port", "42345")
			desiredCapabilities := `{"desiredCapabilities": {"browserName": "firefox"}}`
			service = &Service{Address: "127.0.0.1:42345/wd/hub",
				Timeout:             5 * time.Second,
				ServiceName:         "selenium-server",
				Command:             command,
				DesiredCapabilities: desiredCapabilities}
		})

		Describe("#Start", func() {
			var err error

			Context("when Selenium is started multiple times", func() {
				It("returns an error indicating that Selenium is already running", func() {
					defer service.Stop()
					service.Start()
					err = service.Start()
					Expect(err).To(MatchError("selenium-server is already running"))
				})
			})

			Context("when the Selenium binary is available in PATH", func() {
				BeforeEach(func() {
					err = service.Start()
				})

				AfterEach(func() {
					service.Stop()
				})

				It("does not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("starts a Selenium webdriver server on the provided port", func() {
					response, _ := http.Get("http://127.0.0.1:42345/wd/hub/status")
					body, _ := ioutil.ReadAll(response.Body)
					Expect(string(body)).To(ContainSubstring(`"state":"success"`))
				})
			})

			Context("when the Selenium binary is not available in PATH", func() {
				It("returns an error indicating the Selenium needs to be installed", func() {
					oldPATH := os.Getenv("PATH")
					os.Setenv("PATH", "")
					err := service.Start()
					Expect(err).To(MatchError("selenium-server binary not found"))
					os.Setenv("PATH", oldPATH)
				})
			})

			Context("when the Selenium server fails to start after the provided timeout", func() {
				It("returns an error indicating that it failed to start", func() {
					service.Timeout = 0
					Expect(service.Start()).To(MatchError("selenium-server webdriver failed to start"))
				})
			})
		})

		Describe("#Stop", func() {
			It("stops a running server", func() {
				service.Start()
				service.Stop()
				_, err := http.Get("http://127.0.0.1:42345/wd/hub/status")
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
						Expect(session.URL).To(MatchRegexp(`http://127\.0\.0\.1:42345/wd/hub/session/([0-9a-f]+-)+[0-9a-f]+`))
					})
				})

				Context("if the request fails", func() {
					It("returns the request error", func() {
						service.Address = "#/wd/hub"
						_, err := service.CreateSession()
						Expect(err.Error()).To(ContainSubstring("Post http://#/wd/hub/session"))
					})
				})

				Context("if the request does not contain a session ID", func() {
					It("returns an error indicating that it failed to receive a session ID", func() {
						fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
							response.Write([]byte("{}"))
						}))
						service.Address = strings.Split(fakeServer.URL, "/")[2]
						_, err := service.CreateSession()
						Expect(err).To(MatchError("selenium-server webdriver failed to return a session ID"))
						fakeServer.Close()
					})
				})
			})

			Context("without a running server", func() {
				It("returns an error", func() {
					service.Start()
					service.Stop()
					_, err := service.CreateSession()
					Expect(err).To(MatchError("selenium-server not running"))
				})
			})
		})
	})
})
