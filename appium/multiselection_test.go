package appium_test

import "github.com/sclevine/agouti/appium"

var _ = Describe("MultiSelection", func() {
	var dev *appium.Device

	BeforeEach(func() {
		session := &mockMobileSession{}
		dev = appium.NewTestDevice(session)
	})

	It("should successfully chain CSS through Find", func() {
		Expect(dev.All(".root").Find(".go#css").String()).To(Equal(`selection 'CSS: .root .go#css [single]'`))
	})

	It("should successfully Find, without chaining because of the At() call", func() {
		Expect(dev.All(".root").At(0).Find(".go#css").String()).To(Equal(`selection 'CSS: .root [0] | CSS: .go#css [single]'`))
	})

})
