package bus_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api/internal/bus"
)

var _ = Describe(".Connect", func() {
	var (
		requestPath        string
		requestMethod      string
		requestBody        string
		requestContentType string
		responseBody       string
		server             *httptest.Server
	)

	BeforeEach(func() {
		responseBody = ""
		requestPath, requestMethod, requestBody, requestContentType = "", "", "", ""
		server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			requestPath = request.URL.Path // TODO: use
			requestMethod = request.Method
			requestBodyBytes, _ := ioutil.ReadAll(request.Body)
			requestBody = string(requestBodyBytes)
			requestContentType = request.Header.Get("Content-Type")
			response.Write([]byte(responseBody))
		}))
	})

	AfterEach(func() {
		server.Close()
	})

	It("should make a POST request", func() {
		Connect(server.URL, nil)
		Expect(requestMethod).To(Equal("POST"))
	})

	It("should make the request to the session endpoint", func() {
		Connect(server.URL, nil)
		Expect(requestPath).To(Equal("/session"))
	})

	It("should make the request with the provided desired capabilities", func() {
		Connect(server.URL, map[string]interface{}{"some": "json"})
		Expect(requestBody).To(MatchJSON(`{"desiredCapabilities": {"some": "json"}}`))
	})

	It("should make the request with content type application/json", func() {
		Connect(server.URL, nil)
		Expect(requestContentType).To(Equal("application/json"))
	})

	Context("when the capabilities are invalid", func() {
		It("should return an error", func() {
			_, err := Connect(server.URL, map[string]interface{}{"some": func() {}})
			Expect(err).To(MatchError("json: unsupported type: func()"))
		})
	})

	Context("when the capabilities are nil", func() {
		It("should make the request with empty capabilities", func() {
			Connect(server.URL, nil)
			Expect(requestBody).To(MatchJSON(`{"desiredCapabilities": {}}`))
		})
	})

	Context("when the request is invalid", func() {
		It("should return an error", func() {
			_, err := Connect("%@#$%", nil)
			Expect(err.Error()).To(ContainSubstring(`parse %@: invalid URL escape "%@"`))
		})
	})

	Context("when the request fails", func() {
		It("should return an error", func() {
			_, err := Connect("http://#", nil)
			Expect(err.Error()).To(ContainSubstring("Post http://#/session"))
		})
	})

	Context("if the request succeeds", func() {
		Context("with a valid response body", func() {
			It("should return a session with session URL", func() {
				responseBody = `{"sessionId": "some-id"}`
				client, err := Connect(server.URL, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(client.SessionURL).To(ContainSubstring("/session/some-id"))
			})
		})

		Context("with an response that is invalid JSON", func() {
			It("should return an error", func() {
				responseBody = "$$$"
				_, err := Connect(server.URL, nil)
				Expect(err).To(MatchError("invalid character '$' looking for beginning of value"))
			})
		})

		Context("with a response that does not contain a session ID", func() {
			It("should return an error", func() {
				responseBody = "{}"
				_, err := Connect(server.URL, nil)
				Expect(err).To(MatchError("failed to retrieve a session ID"))
			})
		})
	})
})
