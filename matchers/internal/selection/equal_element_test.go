package selection_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"
)

var _ = Describe("EqualElementMatcher", func() {
	var (
		matcher        *EqualElementMatcher
		selection      *mocks.Selection
		equalSelection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		equalSelection = &mocks.Selection{}
		selection.StringCall.ReturnString = "CSS: #selector"
		equalSelection.StringCall.ReturnString = "XPath: //selector"
		matcher = &EqualElementMatcher{ExpectedSelection: equalSelection}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			It("should compare the selections for element equality", func() {
				matcher.Match(selection)
				Expect(selection.EqualsElementCall.Selection).To(Equal(equalSelection))
			})

			Context("when the expected element equals the actual element", func() {
				BeforeEach(func() {
					selection.EqualsElementCall.ReturnEquals = true
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

			Context("when the expected element does not equal the actual element", func() {
				BeforeEach(func() {
					selection.EqualsElementCall.ReturnEquals = false
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
				Expect(err).To(MatchError("EqualElement matcher requires a Selection.  Got:\n    <string>: not a selection"))
			})
		})

		Context("when the comparison fails", func() {
			It("should return an error", func() {
				selection.EqualsElementCall.Err = errors.New("some error")
				_, err := matcher.Match(selection)
				Expect(err).To(MatchError("EqualElement matcher failed to compare Selections: some error"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			selection.EqualsElementCall.ReturnEquals = false
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' to equal element\n    XPath: //selector"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			selection.EqualsElementCall.ReturnEquals = true
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' not to equal element\n    XPath: //selector"))
		})
	})
})
