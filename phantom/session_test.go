package phantom_test

import (
	. "github.com/sclevine/agouti/phantom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
)

var _ = Describe("Session", func() {
	Describe("#Execute", func() {
		var (
			requestPath string
			requestMethod string
			requestBody interface{}
			session Session
			result struct{Value string}
			server *httptest.Server
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				requestPath = request.URL.Path
				requestMethod = request.Method
				requestBody, _ = ioutil.ReadAll(request.Body)
				response.Write([]byte(`{"value": "some response value"}`))
			}))
			session = Session(server.URL + "/session/some-id")
		})

		AfterEach(func() {
			server.Close()
		})

		Context("for a GET request", func() {
			BeforeEach(func() {
				session.Execute("/some/endpoint", "GET", nil, &result)
			})

			It("makes a request without a body", func() {
				Expect(requestBody).To(Equal(""))
			})

			It("makes a GET request", func() {
				Expect(requestMethod).To(Equal("GET"))
			})

			It("makes a request with the full session endpoint", func() {
				Expect(requestPath).To(Equal("/session/some-id/some/endpoint"))
			})

			It("unmashalls the returned JSON into the provided result", func() {
				Expect(result.Value).To(Equal("some response value"))
			})
		})

		Context("for a POST request", func() {
			BeforeEach(func() {
				body := struct{SomeValue string}{"some request value"}
				session.Execute("/some/endpoint", "POST", body, &result)
			})

			It("makes a request with the provided body", func() {
				Expect(requestBody).To(Equal(`{"someValue": "some request value"}`))
			})

			It("makes a POST request", func() {
				Expect(requestMethod).To(Equal("POST"))
			})

			It("makes a request with the full session endpoint", func() {
				Expect(requestPath).To(Equal("/session/some-id/some/endpoint"))
			})

			It("unmashalls the returned JSON into the provided result", func() {
				Expect(result.Value).To(Equal("some response value"))
			})
		})
	})
})
