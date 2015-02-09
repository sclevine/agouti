package dsl_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDSL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DSL Suite")
}
