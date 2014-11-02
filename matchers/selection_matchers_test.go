package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("Selection Matchers", func() {
	var selection *mocks.Selection

	BeforeEach(func() {
		selection = &mocks.Selection{}
	})

	Describe("#HaveText", func() {
		It("should call the selection#HaveText matcher", func() {
			selection.TextCall.ReturnText = "some text"
			Expect(selection).To(HaveText("some text"))
			Expect(selection).NotTo(HaveText("some other text"))
		})
	})

	Describe("#MatchText", func() {
		It("should call the selection#MatchText matcher", func() {
			selection.TextCall.ReturnText = "some text"
			Expect(selection).To(MatchText("s[^t]+text"))
			Expect(selection).NotTo(MatchText("so*text"))
		})
	})

	Describe("#HaveAttribute", func() {
		It("should call the selection#HaveAttribute matcher", func() {
			selection.AttributeCall.ReturnValue = "some value"
			Expect(selection).To(HaveAttribute("some-attribute", "some value"))
			Expect(selection).NotTo(HaveAttribute("some-attribute", "some other value"))
		})
	})

	Describe("#HaveCSS", func() {
		It("should call the selection#HaveCSS matcher", func() {
			selection.CSSCall.ReturnValue = "some value"
			Expect(selection).To(HaveCSS("some-property", "some value"))
			Expect(selection).NotTo(HaveCSS("some-property", "some other value"))
		})
	})

	Describe("#BeSelected", func() {
		It("should call the selection#BeSelected matcher", func() {
			selection.SelectedCall.ReturnSelected = true
			Expect(selection).To(BeSelected())
			selection.SelectedCall.ReturnSelected = false
			Expect(selection).NotTo(BeSelected())
		})
	})

	Describe("#BeVisible", func() {
		It("should call the selection#BeVisible matcher", func() {
			selection.VisibleCall.ReturnVisible = true
			Expect(selection).To(BeVisible())
			selection.VisibleCall.ReturnVisible = false
			Expect(selection).NotTo(BeVisible())
		})
	})

	Describe("#BeEnabled", func() {
		It("should call the selection#BeEnabled matcher", func() {
			selection.EnabledCall.ReturnEnabled = true
			Expect(selection).To(BeEnabled())
			selection.EnabledCall.ReturnEnabled = false
			Expect(selection).NotTo(BeEnabled())
		})
	})

	Describe("#BeActive", func() {
		It("should call the selection#BeActive matcher", func() {
			selection.ActiveCall.ReturnActive = true
			Expect(selection).To(BeActive())
			selection.ActiveCall.ReturnActive = false
			Expect(selection).NotTo(BeActive())
		})
	})

	Describe("#BeFound", func() {
		It("should call the selection#BeFound matcher", func() {
			selection.CountCall.ReturnCount = 1
			Expect(selection).To(BeFound())
			selection.CountCall.ReturnCount = 0
			Expect(selection).NotTo(BeFound())
		})
	})

	Describe("#EqualElement", func() {
		It("should call the selection#EqualElement matcher", func() {
			selection.EqualsElementCall.ReturnEquals = true
			Expect(selection).To(EqualElement(selection))
			selection.EqualsElementCall.ReturnEquals = false
			Expect(selection).NotTo(EqualElement(selection))
		})
	})
})
