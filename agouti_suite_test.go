package agouti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"

	"testing"
)

func TestAgouti(t *testing.T) {
	RegisterFailHandler(Fail)
	defer CleanupAgouti(SetupAgouti())
	RunSpecs(t, "Agouti Suite")
}
