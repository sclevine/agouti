package agouti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAgouti(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Agouti Suite")
}
