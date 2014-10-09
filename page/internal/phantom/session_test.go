package phantom_test

import (
	. "github.com/sclevine/agouti/page/internal/phantom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Session", func() {
	Describe("#Execute", func() {
		var (
			requestPath        string
			requestMethod      string
			requestBody        string
			requestContentType string
			responseBody       string
			responseStatus     int
			session            *Session
			result             struct{ Some string }
			server             *httptest.Server
			err                error
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				requestPath = request.URL.Path
				requestMethod = request.Method
				requestBodyBytes, _ := ioutil.ReadAll(request.Body)
				requestBody = string(requestBodyBytes)
				requestContentType = request.Header.Get("Content-Type")
				response.WriteHeader(responseStatus)
				response.Write([]byte(responseBody))
			}))
			session = &Session{server.URL + "/session/some-id"}
			responseBody = `{"value": {"some": "response value"}}`
			responseStatus = 200
		})

		AfterEach(func() {
			server.Close()
		})

		It("makes a request with the full session endpoint", func() {
			session.Execute("some/endpoint", "GET", nil, &result)
			Expect(requestPath).To(Equal("/session/some-id/some/endpoint"))
		})

		It("makes a request with the given method", func() {
			session.Execute("some/endpoint", "GET", nil, &result)
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
			BeforeEach(func() {
				err = session.Execute("some/endpoint", "GET", nil, &result)
			})

			It("makes a request without a body", func() {
				Expect(requestBody).To(BeEmpty())
			})

			It("makes a request without a content type", func() {
				Expect(requestContentType).To(BeEmpty())
			})

			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("for a POST request", func() {
			Context("with a valid request body", func() {
				BeforeEach(func() {
					body := struct{ SomeValue string }{"some request value"}
					err = session.Execute("some/endpoint", "POST", body, &result)
				})

				It("makes a request with the provided body", func() {
					Expect(requestBody).To(Equal(`{"SomeValue":"some request value"}`))
				})

				It("makes a request with content type application/json", func() {
					Expect(requestContentType).To(Equal("application/json"))
				})

				It("does not return an error", func() {
					Expect(err).To(BeNil())
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
					Expect(result.Some).To(Equal("response value"))
				})

				It("does not return an error", func() {
					Expect(err).To(BeNil())
				})
			})

			Context("with a response body value that cannot be read", func() {
				It("returns a failed to extract value from response error", func() {
					responseBody = `{"value": "unexpected string"}`
					err = session.Execute("some/endpoint", "GET", nil, &result)
					Expect(err).To(MatchError("failed to parse response value: json: cannot unmarshal string into Go value of type struct { Some string }"))
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
			Context("when the server has a valid error message", func() {
				It("returns an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{"value": {"message": "{\"errorMessage\": \"some error\"}"}}`
					err = session.Execute("some/endpoint", "GET", nil, &result)
					Expect(err).To(MatchError("request unsuccessful: some error"))
				})
			})

			Context("when the server does not have a valid message", func() {
				It("returns an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{}}`
					err = session.Execute("some/endpoint", "GET", nil, &result)
					Expect(err).To(MatchError("request unsuccessful: phantom error unreadable"))
				})
			})

			Context("when the server does not have a valid error message", func() {
				It("returns an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{"value": {"message": "{}}"}}`
					err = session.Execute("some/endpoint", "GET", nil, &result)
					Expect(err).To(MatchError("request unsuccessful: phantom error message unreadable"))
				})
			})
		})
	})
})
