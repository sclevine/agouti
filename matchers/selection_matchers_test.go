package matchers_test

import (
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/sclevine/agouti/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Selection Matchers", func() {
	var selection *mocks.Selection

	BeforeEach(func() {
		selection = &mocks.Selection{}
	})

	Describe("#ContainText", func() {
		It("calls the selection#ContainText matcher", func() {
			selection.TextCall.ReturnText = "some text"
			Expect(selection).To(ContainText("some text"))
			Expect(selection).NotTo(ContainText("some other text"))
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
})
