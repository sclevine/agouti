package webdriver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWebdriver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WebDriver Suite")
}
