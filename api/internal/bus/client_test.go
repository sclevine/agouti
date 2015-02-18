package bus_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api/internal/bus"
)

var _ = Describe("Session", func() {
	var (
		requestPath        string
		requestMethod      string
		requestBody        string
		requestContentType string
		responseBody       string
		responseStatus     int
		server             *httptest.Server
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

		responseStatus = 200
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("#Send", func() {
		var client *Client

		BeforeEach(func() {
			client = &Client{SessionURL: server.URL + "/session/some-id"}
		})

		Context("with a valid request body", func() {
			It("should make a request with the provided body", func() {
				body := struct{ SomeValue string }{"some request value"}
				client.Send("some/endpoint", "POST", body)
				Expect(requestBody).To(Equal(`{"SomeValue":"some request value"}`))
			})

			It("should make a request with content type application/json", func() {
				body := struct{ SomeValue string }{"some request value"}
				client.Send("some/endpoint", "POST", body)
				Expect(requestContentType).To(Equal("application/json"))
			})
		})

		Context("with an invalid request body", func() {
			It("should return an invalid request body error", func() {
				err := client.Send("some/endpoint", "POST", func() {})
				Expect(err).To(MatchError("invalid request body: json: unsupported type: func()"))
			})
		})

		Context("when the provided body is nil", func() {
			It("should make a request without a body", func() {
				Expect(client.Send("some/endpoint", "POST", nil)).To(Succeed())
				Expect(requestBody).To(BeEmpty())
			})
		})

		It("should make a request with the full session endpoint", func() {
			client.Send("some/endpoint", "GET", nil)
			Expect(requestPath).To(Equal("/session/some-id/some/endpoint"))
		})

		It("should make a request with the given method", func() {
			client.Send("some/endpoint", "GET", nil)
			Expect(requestMethod).To(Equal("GET"))
		})

		Context("when the session endpoint is empty", func() {
			It("should make a request to the session itself", func() {
				client.Send("", "GET", nil)
				Expect(requestPath).To(Equal("/session/some-id"))
			})
		})

		Context("with an invalid URL", func() {
			It("should return an invalid request error", func() {
				client.SessionURL = "%@#$%"
				err := client.Send("some/endpoint", "GET", nil)
				Expect(err).To(MatchError(`invalid request: parse %@: invalid URL escape "%@"`))
			})
		})

		Context("when the request fails entirely", func() {
			It("should return an error indicating that the request failed", func() {
				server.Close()
				err := client.Send("some/endpoint", "GET", nil)
				Expect(err.Error()).To(MatchRegexp("request failed: .+ connection refused"))
			})
		})

		Context("when the server responds with a non-2xx status code", func() {
			Context("when the server has a valid error message", func() {
				It("should return an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{"value": {"message": "{\"errorMessage\": \"some error\"}"}}`
					err := client.Send("some/endpoint", "GET", nil)
					Expect(err).To(MatchError("request unsuccessful: some error"))
				})
			})

			Context("when the server does not have a valid message", func() {
				It("should return an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{}}`
					err := client.Send("some/endpoint", "GET", nil)
					Expect(err).To(MatchError("request unsuccessful: error unreadable"))
				})
			})

			Context("when the server does not have a valid error message", func() {
				It("should return an error from the server indicating that the request failed", func() {
					responseStatus = 400
					responseBody = `{"value": {"message": "{}}"}}`
					err := client.Send("some/endpoint", "GET", nil)
					Expect(err).To(MatchError("request unsuccessful: error message unreadable"))
				})
			})
		})

		Context("when the request succeeds", func() {
			var result struct{ Some string }

			BeforeEach(func() {
				responseBody = `{"value": {"some": "response value"}}`
			})

			Context("with a valid response body", func() {
				It("should successfully unmashal the returned JSON into the result, if provided", func() {
					Expect(client.Send("some/endpoint", "GET", nil, &result)).To(Succeed())
					Expect(result.Some).To(Equal("response value"))
				})
			})

			Context("with a response body value that cannot be read", func() {
				It("should return a failed to extract value from response error", func() {
					responseBody = `{"value": "unexpected string"}`
					err := client.Send("some/endpoint", "GET", nil, &result)
					Expect(err).To(MatchError("failed to parse response value: json: cannot unmarshal string into Go value of type struct { Some string }"))
				})
			})
		})
	})

	Describe(".Connect", func() {
		It("should make a POST request with the provided desired capabilities", func() {
			Connect(server.URL, map[string]interface{}{"some": "json"})
			Expect(requestBody).To(MatchJSON(`{"desiredCapabilities": {"some": "json"}}`))
		})

		It("should make the request with content type application/json", func() {
			Connect(server.URL, map[string]interface{}{"some": "json"})
			Expect(requestContentType).To(Equal("application/json"))
		})

		Context("when the capabilities are invalid", func() {
			It("should return an error", func() {
				_, err := Connect(server.URL, map[string]interface{}{"some": func() {}})
				Expect(err).To(MatchError("json: unsupported type: func()"))
			})
		})

		Context("when the capabilities are nil", func() {
			It("should make a POST request with empty capabilities", func() {
				Connect(server.URL, nil)
				Expect(requestBody).To(MatchJSON(`{"desiredCapabilities": {}}`))
			})
		})

		Context("when the request is invalid", func() {
			It("should return the invalid request error", func() {
				_, err := Connect("%@#$%", map[string]interface{}{"some": "json"})
				Expect(err.Error()).To(ContainSubstring(`parse %@: invalid URL escape "%@"`))
			})
		})

		Context("when the request fails", func() {
			It("should return the failed request error", func() {
				_, err := Connect("http://#", map[string]interface{}{"some": "json"})
				Expect(err.Error()).To(ContainSubstring("Post http://#/session"))
			})
		})

		Context("if the request does not contain a session ID", func() {
			It("should return an error indicating that it failed to receive a session ID", func() {
				responseBody = "{}"
				_, err := Connect(server.URL, map[string]interface{}{"some": "json"})
				Expect(err).To(MatchError("failed to retrieve a session ID"))
			})
		})

		Context("if the request succeeds", func() {
			It("should return a session with session URL", func() {
				responseBody = `{"sessionId": "some-id"}`
				client, err := Connect(server.URL, map[string]interface{}{"some": "json"})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.SessionURL).To(ContainSubstring("/session/some-id"))
			})
		})
	})
})
