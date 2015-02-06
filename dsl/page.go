package dsl

import "github.com/sclevine/agouti/core"

// Destroy is comparable to Expect(page.Destroy()).To(Succeed())
func Destroy(page core.Page) {
	checkFailure(page.Destroy())
}

// Navigate is comparable to Expect(page.Navigate()).To(Succeed())
func Navigate(page core.Page) {
	checkFailure(page.Navigate())
}

// SetCookie is comparable to Expect(page.SetCookie(cookie)).To(Succeed())
func SetCookie(page core.Page, cookie core.WebCookie) {
	checkFailure(page.SetCookie(cookie))
}

// DeleteCookie is comparable to Expect(page.DeleteCookie("cookie-name")).To(Succeed())
func DeleteCookie(page core.Page, name string) {
	checkFailure(page.DeleteCookie(name))
}

// ClearCookies is comparable to Expect(page.ClearCookies()).To(Succeed())
func ClearCookies(page core.Page) {
	checkFailure(page.ClearCookies())
}

// Size is comparable to Expect(page.Size(windowWidth, windowHeight)).To(Succeed())
func Size(page core.Page, width, height int) {
	checkFailure(page.Size(width, height))
}

// Screenshot is comparable to Expect(page.Screenshot("screenshot-file.png")).To(Succeed())
func Screenshot(page core.Page, filename string) {
	checkFailure(page.Screenshot(filename))
}

// RunScript is comparable to Expect(page.RunScript(script, args, &result)).To(Succeed())
func RunScript(page core.Page, body string, arguments map[string]interface{}, result interface{}) {
	checkFailure(page.RunScript(body, arguments, result))
}

// EnterPopupText is comparable to Expect(page.EnterPopupText("some text")).To(Succeed())
func EnterPopupText(page core.Page, text string) {
	checkFailure(page.EnterPopupText(text))
}

// ConfirmPopup is comparable to Expect(page.ConfirmPopup()).To(Succeed())
func ConfirmPopup(page core.Page) {
	checkFailure(page.ConfirmPopup())
}

// CancelPopup is comparable to Expect(page.CancelPopup()).To(Succeed())
func CancelPopup(page core.Page) {
	checkFailure(page.CancelPopup())
}

// Forward is comparable to Expect(page.Forward()).To(Succeed())
func Forward(page core.Page) {
	checkFailure(page.Forward())
}

// Back is comparable to Expect(page.Back()).To(Succeed())
func Back(page core.Page) {
	checkFailure(page.Back())
}

// Refresh is comparable to Expect(page.Refresh()).To(Succeed())
func Refresh(page core.Page) {
	checkFailure(page.Refresh())
}

// SwitchToParentFrame is comparable to Expect(page.SwitchToParentFrame()).To(Succeed())
func SwitchToParentFrame(page core.Page) {
	checkFailure(page.SwitchToParentFrame())
}

// SwitchToRootFrame is comparable to Expect(page.SwitchToRootFrame()).To(Succeed())
func SwitchToRootFrame(page core.Page) {
	checkFailure(page.SwitchToRootFrame())
}

// SwitchToWindow is comparable to Expect(page.SwitchToWindow("window name")).To(Succeed())
func SwitchToWindow(page core.Page, name string) {
	checkFailure(page.SwitchToWindow(name))
}

// NextWindow is comparable to Expect(page.NextWindow()).To(Succeed())
func NextWindow(page core.Page) {
	checkFailure(page.NextWindow())
}

// CloseWindow is comparable to Expect(page.CloseWindow()).To(Succeed())
func CloseWindow(page core.Page) {
	checkFailure(page.CloseWindow())
}
