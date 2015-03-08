package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Conflicts", func() {
	It("should allow importing agouti while dot-importing matchers, ginkgo, and gomega", func() {
		Expect(agouti.Capabilities{}).To(Equal(agouti.Capabilities{}))
		Expect(HaveTitle("title")).To(Equal(HaveTitle("title")))
	})
})
