package core

import (
	"time"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/selection"
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
	SetCookie(cookie agouti.Cookie) error

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

	// RunScript runs the JavaScript provided in the body. Any keys present in
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

	// SwitchToParentFrame focuses on the immediate parent frame of a frame selected
	// by Selection#Frame. After switching, all new and existing selections will refer
	// to the parent frame. All further Page methods will apply to this frame as well.
	//
	// This method is not supported by PhantomJS. Please use SwitchToRootFrame instead.
	SwitchToParentFrame() error

	// SwitchToRootFrame focuses on the original, default page frame before any calls
	// to Selection#Frame were made. After switching, all new and existing selections
	// will refer to the root frame. All further Page methods will apply to this frame
	// as well.
	SwitchToRootFrame() error

	// SwitchToWindow switches to the first available window with the provided name
	// (JavaScript `window.name` attribute).
	SwitchToWindow(name string) error

	// NextWindow switches to the next available window.
	NextWindow() error

	// CloseWindow closes the active window.
	CloseWindow() error

	// WindowCount returns the number of available windows.
	WindowCount() (int, error)

	// ReadLogs returns log messages of the provided log type. For example,
	// page.ReadLogs("browser") returns browser console logs, such as JavaScript logs
	// and errors. If the all argument is provided as true, all logs since the session
	// was created are returned. Otherwise, only logs since the last call to ReadLogs
	// are returned. Valid log types may be obtained using the LogTypes method.
	ReadLogs(logType string, all ...bool) ([]Log, error)

	// LogTypes returns all of the valid log types that may be used with a LogReader.
	LogTypes() ([]string, error)
}

// A Log represents a single log message
type Log struct {
	// Message is the text of the log message.
	Message string

	// Location is the code location of the log message, if present
	Location string

	// Level is the log level ("DEBUG", "INFO", "WARNING", or "SEVERE").
	Level string

	// Time is the time the message was logged.
	Time time.Time
}

func newPage(session *api.Session) Page {
	pageSelection := &userSelection{selection.NewSelection(session)}
	return &userPage{&page.Page{Session: session}, pageSelection}
}
