package phantom_test

import (
	. "github.com/sclevine/agouti/phantom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"os"
	"time"
	"strings"
	"strconv"
)

func makeRequest(url string) string {
	response, err := http.Get(url)
	if err != nil {
		return ""
	}
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}

var _ = Describe("Phantom server", func() {
	var server Phantom

	BeforeEach(func() {
		server = Phantom{Host: "127.0.0.1", Port: 8910, Timeout: time.Second * 5}
	})

	Describe("#Start", func() {
		var err error

		Context("when the phantomjs binary is available in PATH", func() {
			BeforeEach(func() {
				err = server.Start()
			})

			AfterEach(func() {
				server.Stop()
			})

			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("starts a phantom webdriver server on port 8910", func() {
				Expect(makeRequest("http://127.0.0.1:8910/status")).To(ContainSubstring(`"status":0`))
			})
		})

		Context("when the phantomjs binary is not available in PATH", func() {
			It("returns an error indicating the phantomjs needs to be installed", func() {
				oldPATH := os.Getenv("PATH")
				os.Setenv("PATH", "")
				err := server.Start()
				Expect(err).To(MatchError("phantomjs not found"))
				os.Setenv("PATH", oldPATH)
			})
		})

		Context("when the phantomjs server fails to start after the provided timeout", func() {
			It("returns an error indicating that it failed to start", func() {
				server.Timeout = 0
				Expect(server.Start()).To(MatchError("phantomjs webdriver failed to start"))
			})
		})
	})

	Describe("#Stop", func() {
		It("stops a running server", func() {
			server.Start()
			server.Stop()
			Expect(makeRequest("http://127.0.0.1:8910/status")).NotTo(ContainSubstring(`"status":0`))
		})
	})

	Describe("#CreateSession", func() {
		Context("with a running server", func() {
			BeforeEach(func() {
				server.Start()
			})

			AfterEach(func() {
				server.Stop()
			})

			Context("if the request succeeds", func() {
				It("returns a session URL", func() {
					url, err := server.CreateSession()
					Expect(err).To(BeNil())
					Expect(url).To(MatchRegexp(`http://127\.0\.0\.1:8910/session/([0-9a-f]+-)+[0-9a-f]+`))
				})
			})

			Context("if the request fails", func() {
				It("returns the request error", func() {
					server.Port = 0
					_, err := server.CreateSession()
					Expect(err).To(MatchError("Post http://127.0.0.1:0/session: dial tcp 127.0.0.1:0: can't assign requested address"))
				})
			})

			Context("if the request does not contain a session ID", func() {
				It("returns an error indicating that it failed to receive a session ID", func() {
					fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
						response.Write([]byte("{}"))
					}))
					server.Port, _ = strconv.Atoi(strings.Split(fakeServer.URL, ":")[2])
					_, err := server.CreateSession()
					Expect(err).To(MatchError("phantomjs webdriver failed to return a session ID"))
				})
			})
		})

		Context("without a running server", func() {
			It("returns an error", func() {
				server.Start()
				server.Stop()
				_, err := server.CreateSession()
				Expect(err).To(MatchError("phantomjs not running"))
			})
		})
	})
})
