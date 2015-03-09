package appium_test

import (
	i "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/appium"
)

var _ = i.Describe("TouchAction", func() {
	session := &mockMobileSession{}

	i.It("should work", func() {
		ta := appium.NewTouchAction(session)
		ta.Tap()
		Expect(ta.String()).To(Equal("tap"))
	})
})
