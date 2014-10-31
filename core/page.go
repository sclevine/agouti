package core

import (
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
)

// A Page represents an open browser session. Pages may be created using the
// WebDriver#Page() method or by calling the Connect or SauceLabs functions.
type Page interface {
	// Selections are created using the Selectable page methods (ex. Find()).
	Selectable

	// Destroy closes the session and any open browsers processes.
	Destroy() error

	// Navigate navigates to the provided URL.
	Navigate(url string) error

	// SetCookie sets a cookie on the page.
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error

	// DeleteCookie deletes a cookie on the page by name.
	DeleteCookie(name string) error

	// ClearCookies deletes all cookies on the page.
	ClearCookies() error

	// URL returns the current page URL.
	URL() (string, error)

	// Size sets the current page size in pixels.
	Size(width, height int) error

	// Screenshot takes a screenshot and saves it to the provided filename.
	Screenshot(filename string) error

	// Title returns the page title.
	Title() (string, error)

	// HTML returns the current contents of the DOM for the entire page.
	HTML() (string, error)

	// RunScript runs the javascript provided in the body. Any keys present in
	// the arguments map will be available as variables in the body.
	// Values provided in arguments are converted into javascript objects.
	// If the body returns a value, it will be unmarshalled into the result argument.
	// Simple example:
	//    var number int
	//    page.RunScript("return test;", map[string]interface{}{"test": 100}, &number)
	//    fmt.Println(number)
	// -> 100
	RunScript(body string, arguments map[string]interface{}, result interface{}) error

	// PopupText returns the current alert, confirm, or prompt popup text.
	PopupText() (string, error)

	// EnterPopupText enters text into an open prompt popup.
	EnterPopupText(text string) error

	// ConfirmPopup confirms an alert, confirm, or prompt popup.
	ConfirmPopup() error

	// CancelPopup cancels an alert, confirm, or prompt popup.
	CancelPopup() error

	// Forward navigates forward in history.
	Forward() error

	// Back navigates backwards in history.
	Back() error

	// Refresh refreshes the page.
	Refresh() error
}

func newPage(client types.Client) Page {
	return struct {
		*page.Page
		*baseSelection
	}{
		&page.Page{Client: client},
		&baseSelection{&selection.Selection{Client: client}},
	}
}
