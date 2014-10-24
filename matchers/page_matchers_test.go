package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Describe("#HavePopupText", func() {
		It("should call the page#HavePopupText matcher", func() {
			page.PopupTextCall.ReturnText = "some text"
			Expect(page).To(HavePopupText("some text"))
			Expect(page).NotTo(HavePopupText("some other text"))
		})
	})
})
