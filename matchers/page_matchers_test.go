package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/internal/mocks"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Page Matchers", func() {
	var page *mocks.Page

	BeforeEach(func() {
		page = &mocks.Page{}
	})

	Describe("#HaveTitle", func() {
		It("calls the page#HaveTitle matcher", func() {
			page.TitleCall.ReturnTitle = "Some Title"
			Expect(page).To(HaveTitle("Some Title"))
			Expect(page).NotTo(HaveTitle("Some Other Title"))
		})
	})
})
