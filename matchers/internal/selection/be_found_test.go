package selection_test

import (
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeFoundMatcher", func() {
	var (
		matcher   *BeFoundMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.SelectorCall.ReturnSelector = "#selector"
		matcher = &BeFoundMatcher{}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			Context("when the element is found", func() {
				BeforeEach(func() {
					selection.CountCall.ReturnCount = 1
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

			Context("when the element is not found", func() {
				BeforeEach(func() {
					selection.CountCall.ReturnCount = 0
				})

				It("returns false", func() {
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
				Expect(err).To(MatchError("BeFound matcher requires a Selection.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("returns a failure message", func() {
			selection.CountCall.ReturnCount = 0
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(Equal("Expected selection '#selector' to be found"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("returns a negated failure message", func() {
			selection.CountCall.ReturnCount = 1
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(Equal("Expected selection '#selector' not to be found"))
		})
	})
})
