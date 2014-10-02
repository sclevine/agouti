package webdriver_test

import (
	. "github.com/sclevine/agouti/webdriver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webdriver", func() {
	var api *Webdriver

	BeforeEach(func() {
		api = &Webdriver{}
	})

	It("exists", func() {
		Expect(api).To(Equal(api))
	})
})
