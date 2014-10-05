package page_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/page"


	"testing"
)

func TestPage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Page Suite")
}

type Do func(Selection)
func (f Do) Call(selection Selection) {
	f(selection)
}
