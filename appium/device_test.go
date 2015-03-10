package appium_test

import "github.com/sclevine/agouti/appium"

var _ = Describe("Device action methods", func() {
})

var _ = Describe("Device selection methods", func() {
	var dev *appium.Device

	BeforeEach(func() {
		session := &mockMobileSession{}
		dev = appium.NewTestDevice(session)
	})

	It("should successfully Find", func() {
		Expect(dev.Find(".go#css").String()).To(Equal(`selection 'CSS: .go#css [single]'`))
	})

	It("should successfully FindByXPath", func() {
		Expect(dev.FindByXPath("//node").String()).To(Equal(`selection 'XPath: //node [single]'`))
	})

	It("should successfully FindByA11yID", func() {
		Expect(dev.FindByA11yID("this-id").String()).To(Equal(`selection 'Accessibility ID: this-id [single]'`))
	})

	It("should successfully FindByAndroidUI", func() {
		Expect(dev.FindByAndroidUI("this-ui").String()).To(Equal(`selection 'Android UIAut.: this-ui [single]'`))
	})

	It("should successfully FindByiOSUI", func() {
		Expect(dev.FindByiOSUI("this-ui").String()).To(Equal(`selection 'iOS UIAut.: this-ui [single]'`))
	})

	It("should successfully FindByLink", func() {
		Expect(dev.FindByLink("a link").String()).To(Equal(`selection 'Link: "a link" [single]'`))
	})
})
