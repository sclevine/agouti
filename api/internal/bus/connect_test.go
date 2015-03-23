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
		responseBody = `{"sessionId": "some-id"}`
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

	It("should successfully make an application/json POST request to the session endpoint", func() {
		_, err := Connect(server.URL, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(requestMethod).To(Equal("POST"))
		Expect(requestPath).To(Equal("/session"))
		Expect(requestContentType).To(Equal("application/json"))
	})

	It("should return a client with a session URL", func() {
		client, err := Connect(server.URL, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(client.SessionURL).To(ContainSubstring("/session/some-id"))
	})

	It("should make the request with the provided desired capabilities", func() {
		_, err := Connect(server.URL, map[string]interface{}{"some": "json"})
		Expect(err).NotTo(HaveOccurred())
		Expect(requestBody).To(MatchJSON(`{"desiredCapabilities": {"some": "json"}}`))
	})

	Context("when the capabilities are nil", func() {
		It("should make the request with empty capabilities", func() {
			_, err := Connect(server.URL, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(requestBody).To(MatchJSON(`{"desiredCapabilities": {}}`))
		})
	})

	Context("when the capabilities are invalid", func() {
		It("should return an error", func() {
			_, err := Connect(server.URL, map[string]interface{}{"some": func() {}})
			Expect(err).To(MatchError("json: unsupported type: func()"))
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

	Context("when the response contains invalid JSON", func() {
		It("should return an error", func() {
			responseBody = "$$$"
			_, err := Connect(server.URL, nil)
			Expect(err).To(MatchError("invalid character '$' looking for beginning of value"))
		})
	})

	Context("when the response does not contain a session ID", func() {
		It("should return an error", func() {
			responseBody = "{}"
			_, err := Connect(server.URL, nil)
			Expect(err).To(MatchError("failed to retrieve a session ID"))
		})
	})
})
