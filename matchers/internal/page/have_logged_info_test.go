package page_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/matchers/internal/mocks"
	. "github.com/sclevine/agouti/matchers/internal/page"
)

var _ = Describe("HaveLoggedInfoMatcher", func() {
	var (
		matcher *HaveLoggedInfoMatcher
		page    *mocks.Page
	)

	BeforeEach(func() {
		page = &mocks.Page{}
		matcher = &HaveLoggedInfoMatcher{ExpectedMessage: "some log"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a page", func() {
			It("should request all of the browser logs", func() {
				matcher.Match(page)
				Expect(page.ReadAllLogsCall.LogType).To(Equal("browser"))
			})

			Context("when the expected log has been logged with the INFO level", func() {
				It("should successfully return true", func() {
					page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some log", "", "INFO", time.Time{}}}
					Expect(matcher.Match(page)).To(BeTrue())
				})
			})

			Context("when the expected log has been logged with any other level", func() {
				It("should successfully return false", func() {
					page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some log", "", "WARNING", time.Time{}}}
					Expect(matcher.Match(page)).To(BeFalse())
				})
			})

			Context("when the expected log has not been logged", func() {
				It("should successfully return false", func() {
					page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some other log", "", "INFO", time.Time{}}}
					Expect(matcher.Match(page)).To(BeFalse())
				})
			})

			Context("when any log is expected", func() {
				BeforeEach(func() {
					matcher.ExpectedMessage = ""
				})

				Context("when any error log is logged", func() {
					It("should successfully return true", func() {
						page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some log", "", "INFO", time.Time{}}}
						Expect(matcher.Match(page)).To(BeTrue())
					})
				})

				Context("when no error logs are logged", func() {
					It("should successfully return false", func() {
						page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some log", "", "WARNING", time.Time{}}}
						Expect(matcher.Match(page)).To(BeFalse())
					})
				})
			})

			Context("when retrieving the logs fails", func() {
				It("should return an error", func() {
					page.ReadAllLogsCall.Err = errors.New("some error")
					_, err := matcher.Match(page)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a page", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a page")
				Expect(err).To(MatchError("HaveLoggedInfo matcher requires a Page.  Got:\n    <string>: not a page"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some log", "", "WARNING", time.Time{}}}
			matcher.Match(page)
			message := matcher.FailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page to have info log matching\n    some log"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some other log", "", "WARNING", time.Time{}}}
			matcher.Match(page)
			message := matcher.NegatedFailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page not to have info log matching\n    some log"))
		})
	})
})
