package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/sclevine/agouti/core"
)

var _ = Describe("Capabilities", func() {
	var capabilities Capabilities

	BeforeEach(func() {
		capabilities = Use()
	})

	It("should successfully encode all provided options into JSON", func() {
		capabilities.Browser("some-browser").Version("v100").Platform("some-os")
		capabilities.With("enabledThing").Without("disabledThing")
		capabilities.Custom("custom", "value")
		Expect(capabilities.JSON()).To(MatchJSON(`{
			"browserName": "some-browser",
			"version": "v100",
			"platform": "some-os",
			"enabledThing": true,
			"disabledThing": false,
			"custom": "value"
		}`))
	})

	Context("when the provided options cannot be converted to JSON", func() {
		It("should return an error", func() {
			capabilities.Custom("custom", func() {})
			_, err := capabilities.JSON()
			Expect(err).To(MatchError("json: unsupported type: func()"))
		})
	})
})
