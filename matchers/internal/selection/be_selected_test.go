package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"
)

var _ = Describe("BeSelectedMatcher", func() {
	var (
		matcher   *BeSelectedMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.StringCall.ReturnString = "CSS: #selector"
		matcher = &BeSelectedMatcher{}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			Context("when the element is selected", func() {
				It("should successfully return true", func() {
					selection.SelectedCall.ReturnSelected = true
					Expect(matcher.Match(selection)).To(BeTrue())
				})
			})

			Context("when the element is not selected", func() {
				It("should successfully return false", func() {
					selection.SelectedCall.ReturnSelected = false
					Expect(matcher.Match(selection)).To(BeFalse())
				})
			})

			Context("when determining whether the element is selected fails", func() {
				It("should return an error", func() {
					selection.SelectedCall.Err = errors.New("some error")
					_, err := matcher.Match(selection)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a selection", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a selection")
				Expect(err).To(MatchError("BeSelected matcher requires a Selection.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			selection.SelectedCall.ReturnSelected = false
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' to be selected"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			selection.SelectedCall.ReturnSelected = true
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' not to be selected"))
		})
	})
})
