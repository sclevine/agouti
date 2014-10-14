package chrome_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"

	"os"
	"testing"
)

func TestChrome(t *testing.T) {
	RegisterFailHandler(Fail)
	if os.Getenv("HEADLESS_ONLY") != "true" {
		RunSpecs(t, "Chrome Suite")
	}
}

var _ = BeforeSuite(func() {
	StartChrome()
})

var _ = AfterSuite(func() {
	StopWebdriver()
})
