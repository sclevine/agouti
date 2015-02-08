package page_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/page"
)

var _ = Describe("HaveURLMatcher", func() {
	var (
		matcher *HaveURLMatcher
		page    *mocks.Page
	)

	BeforeEach(func() {
		page = &mocks.Page{}
		page.URLCall.ReturnURL = "Some URL"
		matcher = &HaveURLMatcher{ExpectedURL: "Some URL"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a page", func() {
			Context("when the expected URL matches the actual URL", func() {
				It("should successfully return true", func() {
					page.URLCall.ReturnURL = "Some URL"
					Expect(matcher.Match(page)).To(BeTrue())
				})
			})

			Context("when the expected URL does not match the actual URL", func() {
				It("should successfully return false", func() {
					page.URLCall.ReturnURL = "Some Other URL"
					Expect(matcher.Match(page)).To(BeFalse())
				})
			})

			Context("when retrieving the page URL fails", func() {
				It("should return an error", func() {
					page.URLCall.Err = errors.New("some error")
					_, err := matcher.Match(page)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a page", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a page")
				Expect(err).To(MatchError("HaveURL matcher requires a Page.  Got:\n    <string>: not a page"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			page.URLCall.ReturnURL = "Some Other URL"
			matcher.Match(page)
			message := matcher.FailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page to have URL matching\n    Some URL"))
			Expect(message).To(ContainSubstring("but found\n    Some Other URL"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			page.URLCall.ReturnURL = "Some URL"
			matcher.Match(page)
			message := matcher.NegatedFailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page not to have URL matching\n    Some URL"))
			Expect(message).To(ContainSubstring("but found\n    Some URL"))
		})
	})
})
