package chrome_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"

	"testing"
)

func TestChrome(t *testing.T) {
	RegisterFailHandler(Fail)
	//RunSpecs(t, "Chrome Suite")
}

var _ = BeforeSuite(func() {
	StartChrome()
})

var _ = AfterSuite(func() {
	StopWebdriver()
})
