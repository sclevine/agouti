package integration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
)

var phantomDriver WebDriver

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	phantomDriver, _ = PhantomJS()
	Expect(phantomDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(phantomDriver.Stop()).To(Succeed())
})
