package colorparser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestColorparser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Colorparser Suite")
}
