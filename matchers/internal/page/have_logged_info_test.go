package page_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core"
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
				Expect(page.ReadLogsCall.LogType).To(Equal("browser"))
				Expect(page.ReadLogsCall.All).To(BeTrue())
			})

			Context("when the expected log has been logged with the INFO level", func() {
				It("should successfully return true", func() {
					page.ReadLogsCall.ReturnLogs = []core.Log{core.Log{"some log", "", "INFO", time.Time{}}}
					success, err := matcher.Match(page)
					Expect(success).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the expected log has been logged with any other level", func() {
				It("should successfully return false", func() {
					page.ReadLogsCall.ReturnLogs = []core.Log{core.Log{"some log", "", "WARNING", time.Time{}}}
					success, err := matcher.Match(page)
					Expect(success).To(BeFalse())
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the expected log has not been logged", func() {
				It("should successfully return false", func() {
					page.ReadLogsCall.ReturnLogs = []core.Log{core.Log{"some other log", "", "INFO", time.Time{}}}
					success, err := matcher.Match(page)
					Expect(success).To(BeFalse())
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when any log is expected", func() {
				BeforeEach(func() {
					matcher.ExpectedMessage = ""
				})

				Context("when any error log is logged", func() {
					It("should successfully return true", func() {
						page.ReadLogsCall.ReturnLogs = []core.Log{core.Log{"some log", "", "INFO", time.Time{}}}
						success, err := matcher.Match(page)
						Expect(success).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())
					})
				})

				Context("when no error logs are logged", func() {
					It("should successfully return false", func() {
						page.ReadLogsCall.ReturnLogs = []core.Log{core.Log{"some log", "", "WARNING", time.Time{}}}
						success, err := matcher.Match(page)
						Expect(success).To(BeFalse())
						Expect(err).NotTo(HaveOccurred())
					})
				})
			})

			Context("when retrieving the logs fails", func() {
				It("should return an error", func() {
					page.ReadLogsCall.Err = errors.New("some error")
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
			page.ReadLogsCall.ReturnLogs = []core.Log{core.Log{"some log", "", "WARNING", time.Time{}}}
			matcher.Match(page)
			message := matcher.FailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page to have info log matching\n    some log"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			page.ReadLogsCall.ReturnLogs = []core.Log{core.Log{"some other log", "", "WARNING", time.Time{}}}
			matcher.Match(page)
			message := matcher.NegatedFailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page not to have info log matching\n    some log"))
		})
	})
})
