package appium_test

import "github.com/sclevine/agouti/appium"

var _ = Describe("TouchAction", func() {
	session := &mockMobileSession{}

	It("should work", func() {
		ta := appium.NewTouchAction(session)
		ta.Tap().Tap()
		Expect(ta.String()).To(Equal("tap, tap"))
	})
})
