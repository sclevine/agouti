package window_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWindow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Window Suite")
}
