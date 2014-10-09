package selection_test

import (
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContainTextMatcher", func() {
	var (
		matcher   *ContainTextMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.SelectorCall.ReturnSelector = "#selector"
		matcher = &ContainTextMatcher{"some text"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			Context("when the expected text matches the actual text", func() {
				BeforeEach(func() {
					selection.TextCall.ReturnText = "some text"
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

			Context("when the expected text does not match the actual text", func() {
				BeforeEach(func() {
					selection.TextCall.ReturnText = "some other text"
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
				Expect(err).To(MatchError("ContainText matcher requires a Selection or Page.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("return a failure message", func() {
			selection.TextCall.ReturnText = "some other text"
			message := matcher.FailureMessage(selection)
			Expect(message).To(ContainSubstring("<selection.SelectorText>: #selector"))
			Expect(message).To(ContainSubstring("to have text matching"))
			Expect(message).To(ContainSubstring("<string>: some text"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("return a negated failure message", func() {
			selection.TextCall.ReturnText = "some other text"
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(ContainSubstring("<selection.SelectorText>: #selector"))
			Expect(message).To(ContainSubstring("not to have text matching"))
			Expect(message).To(ContainSubstring("<string>: some text"))
		})
	})
})
