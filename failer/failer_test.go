package failer_test

import (
	. "github.com/sclevine/agouti/failer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Failer", func() {
	Describe("#Fail", func() {
		It("panics when asynchronous", func() {
			failer := &Failer{}
			failer.Async()
			defer func() {
				Expect(recover()).To(Equal("async panic"))
			}()
			failer.Fail("async panic")
		})

		// NOTE: cannot test actually failing the specs using Ginkgo
	})

	// NOTE: cannot test #Skip or #UnSkip for the same reason
})
