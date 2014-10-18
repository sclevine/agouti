package session_test

import (
	. "github.com/sclevine/agouti/core/internal/session"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Session", func() {
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

	Describe("#Execute", func() {
		Context("with an invalid request body", func() {
			It("should return an invalid request body error", func() {
				err = session.Execute("some/endpoint", "POST", func() {})
				Expect(err).To(MatchError("invalid request body: json: unsupported type: func()"))
			})
		})

		Context("with a valid request body", func() {
			It("should make a request with the provided body", func() {
				body := struct{ SomeValue string }{"some request value"}
				session.Execute("some/endpoint", "POST", body)
				Expect(requestBody).To(Equal(`{"SomeValue":"some request value"}`))
			})
		})

		Context("when the provided body is nil", func() {
			It("should make a request without a body", func() {
				err := session.Execute("some/endpoint", "POST", nil)
				Expect(requestBody).To(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
			})
		})

		It("should make a request with the full session endpoint", func() {
			session.Execute("some/endpoint", "GET", nil)
			Expect(requestPath).To(Equal("/session/some-id/some/endpoint"))
		})

		Context("when the session endpoint is empty", func() {
			It("should make a request to the session itself", func() {
				session.Execute("", "GET", nil)
				Expect(requestPath).To(Equal("/session/some-id"))
			})
		})

		It("should make a request with the given method", func() {
			session.Execute("some/endpoint", "GET", nil)
			Expect(requestMethod).To(Equal("GET"))
		})

		Context("with an invalid URL", func() {
			It("should return an invalid request error", func() {
				session.URL = "%@#$%"
				err = session.Execute("some/endpoint", "GET", nil)
				Expect(err).To(MatchError(`invalid request: parse %@: invalid URL escape "%@"`))
			})
		})

		Context("for a POST request", func() {
			It("should make a request with content type application/json", func() {
				session.Execute("some/endpoint", "POST", nil)
				Expect(requestContentType).To(Equal("application/json"))
			})
		})

		Context("when the request fails entirely", func() {
			It("should return an error indicating that the request failed", func() {
				server.Close()
				err = session.Execute("some/endpoint", "GET", nil)
				Expect(err.Error()).To(MatchRegexp("request failed: .+ connection refused"))
			})
		})

		Context("when the server responds with a non-2xx status code", func() {
			Context("when the server has a valid error message", func() {
				It("should return an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{"value": {"message": "{\"errorMessage\": \"some error\"}"}}`
					err = session.Execute("some/endpoint", "GET", nil)
					Expect(err).To(MatchError("request unsuccessful: some error"))
				})
			})

			Context("when the server does not have a valid message", func() {
				It("should return an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{}}`
					err = session.Execute("some/endpoint", "GET", nil)
					Expect(err).To(MatchError("request unsuccessful: error unreadable"))
				})
			})

			Context("when the server does not have a valid error message", func() {
				It("should return an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{"value": {"message": "{}}"}}`
					err = session.Execute("some/endpoint", "GET", nil)
					Expect(err).To(MatchError("request unsuccessful: error message unreadable"))
				})
			})
		})

		Context("when the request succeeds", func() {
			Context("with a valid response body", func() {
				BeforeEach(func() {
					err = session.Execute("some/endpoint", "GET", nil, &result)
				})

				It("should unmashal the returned JSON into the result, if provided", func() {
					Expect(result.Some).To(Equal("response value"))
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("with a response body value that cannot be read", func() {
				It("should return a failed to extract value from response error", func() {
					responseBody = `{"value": "unexpected string"}`
					err = session.Execute("some/endpoint", "GET", nil, &result)
					Expect(err).To(MatchError("failed to parse response value: json: cannot unmarshal string into Go value of type struct { Some string }"))
				})
			})
		})
	})

	Describe(".Open", func() {
		var capabilities map[string]interface{}

		BeforeEach(func() {
			capabilities = map[string]interface{}{}
		})

		It("should make a POST request using the desired browser name", func() {
			var requestBody string

			fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				requestBodyBytes, _ := ioutil.ReadAll(request.Body)
				requestBody = string(requestBodyBytes)
			}))
			defer fakeServer.Close()
			capabilities["browserName"] = "some-browser"
			Open(fakeServer.URL, capabilities)
			Expect(requestBody).To(MatchJSON(`{"desiredCapabilities": {"browserName": "some-browser"}}`))
		})

		Context("when the request is invalid", func() {
			It("should return the invalid request error", func() {
				_, err := Open("%@#$%", capabilities)
				Expect(err.Error()).To(ContainSubstring(`parse %@: invalid URL escape "%@"`))
			})
		})

		Context("when the request fails", func() {
			It("should return the failed request error", func() {
				_, err := Open("http://#", capabilities)
				Expect(err.Error()).To(ContainSubstring("Post http://#/session"))
			})
		})

		Context("if the request does not contain a session ID", func() {
			It("should return an error indicating that it failed to receive a session ID", func() {
				fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
					response.Write([]byte("{}"))
				}))
				defer fakeServer.Close()
				_, err := Open(fakeServer.URL, capabilities)
				Expect(err).To(MatchError("failed to retrieve a session ID"))
			})
		})

		Context("if the request succeeds", func() {
			It("should return a session with session URL", func() {
				fakeServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
					response.Write([]byte(`{"sessionId": "some-id"}`))
				}))
				defer fakeServer.Close()
				session, err := Open(fakeServer.URL, capabilities)
				Expect(err).NotTo(HaveOccurred())
				Expect(session.URL).To(ContainSubstring("/session/some-id"))
			})
		})
	})
})
