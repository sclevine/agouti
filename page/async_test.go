package page_test

import (
	. "github.com/sclevine/agouti/page"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/webdriver"
	"github.com/sclevine/agouti/mocks"
	"time"
)

var _ = Describe("Async", func() {
	var (
		async FinalSelection
		failer    *mocks.Failer
		driver    *mocks.Driver
		element   *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		failer = &mocks.Failer{}
		element = &mocks.Element{}
		async = NewPage(driver, failer.Fail).Within("#selector").ShouldEventually()
	})

	Describe("#ContainText", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []webdriver.Element{element}
			element.GetTextCall.ReturnText = "no match"
		})

		Context("if #ContainText eventually passes", func() {
			It("passes the test", func(done Done) {
				go func() {
					defer GinkgoRecover()
					Expect(func() { async.ContainText("text") }).NotTo(Panic())
					close(done)
				}()
				time.Sleep(400 * time.Millisecond)
				element.GetTextCall.ReturnText = "text"
			})
		})

		Context("if #ContainText never passes", func() {
			It("fails the test", func(done Done) {
				go func() {
					defer GinkgoRecover()
					Expect(func() { async.ContainText("text") }).To(Panic())
					Expect(failer.Message).To(Equal("After 0.5 seconds:\n FAILED"))
					Expect(failer.CallerSkip).To(Equal(100))
					close(done)
				}()
				time.Sleep(600 * time.Millisecond)
				element.GetTextCall.ReturnText = "text"
			})
		})
	})
})
