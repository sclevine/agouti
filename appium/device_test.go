package appium

import (
	"github.com/onsi/gomega/types"
	"github.com/sclevine/agouti"
)

var _ = Describe("Device", func() {
	var dev *Device
	var baseSel *Selection

	BeforeEach(func() {
		session := &mockMobileSession{}
		dev = newDevice(session, &agouti.Page{})
		baseSel = dev.Find(".root")
	})

	bePrefixed := func(prefix string) types.GomegaMatcher {
		return Equal("selection '" + prefix + "'")
	}
	beBasePrefixed := func(prefix string) types.GomegaMatcher {
		return Equal("selection 'CSS: .root [single] | " + prefix + "'")
	}

	expectDeviceAndSelectionOutput := func(deviceSelection *Selection, chainedSelection *Selection, output string) {
		Expect(deviceSelection.String()).To(bePrefixed(output))
		Expect(chainedSelection.String()).To(beBasePrefixed(output))
	}

	Describe("#Find... methods", func() {
		It("should Find", func() {
			expectDeviceAndSelectionOutput(
				dev.Find(".go#css"),
				baseSel.Find(".go#css"),
				`CSS: .go#css [single]`,
			)
		})

		It("should FindByXPath", func() {
			expectDeviceAndSelectionOutput(
				dev.FindByXPath("//node"),
				baseSel.FindByXPath("//node"),
				`XPath: //node [single]`,
			)
		})

		It("should FindByA11yID", func() {
			expectDeviceAndSelectionOutput(
				dev.FindByA11yID("this-id"),
				baseSel.FindByA11yID("this-id"),
				`Accessibility ID: this-id [single]`,
			)
		})

		It("should FindByAndroidUI", func() {
			expectDeviceAndSelectionOutput(
				dev.FindByAndroidUI("this-ui"),
				baseSel.FindByAndroidUI("this-ui"),
				`Android UIAut.: this-ui [single]`,
			)
		})

		It("should FindByiOSUI", func() {
			expectDeviceAndSelectionOutput(
				dev.FindByiOSUI("this-ui"),
				baseSel.FindByiOSUI("this-ui"),
				`iOS UIAut.: this-ui [single]`,
			)
		})

		It("should FindByLink", func() {
			expectDeviceAndSelectionOutput(
				dev.FindByLink("a link"),
				baseSel.FindByLink("a link"),
				`Link: "a link" [single]`,
			)
		})

	})
})
