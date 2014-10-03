package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	defer CleanupAgouti(SetupAgouti())
	RunSpecs(t, "Integration Suite")
}
