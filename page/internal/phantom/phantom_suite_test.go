package phantom_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPhantom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Phantom Suite")
}
