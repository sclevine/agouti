package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sclevine/agouti/core/internal/types"
)

// A Page represents a browser session. Pages may be created using the
// WebDriver#Page() method or by calling the Connect or SauceLabs functions.
type Page interface {
	// Selections are initiated using the Selectable page methods
	Selectable

	// Destroy closes the session and associated open browser instances
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

	// Size sets the current page size.
	Size(width, height int) error

	// Screenshot takes a screenshot and saves it to the provided filename.
	Screenshot(filename string) error

	// Title returns the page title.
	Title() (string, error)

	// HTML returns the current contents of the DOM for the entire page.
	HTML() (string, error)

	// RunScript runs the javascript provided in the body argument. Any keys present
	// in the arguments map will be available as variables in the body argument.
	// Arguments values are converted into javascript objects.
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

	// Refresh refreshes the page
	Refresh() error
}

type page struct {
	*baseSelection
}

func (p *page) Destroy() error {
	if err := p.Client.DeleteSession(); err != nil {
		return fmt.Errorf("failed to destroy session: %s", err)
	}
	return nil
}

func (p *page) Navigate(url string) error {
	if err := p.Client.SetURL(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (p *page) SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error {
	cookie := types.Cookie{Name: name, Value: value, Path: path, Domain: domain, Secure: secure, HTTPOnly: httpOnly, Expiry: expiry}
	if err := p.Client.SetCookie(&cookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
	}
	return nil
}

func (p *page) DeleteCookie(name string) error {
	if err := p.Client.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %s", name, err)
	}
	return nil
}

func (p *page) ClearCookies() error {
	if err := p.Client.DeleteCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %s", err)
	}
	return nil
}

func (p *page) URL() (string, error) {
	url, err := p.Client.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %s", err)
	}
	return url, nil
}

func (p *page) Size(width, height int) error {
	window, err := p.Client.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to retrieve window: %s", err)
	}

	if err := window.SetSize(width, height); err != nil {
		return fmt.Errorf("failed to set window size: %s", err)
	}

	return nil
}

func (p *page) Screenshot(filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0750); err != nil {
		return fmt.Errorf("failed to create directory for screenshot: %s", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file for screenshot: %s", err)
	}
	defer file.Close()

	screenshot, err := p.Client.GetScreenshot()
	if err != nil {
		os.Remove(filename)
		return fmt.Errorf("failed to retrieve screenshot: %s", err)
	}

	if _, err := file.Write(screenshot); err != nil {
		return fmt.Errorf("failed to write file for screenshot: %s", err)
	}

	return nil
}

func (p *page) Title() (string, error) {
	title, err := p.Client.GetTitle()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page title: %s", err)
	}
	return title, nil
}

func (p *page) HTML() (string, error) {
	html, err := p.Client.GetSource()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page HTML: %s", err)
	}
	return html, nil
}

func (p *page) RunScript(body string, arguments map[string]interface{}, result interface{}) error {
	var (
		keys   []string
		values []interface{}
	)

	for key, value := range arguments {
		keys = append(keys, key)
		values = append(values, value)
	}

	argumentList := strings.Join(keys, ", ")
	cleanBody := fmt.Sprintf("return (function(%s) { %s; }).apply(this, arguments);", argumentList, body)

	if err := p.Client.Execute(cleanBody, values, result); err != nil {
		return fmt.Errorf("failed to run script: %s", err)
	}

	return nil
}

func (p *page) PopupText() (string, error) {
	text, err := p.Client.GetAlertText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve popup text: %s", err)
	}
	return text, nil
}

func (p *page) EnterPopupText(text string) error {
	if err := p.Client.SetAlertText(text); err != nil {
		return fmt.Errorf("failed to enter popup text: %s", err)
	}
	return nil
}

func (p *page) ConfirmPopup() error {
	if err := p.Client.SetAlertText("\u000d"); err != nil {
		return fmt.Errorf("failed to confirm popup: %s", err)
	}
	return nil
}

func (p *page) CancelPopup() error {
	if err := p.Client.SetAlertText("\u001b"); err != nil {
		return fmt.Errorf("failed to cancel popup: %s", err)
	}
	return nil
}

func (p *page) Forward() error {
	if err := p.Client.Forward(); err != nil {
		return fmt.Errorf("failed to navigate forward in history: %s", err)
	}
	return nil
}

func (p *page) Back() error {
	if err := p.Client.Back(); err != nil {
		return fmt.Errorf("failed to navigate backwards in history: %s", err)
	}
	return nil
}

func (p *page) Refresh() error {
	if err := p.Client.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh page: %s", err)
	}
	return nil
}
