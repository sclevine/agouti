package matchers

import (
	"github.com/onsi/gomega/types"
	"github.com/sclevine/agouti/matchers/internal"
)

// HaveTitle passes when the expected title is equivalent to the
// title of the provided page.
func HaveTitle(title string) types.GomegaMatcher {
	return &internal.ValueMatcher{Method: "Title", Property: "title", Expected: title}
}

// HaveURL passes when the expected URL is equivalent to the
// current URL of the provided page.
func HaveURL(URL string) types.GomegaMatcher {
	return &internal.ValueMatcher{Method: "URL", Property: "URL", Expected: URL}
}

// HavePopupText passes when the expected text is equivalent to the
// text contents of an open alert, confirm, or prompt popup.
func HavePopupText(text string) types.GomegaMatcher {
	return &internal.ValueMatcher{Method: "PopupText", Property: "popup text", Expected: text}
}

// HaveWindowCount passes when the expected window count is equivalent
// to the number of open windows.
func HaveWindowCount(count int) types.GomegaMatcher {
	return &internal.ValueMatcher{Method: "WindowCount", Property: "window count", Expected: count}
}

// HaveLoggedError passes when the expected log message is logged as
// an error in the browser console.
func HaveLoggedError(messageOrEmpty ...string) types.GomegaMatcher {
	message := ""
	if len(messageOrEmpty) > 0 {
		message = messageOrEmpty[0]
	}
	return &internal.HaveLoggedErrorMatcher{ExpectedMessage: message}
}

// HaveLoggedInfo passes when the expected log message is logged as
// info in the browser console.
func HaveLoggedInfo(messageOrEmpty ...string) types.GomegaMatcher {
	message := ""
	if len(messageOrEmpty) > 0 {
		message = messageOrEmpty[0]
	}
	return &internal.HaveLoggedInfoMatcher{ExpectedMessage: message}
}
