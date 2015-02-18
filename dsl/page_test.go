package dsl_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/dsl"
	"github.com/sclevine/agouti/dsl/internal/mocks"
)

var _ = Describe("Page", func() {
	var (
		page        *mocks.Page
		failMessage string
	)

	BeforeEach(func() {
		failMessage = ""
		RegisterAgoutiFailHandler(func(message string, callerSkip ...int) {
			failMessage = message
			ExpectWithOffset(3, callerSkip[0]).To(Equal(2))
			panic("Failed to catch test panic.")
		})
		page = &mocks.Page{}
	})

	AfterEach(func() {
		RegisterAgoutiFailHandler(Fail)
	})

	Describe(".Destroy", func() {
		It("should call page.Destroy", func() {
			Destroy(page)
			Expect(page.DestroyCall.Called).To(BeTrue())
		})

		It("should fail when page.Destroy returns an error", func() {
			page.DestroyCall.Err = errors.New("some error")
			Expect(func() { Destroy(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".Navigate", func() {
		It("should call page.Navigate", func() {
			Navigate(page, "some URL")
			Expect(page.NavigateCall.URL).To(Equal("some URL"))
		})

		It("should fail when page.Navigate returns an error", func() {
			page.NavigateCall.Err = errors.New("some error")
			Expect(func() { Navigate(page, "some URL") }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".SetCookie", func() {
		It("should call page.SetCookie", func() {
			cookie := agouti.NewCookie("some", "cookie")
			SetCookie(page, cookie)
			Expect(page.SetCookieCall.Cookie).To(Equal(cookie))
		})

		It("should fail when page.SetCookie returns an error", func() {
			page.SetCookieCall.Err = errors.New("some error")
			Expect(func() { SetCookie(page, agouti.NewCookie("some", "cookie")) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".DeleteCookie", func() {
		It("should call page.DeleteCookie", func() {
			DeleteCookie(page, "some cookie name")
			Expect(page.DeleteCookieCall.Name).To(Equal("some cookie name"))
		})

		It("should fail when page.DeleteCookie returns an error", func() {
			page.DeleteCookieCall.Err = errors.New("some error")
			Expect(func() { DeleteCookie(page, "some cookie name") }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".ClearCookies", func() {
		It("should call page.ClearCookies", func() {
			ClearCookies(page)
			Expect(page.ClearCookiesCall.Called).To(BeTrue())
		})

		It("should fail when page.ClearCookies returns an error", func() {
			page.ClearCookiesCall.Err = errors.New("some error")
			Expect(func() { ClearCookies(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".Size", func() {
		It("should call page.Size", func() {
			Size(page, 100, 200)
			Expect(page.SizeCall.Width).To(Equal(100))
			Expect(page.SizeCall.Height).To(Equal(200))
		})

		It("should fail when page.Size returns an error", func() {
			page.SizeCall.Err = errors.New("some error")
			Expect(func() { Size(page, 100, 200) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".Screenshot", func() {
		It("should call page.Screenshot", func() {
			Screenshot(page, "some filename")
			Expect(page.ScreenshotCall.Filename).To(Equal("some filename"))
		})

		It("should fail when page.Screenshot returns an error", func() {
			page.ScreenshotCall.Err = errors.New("some error")
			Expect(func() { Screenshot(page, "some filename") }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".RunScript", func() {
		It("should call page.RunScript", func() {
			arguments := map[string]interface{}{"some": "arguments"}
			RunScript(page, "some body", arguments, "some result")
			Expect(page.RunScriptCall.Body).To(Equal("some body"))
			Expect(page.RunScriptCall.Arguments).To(Equal(arguments))
			Expect(page.RunScriptCall.Result).To(Equal("some result"))
		})

		It("should fail when page.RunScript returns an error", func() {
			page.RunScriptCall.Err = errors.New("some error")
			Expect(func() { RunScript(page, "some body", nil, nil) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".EnterPopupText", func() {
		It("should call page.EnterPopupText", func() {
			EnterPopupText(page, "some text")
			Expect(page.EnterPopupTextCall.Text).To(Equal("some text"))
		})

		It("should fail when page.EnterPopupText returns an error", func() {
			page.EnterPopupTextCall.Err = errors.New("some error")
			Expect(func() { EnterPopupText(page, "some text") }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".ConfirmPopup", func() {
		It("should call page.ClearCookies", func() {
			ClearCookies(page)
			Expect(page.ClearCookiesCall.Called).To(BeTrue())
		})

		It("should fail when page.ClearCookies returns an error", func() {
			page.ClearCookiesCall.Err = errors.New("some error")
			Expect(func() { ClearCookies(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".CancelPopup", func() {
		It("should call page.CancelPopup", func() {
			CancelPopup(page)
			Expect(page.CancelPopupCall.Called).To(BeTrue())
		})

		It("should fail when page.CancelPopup returns an error", func() {
			page.CancelPopupCall.Err = errors.New("some error")
			Expect(func() { CancelPopup(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".Forward", func() {
		It("should call page.Forward", func() {
			Forward(page)
			Expect(page.ForwardCall.Called).To(BeTrue())
		})

		It("should fail when page.Forward returns an error", func() {
			page.ForwardCall.Err = errors.New("some error")
			Expect(func() { Forward(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".Back", func() {
		It("should call page.Back", func() {
			Back(page)
			Expect(page.BackCall.Called).To(BeTrue())
		})

		It("should fail when page.Back returns an error", func() {
			page.BackCall.Err = errors.New("some error")
			Expect(func() { Back(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".Refresh", func() {
		It("should call page.Refresh", func() {
			Refresh(page)
			Expect(page.RefreshCall.Called).To(BeTrue())
		})

		It("should fail when page.Refresh returns an error", func() {
			page.RefreshCall.Err = errors.New("some error")
			Expect(func() { Refresh(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".SwitchToParentFrame", func() {
		It("should call page.SwitchToParentFrame", func() {
			SwitchToParentFrame(page)
			Expect(page.SwitchToParentFrameCall.Called).To(BeTrue())
		})

		It("should fail when page.SwitchToParentFrame returns an error", func() {
			page.SwitchToParentFrameCall.Err = errors.New("some error")
			Expect(func() { SwitchToParentFrame(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".SwitchToRootFrame", func() {
		It("should call page.SwitchToRootFrame", func() {
			SwitchToRootFrame(page)
			Expect(page.SwitchToRootFrameCall.Called).To(BeTrue())
		})

		It("should fail when page.SwitchToRootFrame returns an error", func() {
			page.SwitchToRootFrameCall.Err = errors.New("some error")
			Expect(func() { SwitchToRootFrame(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".SwitchToWindow", func() {
		It("should call page.SwitchToWindow", func() {
			SwitchToWindow(page, "some window name")
			Expect(page.SwitchToWindowCall.Name).To(Equal("some window name"))
		})

		It("should fail when page.SwitchToWindow returns an error", func() {
			page.SwitchToWindowCall.Err = errors.New("some error")
			Expect(func() { SwitchToWindow(page, "some window name") }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".NextWindow", func() {
		It("should call page.NextWindow", func() {
			NextWindow(page)
			Expect(page.NextWindowCall.Called).To(BeTrue())
		})

		It("should fail when page.NextWindow returns an error", func() {
			page.NextWindowCall.Err = errors.New("some error")
			Expect(func() { NextWindow(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})

	Describe(".CloseWindow", func() {
		It("should call page.CloseWindow", func() {
			CloseWindow(page)
			Expect(page.CloseWindowCall.Called).To(BeTrue())
		})

		It("should fail when page.CloseWindow returns an error", func() {
			page.CloseWindowCall.Err = errors.New("some error")
			Expect(func() { CloseWindow(page) }).To(Panic())
			Expect(failMessage).To(Equal("Agouti failure: some error"))
		})
	})
})
