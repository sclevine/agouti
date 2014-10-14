package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Selection Matchers", func() {
	var selection *mocks.Selection

	BeforeEach(func() {
		selection = &mocks.Selection{}
	})

	Describe("#HaveText", func() {
		It("calls the selection#HaveText matcher", func() {
			selection.TextCall.ReturnText = "some text"
			Expect(selection).To(HaveText("some text"))
			Expect(selection).NotTo(HaveText("some other text"))
		})
	})

	Describe("#HaveAttribute", func() {
		It("calls the selection#HaveAttribute matcher", func() {
			selection.AttributeCall.ReturnValue = "some value"
			Expect(selection).To(HaveAttribute("some-attribute", "some value"))
			Expect(selection).NotTo(HaveAttribute("some-attribute", "some other value"))
		})
	})

	Describe("#HaveCSS", func() {
		It("calls the selection#HaveCSS matcher", func() {
			selection.CSSCall.ReturnValue = "some value"
			Expect(selection).To(HaveCSS("some-property", "some value"))
			Expect(selection).NotTo(HaveCSS("some-property", "some other value"))
		})
	})

	Describe("#BeSelected", func() {
		It("calls the selection#BeSelected matcher", func() {
			selection.SelectedCall.ReturnSelected = true
			Expect(selection).To(BeSelected())
			selection.SelectedCall.ReturnSelected = false
			Expect(selection).NotTo(BeSelected())
		})
	})

	Describe("#BeVisible", func() {
		It("calls the selection#BeVisible matcher", func() {
			selection.VisibleCall.ReturnVisible = true
			Expect(selection).To(BeVisible())
			selection.VisibleCall.ReturnVisible = false
			Expect(selection).NotTo(BeVisible())
		})
	})

	Describe("#BeFound", func() {
		It("calls the selection#BeFound matcher", func() {
			selection.CountCall.ReturnCount = 1
			Expect(selection).To(BeFound())
			selection.CountCall.ReturnCount = 0
			Expect(selection).NotTo(BeFound())
		})
	})
})
