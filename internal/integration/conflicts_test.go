package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Feature("Conflicts", func() {
	It("should allow dot-importing core, matchers, dsl, ginkgo, and gomega", func() {
		Expect(HaveTitle).To(BeAssignableToTypeOf(HaveTitle))
		Expect(Use()).To(BeAssignableToTypeOf(Use()))
	})
})
