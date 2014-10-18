package selection_test

import (
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveTextMatcher", func() {
	var (
		matcher   *HaveTextMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.StringCall.ReturnString = "CSS: #selector"
		matcher = &HaveTextMatcher{ExpectedText: "some text"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			Context("when the expected text matches the actual text", func() {
				BeforeEach(func() {
					selection.TextCall.ReturnText = "some text"
				})

				It("should return true", func() {
					success, _ := matcher.Match(selection)
					Expect(success).To(BeTrue())
				})

				It("should not return an error", func() {
					_, err := matcher.Match(selection)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the expected text does not equal the actual text", func() {
				BeforeEach(func() {
					selection.TextCall.ReturnText = "some other text"
				})

				It("should return false", func() {
					success, _ := matcher.Match(selection)
					Expect(success).To(BeFalse())
				})

				It("should not return an error", func() {
					_, err := matcher.Match(selection)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the actual object is not a selection", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a selection")
				Expect(err).To(MatchError("HaveText matcher requires a Selection.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			selection.TextCall.ReturnText = "some other text"
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' to have text equaling\n    some text"))
			Expect(message).To(ContainSubstring("but found\n    some other text"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			selection.TextCall.ReturnText = "some text"
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' not to have text equaling\n    some text"))
			Expect(message).To(ContainSubstring("but found\n    some text"))
		})
	})
})
