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

	Describe("#Browser", func() {
		It("should set a browserName capability", func() {
			json, _ := capabilities.Browser("some-browser").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"browserName": "some-browser"}}`))
		})
	})

	Describe("#Version", func() {
		It("should set a version capability", func() {
			json, _ := capabilities.Version("v100").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"version": "v100"}}`))
		})
	})

	Describe("#Platform", func() {
		It("should set a platform capability", func() {
			json, _ := capabilities.Platform("some-os").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"platform": "some-os"}}`))
		})
	})

	Describe("#With", func() {
		It("should set the given feature to true", func() {
			json, _ := capabilities.With("somethingEnabled").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"somethingEnabled": true}}`))
		})
	})

	Describe("#Without", func() {
		It("should set the given feature to false", func() {
			json, _ := capabilities.Without("somethingEnabled").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"somethingEnabled": false}}`))
		})
	})

	Describe("#Custom", func() {
		It("should set a custom capability", func() {
			json, _ := capabilities.Custom("some", "value").JSON()
			Expect(json).To(MatchJSON(`{"desiredCapabilities": {"some": "value"}}`))
		})
	})

	Describe("#JSON", func() {
		It("should successfully encode all provided options into JSON", func() {
			capabilities.Browser("some-browser").Version("v100").Platform("some-os")
			capabilities.With("enabledThing").Without("disabledThing")
			capabilities.Custom("custom", "value")
			json, err := capabilities.JSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(json).To(MatchJSON(`{
				"desiredCapabilities": {
					"browserName": "some-browser",
					"version": "v100",
					"platform": "some-os",
					"enabledThing": true,
					"disabledThing": false,
					"custom": "value"
				}
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
})
