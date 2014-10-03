package webdriver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"testing"
)

func TestWebdriver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webdriver Suite")
}

type mockSession struct {
	endpoint string
	method   string
	bodyJSON []byte
	result   string
	err      error
}

func (m *mockSession) Execute(endpoint, method string, body, result interface{}) error {
	m.endpoint = endpoint
	m.method = method
	m.bodyJSON, _ = json.Marshal(body)
	json.Unmarshal([]byte(m.result), result)
	return m.err
}
