package phantom_test

import (
	. "github.com/sclevine/agouti/phantom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Session", func() {
	Describe("#Execute", func() {
		var (
			requestPath   string
			requestMethod string
			requestBody   string
			responseBody  string
			responseStatus int
			session       *Session
			result        struct{ Value string }
			server        *httptest.Server
			err           error
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				requestPath = request.URL.Path
				requestMethod = request.Method
				requestBodyBytes, _ := ioutil.ReadAll(request.Body)
				requestBody = string(requestBodyBytes)
				response.WriteHeader(responseStatus)
				response.Write([]byte(responseBody))
			}))
			session = &Session{server.URL + "/session/some-id"}
			responseBody = `{"value": "some response value"}`
			responseStatus = 200
		})

		AfterEach(func() {
			server.Close()
		})

		It("makes a request with the full session endpoint", func() {
			err = session.Execute("some/endpoint", "GET", nil, &result)
			Expect(requestPath).To(Equal("/session/some-id/some/endpoint"))
		})

		It("makes a request with the given method", func() {
			err = session.Execute("some/endpoint", "GET", nil, &result)
			Expect(requestMethod).To(Equal("GET"))
		})

		Context("with an invalid URL", func() {
			It("returns an invalid request error", func() {
				session.URL = "%@#$%"
				err = session.Execute("some/endpoint", "GET", nil, &result)
				Expect(err).To(MatchError(`invalid request: parse %@: invalid URL escape "%@"`))
			})
		})

		Context("for a GET request", func() {
			It("makes a request without a body", func() {
				err = session.Execute("some/endpoint", "GET", nil, &result)
				Expect(requestBody).To(Equal(""))
			})
		})

		Context("for a POST request", func() {
			Context("with a request valid body", func() {
				It("makes a request with the provided body", func() {
					body := struct{ SomeValue string }{"some request value"}
					err = session.Execute("some/endpoint", "POST", body, &result)
					Expect(requestBody).To(Equal(`{"SomeValue":"some request value"}`))
				})
			})

			Context("with an invalid request body", func() {
				It("returns an invalid request body error", func() {
					err = session.Execute("some/endpoint", "POST", func() {}, &result)
					Expect(err).To(MatchError("invalid request body: json: unsupported type: func()"))
				})
			})
		})

		Context("when the request succeeds", func() {
			Context("with a valid response body", func() {
				BeforeEach(func() {
					err = session.Execute("some/endpoint", "GET", nil, &result)
				})

				It("unmashals the returned JSON into the provided result", func() {
					Expect(result.Value).To(Equal("some response value"))
				})

				It("does not return an error", func() {
					Expect(err).To(BeNil())
				})
			})

			Context("with an invalid response body", func() {
				It("returns an invalid response body error", func() {
					responseBody = "}{"
					err = session.Execute("some/endpoint", "GET", nil, &result)
					Expect(err).To(MatchError("invalid response body: invalid character '}' looking for beginning of value"))
				})
			})
		})

		Context("when the request fails entirely", func() {
			It("returns an error indicating that the request failed", func() {
				server.Close()
				err = session.Execute("some/endpoint", "GET", nil, &result)
				Expect(err.Error()).To(MatchRegexp("request failed: .+ connection refused"))
			})
		})

		Context("when the server responds with a non-2xx status code", func() {
			It("returns an error indicating that the request failed", func() {
				responseStatus = 400
				err = session.Execute("some/endpoint", "GET", nil, &result)
				Expect(err).To(MatchError("request unsuccessful: 400 - Bad Request"))
			})
		})
	})
})
