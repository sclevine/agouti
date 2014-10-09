package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/page"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	defer StopPhantom(StartPhantom())
	RunSpecs(t, "Integration Suite")
}
