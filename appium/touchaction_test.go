package appium_test

import "github.com/sclevine/agouti/appium"

var _ = Describe("TouchAction", func() {
	session := &mockMobileSession{}

	It("should chain taps", func() {
		ta := appium.NewTouchAction(session)

		ta = ta.Tap().Tap()

		Expect(ta.String()).To(Equal("tap() -> tap()"))
	})

	It("should moveTo a position", func() {
		ta := appium.NewTouchAction(session)

		ta = ta.MoveToPosition(1, 2)

		Expect(ta.String()).To(Equal(`moveTo(x=1, y=2)`))
	})

	It("should chain tap and moveTo a position", func() {
		ta := appium.NewTouchAction(session)

		ta = ta.Tap().MoveToPosition(1, 2)

		Expect(ta.String()).To(Equal(`tap() -> moveTo(x=1, y=2)`))
	})
})
