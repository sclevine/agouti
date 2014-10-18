package page_test

import (
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/page"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveTitleMatcher", func() {
	var (
		matcher *HaveTitleMatcher
		page    *mocks.Page
	)

	BeforeEach(func() {
		page = &mocks.Page{}
		page.TitleCall.ReturnTitle = "Some Title"
		matcher = &HaveTitleMatcher{ExpectedTitle: "Some Title"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a types.PageOnly", func() {
			Context("when the expected title matches the actual title", func() {
				BeforeEach(func() {
					page.TitleCall.ReturnTitle = "Some Title"
				})

				It("should return true", func() {
					success, _ := matcher.Match(page)
					Expect(success).To(BeTrue())
				})

				It("should not return an error", func() {
					_, err := matcher.Match(page)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the expected title does not match the actual title", func() {
				BeforeEach(func() {
					page.TitleCall.ReturnTitle = "Some Other Title"
				})

				It("should return false", func() {
					success, _ := matcher.Match(page)
					Expect(success).To(BeFalse())
				})

				It("should not return an error", func() {
					_, err := matcher.Match(page)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the actual object is not a types.PageOnly", func() {
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
