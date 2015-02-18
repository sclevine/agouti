package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Feature("Conflicts", func() {
	It("should allow dot-importing matchers, dsl, ginkgo, and gomega", func() {
		Expect(agouti.Cookie{}).To(Equal(agouti.Cookie{}))
		Expect(HaveTitle("title")).To(Equal(HaveTitle("title")))
	})
})
