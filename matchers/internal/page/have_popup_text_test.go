package page_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/page"
)

var _ = Describe("HavePopupTextMatcher", func() {
	var (
		matcher *HavePopupTextMatcher
		page    *mocks.Page
	)

	BeforeEach(func() {
		page = &mocks.Page{}
		page.PopupTextCall.ReturnText = "some text"
		matcher = &HavePopupTextMatcher{ExpectedText: "some text"}
	})

	Describe("#Match", func() {
		Context("when the actual object is page", func() {
			Context("when the expected text matches the actual text", func() {
				It("should successfully return true", func() {
					page.PopupTextCall.ReturnText = "some text"
					Expect(matcher.Match(page)).To(BeTrue())
				})
			})

			Context("when the expected text does not match the actual text", func() {
				It("should successfully return false", func() {
					page.PopupTextCall.ReturnText = "some other text"
					Expect(matcher.Match(page)).To(BeFalse())
				})
			})

			Context("when retrieving the popup text fails", func() {
				It("should return an error", func() {
					page.PopupTextCall.Err = errors.New("some error")
					_, err := matcher.Match(page)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a page", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a page")
				Expect(err).To(MatchError("HavePopupText matcher requires a Page.  Got:\n    <string>: not a page"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			page.PopupTextCall.ReturnText = "some other text"
			matcher.Match(page)
			message := matcher.FailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page to have popup text matching\n    some text"))
			Expect(message).To(ContainSubstring("but found\n    some other text"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			page.PopupTextCall.ReturnText = "some text"
			matcher.Match(page)
			message := matcher.NegatedFailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page not to have popup text matching\n    some text"))
			Expect(message).To(ContainSubstring("but found\n    some text"))
		})
	})
})
