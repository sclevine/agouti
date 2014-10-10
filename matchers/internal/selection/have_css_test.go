package selection_test

import (
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveCSS", func() {
	var (
		matcher   *HaveCSSMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.SelectorCall.ReturnSelector = "#selector"
		matcher = &HaveCSSMatcher{"some-property", "some value"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			It("requests the provided page property", func() {
				matcher.Match(selection)
				Expect(selection.CSSCall.Property).To(Equal("some-property"))
			})

			Context("when the expected property value matches the actual property value", func() {
				BeforeEach(func() {
					selection.CSSCall.ReturnValue = "some value"
				})

				It("returns true", func() {
					success, _ := matcher.Match(selection)
					Expect(success).To(BeTrue())
				})

				It("does not return an error", func() {
					_, err := matcher.Match(selection)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the expected property value does not match the actual property value", func() {
				BeforeEach(func() {
					selection.CSSCall.ReturnValue = "some other value"
				})

				It("returns true", func() {
					success, _ := matcher.Match(selection)
					Expect(success).To(BeFalse())
				})

				It("does not return an error", func() {
					_, err := matcher.Match(selection)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the actual object is not a selection", func() {
			It("returns an error", func() {
				_, err := matcher.Match("not a selection")
				Expect(err).To(MatchError("HaveCSS matcher requires a Selection or Page.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("return a failure message", func() {
			selection.CSSCall.ReturnValue = "some other value"
			message := matcher.FailureMessage(selection)
			Expect(message).To(ContainSubstring("<selection.Selector>: #selector"))
			Expect(message).To(ContainSubstring("to have CSS matching"))
			Expect(message).To(ContainSubstring(`<string>: some-property: "some value"`))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("return a negated failure message", func() {
			selection.CSSCall.ReturnValue = "some other value"
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(ContainSubstring("<selection.Selector>: #selector"))
			Expect(message).To(ContainSubstring("not to have CSS matching"))
			Expect(message).To(ContainSubstring(`<string>: some-property: "some value"`))
		})
	})
})
