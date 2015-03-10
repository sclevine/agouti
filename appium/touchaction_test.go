package appium

var _ = Describe("TouchAction", func() {
	session := &mockMobileSession{}

	It("should work", func() {
		ta := NewTouchAction(session)
		ta.Tap().Tap()
		Expect(ta.String()).To(Equal("tap, tap"))
	})
})
