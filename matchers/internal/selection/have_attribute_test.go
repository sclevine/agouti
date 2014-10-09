package selection_test

import (
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveAttributeMatcher", func() {
	var (
		matcher   *HaveAttributeMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.SelectorCall.ReturnSelector = "#selector"
		matcher = &HaveAttributeMatcher{"some-attribute", "some value"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			It("requests the provided page attribute", func() {
				matcher.Match(selection)
				Expect(selection.AttributeCall.Attribute).To(Equal("some-attribute"))
			})

			Context("when the expected attribute value matches the actual attribute value", func() {
				BeforeEach(func() {
					selection.AttributeCall.ReturnValue = "some value"
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

			Context("when the expected attribute value does not match the actual attribute value", func() {
				BeforeEach(func() {
					selection.AttributeCall.ReturnValue = "some other value"
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
				Expect(err).To(MatchError("HaveAttribute matcher requires a Selection or Page.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("return a failure message", func() {
			selection.AttributeCall.ReturnValue = "some other value"
			message := matcher.FailureMessage(selection)
			Expect(message).To(ContainSubstring("<selection.SelectorText>: #selector"))
			Expect(message).To(ContainSubstring("to have attribute matching"))
			Expect(message).To(ContainSubstring(`<string>: [some-attribute="some value"]`))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("return a negated failure message", func() {
			selection.AttributeCall.ReturnValue = "some other value"
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(ContainSubstring("<selection.SelectorText>: #selector"))
			Expect(message).To(ContainSubstring("not to have attribute matching"))
			Expect(message).To(ContainSubstring(`<string>: [some-attribute="some value"]`))
		})
	})
})
