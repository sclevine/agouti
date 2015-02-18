package matchers_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("Page Matchers", func() {
	var page *mocks.Page

	BeforeEach(func() {
		page = &mocks.Page{}
	})

	Describe("#HaveTitle", func() {
		It("should call the page#HaveTitle matcher", func() {
			page.TitleCall.ReturnTitle = "Some Title"
			Expect(page).To(HaveTitle("Some Title"))
			Expect(page).NotTo(HaveTitle("Some Other Title"))
		})
	})

	Describe("#HaveURL", func() {
		It("should call the page#HaveURL matcher", func() {
			page.URLCall.ReturnURL = "some/url"
			Expect(page).To(HaveURL("some/url"))
			Expect(page).NotTo(HaveURL("some/other/url"))
		})
	})

	Describe("#HavePopupText", func() {
		It("should call the page#HavePopupText matcher", func() {
			page.PopupTextCall.ReturnText = "some text"
			Expect(page).To(HavePopupText("some text"))
			Expect(page).NotTo(HavePopupText("some other text"))
		})
	})

	Describe("#HaveLoggedError", func() {
		It("should call the page#HaveLoggedError matcher", func() {
			page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some log", "", "WARNING", time.Time{}}}
			Expect(page).To(HaveLoggedError("some log"))
			Expect(page).NotTo(HaveLoggedError("some other log"))
		})
	})

	Describe("#HaveLoggedInfo", func() {
		It("should call the page#HaveLoggedInfo matcher", func() {
			page.ReadAllLogsCall.ReturnLogs = []agouti.Log{agouti.Log{"some log", "", "INFO", time.Time{}}}
			Expect(page).To(HaveLoggedInfo("some log"))
			Expect(page).NotTo(HaveLoggedInfo("some other log"))
		})
	})
})
