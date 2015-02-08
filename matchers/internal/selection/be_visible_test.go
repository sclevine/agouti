package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/selection"
)

var _ = Describe("BeVisibleMatcher", func() {
	var (
		matcher   *BeVisibleMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.StringCall.ReturnString = "CSS: #selector"
		matcher = &BeVisibleMatcher{}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			Context("when the element is visible", func() {
				It("should successfully return true", func() {
					selection.VisibleCall.ReturnVisible = true
					Expect(matcher.Match(selection)).To(BeTrue())
				})
			})

			Context("when the element is not visible", func() {
				It("should successfully return false", func() {
					selection.VisibleCall.ReturnVisible = false
					Expect(matcher.Match(selection)).To(BeFalse())
				})
			})

			Context("when determining whether the element is visible fails", func() {
				It("should return an error", func() {
					selection.VisibleCall.Err = errors.New("some error")
					_, err := matcher.Match(selection)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a selection", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a selection")
				Expect(err).To(MatchError("BeVisible matcher requires a Selection.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			selection.VisibleCall.ReturnVisible = false
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' to be visible"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			selection.VisibleCall.ReturnVisible = true
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' not to be visible"))
		})
	})
})
