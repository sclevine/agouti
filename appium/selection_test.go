package appium_test

import (
	"github.com/sclevine/agouti/appium"
)

var _ = Describe("Selection", func() {
	var dev *appium.Device
	var baseSel *appium.Selection

	BeforeEach(func() {
		session := &mockMobileSession{}
		dev = appium.NewTestDevice(session)
		baseSel = dev.Find(".root")
	})

	It("should successfully Find", func() {
		Expect(baseSel.Find(".go#css").String()).To(Equal(`selection 'CSS: .root [single] | CSS: .go#css [single]'`))
	})

	It("should successfully FindByID", func() {
		Expect(baseSel.FindByID("an-id").String()).To(Equal(`selection 'CSS: .root [single] | ID: an-id [single]'`))
	})

	It("should successfully FindByXPath", func() {
		Expect(baseSel.FindByXPath("//node").String()).To(Equal(`selection 'CSS: .root [single] | XPath: //node [single]'`))
	})

	It("should successfully FindByA11yID", func() {
		Expect(baseSel.FindByA11yID("this-id").String()).To(Equal(`selection 'CSS: .root [single] | Accessibility ID: this-id [single]'`))
	})

	It("should successfully FindByAndroidUI", func() {
		Expect(baseSel.FindByAndroidUI("this-ui").String()).To(Equal(`selection 'CSS: .root [single] | Android UIAut.: this-ui [single]'`))
	})

	It("should successfully FindByiOSUI", func() {
		Expect(baseSel.FindByiOSUI("this-ui").String()).To(Equal(`selection 'CSS: .root [single] | iOS UIAut.: this-ui [single]'`))
	})

	It("should successfully FindByLink", func() {
		Expect(baseSel.FindByLink("a link").String()).To(Equal(`selection 'CSS: .root [single] | Link: "a link" [single]'`))
	})
})
