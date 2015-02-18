package dsl

import "github.com/sclevine/agouti"

type ActionPage interface {
	Destroy() error
	Navigate(url string) error
	SetCookie(cookie agouti.Cookie) error
	DeleteCookie(name string) error
	ClearCookies() error
	Size(width, height int) error
	Screenshot(filename string) error
	RunScript(body string, arguments map[string]interface{}, result interface{}) error
	EnterPopupText(text string) error
	ConfirmPopup() error
	CancelPopup() error
	Forward() error
	Back() error
	Refresh() error
	SwitchToParentFrame() error
	SwitchToRootFrame() error
	SwitchToWindow(name string) error
	NextWindow() error
	CloseWindow() error
}

// Destroy is comparable to Expect(page.Destroy()).To(Succeed())
func Destroy(page ActionPage) {
	checkFailure(page.Destroy())
}

// Navigate is comparable to Expect(page.Navigate(url)).To(Succeed())
func Navigate(page ActionPage, url string) {
	checkFailure(page.Navigate(url))
}

// SetCookie is comparable to Expect(page.SetCookie(cookie)).To(Succeed())
func SetCookie(page ActionPage, cookie agouti.Cookie) {
	checkFailure(page.SetCookie(cookie))
}

// DeleteCookie is comparable to Expect(page.DeleteCookie("cookie-name")).To(Succeed())
func DeleteCookie(page ActionPage, name string) {
	checkFailure(page.DeleteCookie(name))
}

// ClearCookies is comparable to Expect(page.ClearCookies()).To(Succeed())
func ClearCookies(page ActionPage) {
	checkFailure(page.ClearCookies())
}

// Size is comparable to Expect(page.Size(windowWidth, windowHeight)).To(Succeed())
func Size(page ActionPage, width, height int) {
	checkFailure(page.Size(width, height))
}

// Screenshot is comparable to Expect(page.Screenshot("screenshot-file.png")).To(Succeed())
func Screenshot(page ActionPage, filename string) {
	checkFailure(page.Screenshot(filename))
}

// RunScript is comparable to Expect(page.RunScript(script, args, &result)).To(Succeed())
func RunScript(page ActionPage, body string, arguments map[string]interface{}, result interface{}) {
	checkFailure(page.RunScript(body, arguments, result))
}

// EnterPopupText is comparable to Expect(page.EnterPopupText("some text")).To(Succeed())
func EnterPopupText(page ActionPage, text string) {
	checkFailure(page.EnterPopupText(text))
}

// ConfirmPopup is comparable to Expect(page.ConfirmPopup()).To(Succeed())
func ConfirmPopup(page ActionPage) {
	checkFailure(page.ConfirmPopup())
}

// CancelPopup is comparable to Expect(page.CancelPopup()).To(Succeed())
func CancelPopup(page ActionPage) {
	checkFailure(page.CancelPopup())
}

// Forward is comparable to Expect(page.Forward()).To(Succeed())
func Forward(page ActionPage) {
	checkFailure(page.Forward())
}

// Back is comparable to Expect(page.Back()).To(Succeed())
func Back(page ActionPage) {
	checkFailure(page.Back())
}

// Refresh is comparable to Expect(page.Refresh()).To(Succeed())
func Refresh(page ActionPage) {
	checkFailure(page.Refresh())
}

// SwitchToParentFrame is comparable to Expect(page.SwitchToParentFrame()).To(Succeed())
func SwitchToParentFrame(page ActionPage) {
	checkFailure(page.SwitchToParentFrame())
}

// SwitchToRootFrame is comparable to Expect(page.SwitchToRootFrame()).To(Succeed())
func SwitchToRootFrame(page ActionPage) {
	checkFailure(page.SwitchToRootFrame())
}

// SwitchToWindow is comparable to Expect(page.SwitchToWindow("window name")).To(Succeed())
func SwitchToWindow(page ActionPage, name string) {
	checkFailure(page.SwitchToWindow(name))
}

// NextWindow is comparable to Expect(page.NextWindow()).To(Succeed())
func NextWindow(page ActionPage) {
	checkFailure(page.NextWindow())
}

// CloseWindow is comparable to Expect(page.CloseWindow()).To(Succeed())
func CloseWindow(page ActionPage) {
	checkFailure(page.CloseWindow())
}
