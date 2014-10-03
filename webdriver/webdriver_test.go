package webdriver_test

import (
	. "github.com/sclevine/agouti/webdriver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
	"errors"
)

type mockSession struct {
	endpoint string
	method   string
	bodyJSON []byte
	result   interface{}
	err      error
}

func (m *mockSession) Execute(endpoint, method string, body, result interface{}) error {
	m.endpoint = endpoint
	m.method = method
	m.bodyJSON, _ = json.Marshal(body)
	json.Unmarshal(result, m.result)
	return m.err
}

var _ = Describe("Webdriver", func() {
	var (
		driver	*Driver
		session *mockSession
		err		error
	)

	BeforeEach(func() {
		session = &mockSession{}
		driver = &Driver{session}
	})

	Describe("#Navigate", func() {
		BeforeEach(func() {
			err = driver.Navigate("http://example.com")
		})
		
		It("makes a POST request", func() {
			Expect(session.method).To(Equal("POST"));
		})

		It("hits the /url endpoint", func() {
			Expect(session.endpoint).To(Equal("url"))
		})

		It("includes the new URL in the request body", func() {
			Expect(session.bodyJSON).To(MatchJSON(`{"url": "http://example.com"}`))
		})

		Context("when the response indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the response indicates a failure", func() {
			It("returns an error indicating the page failed to navigate", func() {
				session.err = errors.New("some error")
				err = driver.Navigate("http://example.com")
				Expect(err).To(MatchError("failed to navigate: some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		BeforeEach(func() {
			err = driver.GetElements("#selector")
		})

		It("makes a POST request", func() {
			Expect(session.method).To(Equal("POST"));
		})

		It("hits the /url endpoint", func() {
			Expect(session.endpoint).To(Equal("elements"))
		})

		It("includes the new URL in the request body", func() {
			Expect(session.bodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		Context("when the response indicates a success", func() {
			It("doesn't return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("when the response indicates a failure", func() {
			It("returns an error indicating the page failed to navigate", func() {
				session.err = errors.New("some error")
				err = driver.GetElements("#selector")
				Expect(err).To(MatchError("failed to get elements with selector '#selector': some error"))
			})
		})
	})
})
