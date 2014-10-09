package dsl_test

import (
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/dsl"
)

var _ = Feature("DSL", func() {
	Background(func() {
		Step("works without a body")
	})

	Scenario("The Agouti DSL", func() {
		Step("works with a body", func() {
			Expect(CreatePage).NotTo(BeNil())
		})
	})
})
