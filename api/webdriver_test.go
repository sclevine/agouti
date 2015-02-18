package api_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/internal/mocks"
)

var _ = Describe("WebDriver", func() {
	var (
		webDriver *WebDriver
		service   *mocks.Service
	)

	BeforeEach(func() {
		service = &mocks.Service{}
		webDriver = &WebDriver{Timeout: 2 * time.Second, Service: service}
	})

	Describe("#Open", func() {
		var (
			server        *httptest.Server
			requestBody   string
			requestMethod string
			responseBody  string
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				requestBodyBytes, _ := ioutil.ReadAll(request.Body)
				requestBody = string(requestBodyBytes)
				requestMethod = request.Method
				response.Write([]byte(responseBody))
			}))
			service.URLCall.ReturnURL = server.URL
		})

		AfterEach(func() {
			server.Close()
		})

		It("should successfully return a session with a bus that talks to the WebDriver", func() {
			responseBody = `{"sessionId": "some-id"}`
			session, err := webDriver.Open(nil)
			Expect(err).NotTo(HaveOccurred())
			responseBody = `{"value": "some title"}`
			Expect(session.GetTitle()).To(Equal("some title"))
		})

		It("should open the connection with the provided desired capabilities", func() {
			webDriver.Open(map[string]interface{}{"some": "capability"})
			Expect(requestBody).To(Equal(`{"desiredCapabilities":{"some":"capability"}}`))
		})

		Context("when the service URL cannot be retrieved", func() {
			It("should return an error", func() {
				service.URLCall.Err = errors.New("some error")
				_, err := webDriver.Open(nil)
				Expect(err).To(MatchError("cannot retrieve URL: some error"))
			})
		})

		Context("when we cannot connect to the WebDriver bus", func() {
			It("should return an error", func() {
				_, err := webDriver.Open(nil)
				Expect(err).To(MatchError("failed to connect: failed to retrieve a session ID"))
			})
		})

		Context("when the WebDriver is stopped", func() {
			It("should delete the opened session", func() {
				responseBody = `{"sessionId": "some-id"}`
				webDriver.Open(nil)
				requestMethod = ""
				webDriver.Stop()
				Expect(requestBody).To(Equal(""))
				Expect(requestMethod).To(Equal("DELETE"))
			})
		})
	})

	Describe("#Start", func() {
		It("should successfully start the WebDriver service", func() {
			Expect(webDriver.Start()).To(Succeed())
			Expect(service.StartCall.Timeout).To(Equal(2 * time.Second))
		})

		Context("when the WebDriver service cannot be started", func() {
			It("should return an error", func() {
				service.StartCall.Err = errors.New("some error")
				err := webDriver.Start()
				Expect(err).To(MatchError("failed to start service: some error"))
			})
		})
	})

	Describe("#Stop", func() {
		It("should successfully stop the WebDriver service", func() {
			Expect(webDriver.Stop()).To(Succeed())
			Expect(service.StopCall.Called).To(BeTrue())
		})

		Context("when the WebDriver service cannot be stopped", func() {
			It("should return an error", func() {
				service.StopCall.Err = errors.New("some error")
				err := webDriver.Stop()
				Expect(err).To(MatchError("failed to stop service: some error"))
			})
		})
	})
})
