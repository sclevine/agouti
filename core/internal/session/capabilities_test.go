package session_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core/internal/session"
)

var _ = Describe("Capabilities", func() {
	var capabilities Capabilities

	BeforeEach(func() {
		capabilities = Capabilities{}
	})

	Describe("#Browser", func() {
		It("should set a browserName capability", func() {
			json := capabilities.Browser("some-browser").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"browserName": "some-browser"}}`))
		})
	})

	Describe("#Version", func() {
		It("should set a version capability", func() {
			json := capabilities.Version("v100").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"version": "v100"}}`))
		})
	})

	Describe("#Platform", func() {
		It("should set a platform capability", func() {
			json := capabilities.Platform("some-os").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"platform": "some-os"}}`))
		})
	})

	Describe("#With", func() {
		It("should set the given feature to true", func() {
			json := capabilities.With("somethingEnabled").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"somethingEnabled": true}}`))
		})
	})

	Describe("#Without", func() {
		It("should set the given feature to false", func() {
			json := capabilities.Without("somethingEnabled").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"somethingEnabled": false}}`))
		})
	})

	Describe("#JSON", func() {
		It("should encode all provided options into JSON", func() {
			capabilities.Browser("some-browser").Version("v100").Platform("some-os")
			capabilities.With("enabledThing").Without("disabledThing")
			Expect(capabilities.JSON()).To(MatchJSON(`{
				"desiredCapabilities": {
					"browserName": "some-browser",
					"version": "v100",
					"platform": "some-os",
					"enabledThing": true,
					"disabledThing": false
				}
			}`))
		})
	})
})
