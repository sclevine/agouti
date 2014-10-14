package page_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Page Suite")
}
