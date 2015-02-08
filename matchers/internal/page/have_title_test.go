package page_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/page"
)

var _ = Describe("HaveTitleMatcher", func() {
	var (
		matcher *HaveTitleMatcher
		page    *mocks.Page
	)

	BeforeEach(func() {
		page = &mocks.Page{}
		matcher = &HaveTitleMatcher{ExpectedTitle: "Some Title"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a page", func() {
			Context("when the expected title matches the actual title", func() {
				It("should successfully return true", func() {
					page.TitleCall.ReturnTitle = "Some Title"
					Expect(matcher.Match(page)).To(BeTrue())
				})
			})

			Context("when the expected title does not match the actual title", func() {
				It("should successfully return false", func() {
					page.TitleCall.ReturnTitle = "Some Other Title"
					Expect(matcher.Match(page)).To(BeFalse())
				})
			})

			Context("when retrieving the page title fails", func() {
				It("should return an error", func() {
					page.TitleCall.Err = errors.New("some error")
					_, err := matcher.Match(page)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a page", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a page")
				Expect(err).To(MatchError("HaveTitle matcher requires a Page.  Got:\n    <string>: not a page"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			page.TitleCall.ReturnTitle = "Some Other Title"
			matcher.Match(page)
			message := matcher.FailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page to have title matching\n    Some Title"))
			Expect(message).To(ContainSubstring("but found\n    Some Other Title"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			page.TitleCall.ReturnTitle = "Some Title"
			matcher.Match(page)
			message := matcher.NegatedFailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page not to have title matching\n    Some Title"))
			Expect(message).To(ContainSubstring("but found\n    Some Title"))
		})
	})
})
