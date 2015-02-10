package bus_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bus Suite")
}
