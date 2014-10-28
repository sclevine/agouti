// Agouti core is a WebDriver API for Go.
package core

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/page"
	"github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/service"
	"github.com/sclevine/agouti/core/internal/session"
	"github.com/sclevine/agouti/core/internal/webdriver"
)

// Selection instances refer to a selection of elements.
// All Selection methods are valid MultiSelection methods.
// Examples:
// Check the first checkbox in the third row of the only table:
//    selection.Find("table").All("tr").At(2).First("td input[type=checkbox]").Check()
// Check all checkboxes in the first-and-only cell of each row in the only table:
//    selection.Find("table").All("tr").Find("td").All("input[type=checkbox]").Check()
type Selection interface {
	// The Find<X>(), First<X>(), and All<X>() methods apply their selectors to
	// each element in the selection they are called on. If the selection they are
	// called on refers to multiple elements, the resulting selection will refer
	// to at least that many elements.

	// Therefore, for each element in the current selection:

	// Find finds exactly one element by CSS selector.
	Find(selector string) *selection.Selection

	// FindByXPath finds exactly one element by XPath selector.
	FindByXPath(selector string) *selection.Selection

	// FindByLink finds exactly one element by anchor link text.
	FindByLink(text string) *selection.Selection

	// FindByLabel finds exactly one element by associated label text.
	FindByLabel(text string) *selection.Selection

	// First finds the first element by CSS selector.
	First(selector string) *selection.Selection

	// FirstByXPath finds the first element by XPath selector.
	FirstByXPath(selector string) *selection.Selection

	// FirstByLink finds the first element by anchor link text.
	FirstByLink(text string) *selection.Selection

	// FirstByLabel finds the first element by associated label text.
	FirstByLabel(text string) *selection.Selection

	// All finds zero or more elements by CSS selector.
	All(selector string) *selection.MultiSelection

	// AllByXPath finds zero or more elements by XPath selector.
	AllByXPath(selector string) *selection.MultiSelection

	// AllByLink finds zero or more elements by anchor link text.
	AllByLink(text string) *selection.MultiSelection

	// AllByLabel finds zero or more elements by associated label text.
	AllByLabel(text string) *selection.MultiSelection

	// String returns a string representation of the selection, ex.
	//    CSS: .some-class | XPath: //table [3] | Link "click me" [single]
	String() string

	// Count returns the number of elements the selection refers to.
	Count() (int, error)

	// EqualsElement returns whether or not two selections of exactly
	// one element each refer to the same element.
	EqualsElement(comparable interface{}) (bool, error)

	// Click clicks on all of the elements the selection refers to.
	Click() error

	// DoubleClick double-clicks on all of the elements the selection refers to.
	DoubleClick() error

	// Fill fills all of the input fields the selection refers to.
	Fill(text string) error

	// Check checks all of the unchecked checkboxes that the selection refers to.
	Check() error

	// Uncheck unchecks all of the checked checkboxes that the selection refers to.
	Uncheck() error

	// Select, when called on any number of <select> elements, will select all
	// <options> under those elements that match the provided text.
	Select(text string) error

	// Submit submits a form. The selection may refer to the form itself, or any
	// input element contained within the form.
	Submit() error

	// Text returns text for exactly one element.
	Text() (string, error)

	// Attribute returns an attribute value for exactly one element.
	Attribute(attribute string) (string, error)

	// CSS returns a CSS style property value for exactly one element.
	CSS(property string) (string, error)

	// Selected returns true if all of the elements that the selection
	// refers to are selected.
	Selected() (bool, error)

	// Visible returns true if all of the elements that the selection
	// refers to are visible.
	Visible() (bool, error)

	// Enabled returns true if all of the elements that the selection
	// refers to are enabled.
	Enabled() (bool, error)
}

// MultiSelection instances are Selection instances that can be indexed by element.
// A Selection returned by At(int) or Single() may still refer to multiple
// elements if any parent of the MultiSelection refers to multiple elements.
// Example:
// Submits the second form in each section:
//    selection.All("section").All("form").At(1).Submit()
// Clicks the only h1 in each div, failing if any div does not contain exactly one h1:
//    selection.All("div").Find("h1").Click()
type MultiSelection interface {
	// At specifies a single element to select using the provided index.
	At(index int) *selection.Selection

	// Single specifies that the selection must refer to exactly one element.
	//    selection.Find("#selector")
	// is equivalent to
	//    selection.All("#selector").Single()
	Single() *selection.Selection

	// All Selection methods are valid MultiSelection methods.
	Selection
}

// Page instances represent browser sessions. They can be created using the
// WebDriver#Page methods or by calling the Connect or SauceLabs functions.
type Page interface {
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

	// These Find<type>() and All<type>() selectors are equivalent to their
	// Selection counterparts but without any scoping.
	Find(selector string) *selection.Selection
	FindByXPath(selector string) *selection.Selection
	FindByLink(text string) *selection.Selection
	FindByLabel(text string) *selection.Selection
	All(selector string) *selection.MultiSelection
	AllByXPath(selector string) *selection.MultiSelection
	AllByLink(text string) *selection.MultiSelection
	AllByLabel(text string) *selection.MultiSelection
}

// SauceLabs opens a Sauce Labs session and returns a Page. Does not support Sauce Connect.
func SauceLabs(name, platform, browser, version, username, key string) (Page, error) {
	url := "http://ondemand.saucelabs.com/wd/hub"
	capabilities := session.Capabilities{
		"name":        name,
		"platform":    platform,
		"browserName": browser,
		"version":     version,
		"username":    username,
		"accessKey":   key,
	}

	pageSession, err := session.Open(url, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to Sauce Labs: %s", err)
	}

	client := &api.Client{Session: pageSession}
	return &page.Page{Client: client}, nil
}

// Connect opens a session using the provided WebDriver URL and returns a Page.
func Connect(capabilities Capabilities, url string) (Page, error) {
	pageSession, err := session.Open(url, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to Sauce Labs: %s", err)
	}

	client := &api.Client{Session: pageSession}
	return &page.Page{Client: client}, nil
}

// Capabilities is generated via the Use() function and defines a Page's desired
// capabilities. It conforms to the session.JSONable interface and is intended
// to be used with Webdriver#Page()
type Capabilities interface {
	// Browser sets the desired browser name - {chrome|firefox|safari|iphone|...}.
	Browser(browser string) session.Capabilities

	// Version sets the desired browser version (ex. "3.6").
	Version(version string) session.Capabilities

	// Platform sets the desired browser platform - {WINDOWS|XP|VISTA|MAC|LINUX|UNIX}.
	Platform(platform string) session.Capabilities

	// With enables the provided feature (ex. "handlesAlerts").
	With(feature string) session.Capabilities

	// Without disables the provided feature (ex. "javascriptEnabled").
	Without(feature string) session.Capabilities

	// JSON returns a JSON string representing the desired capabilities.
	JSON() string
}

// Use returns a Capabilities instance.
func Use() Capabilities {
	return session.Capabilities{}
}

// WebDriver represents a Selenium, PhantomJS, or ChromeDriver process.
type WebDriver interface {
	// Start launches the WebDriver process.
	Start() error

	// Stop ends all remaining sessions and stops the WebDriver process.
	Stop()

	// Page returns a new WebDriver session. The optional capabilities
	// argument allows for specification of the desired browser capabilities.
	// For Selenium, the capabilities argument should specify a browser.
	// To generate a JSONable capabilities instance, see the Use() function.
	Page(capabilities ...session.JSONable) (*page.Page, error)
}

// Chrome returns an instance of a ChromeDriver WebDriver.
func Chrome() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s", address)
	command := []string{"chromedriver", "--silent", "--port=" + port}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &webdriver.Driver{Service: service}, nil
}

// PhantomJS returns an instance of a PhantomJS WebDriver.
func PhantomJS() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	url := fmt.Sprintf("http://%s", address)
	command := []string{"phantomjs", fmt.Sprintf("--webdriver=%s", address)}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &webdriver.Driver{Service: service}, nil
}

// Selenium returns an instance of a Selenium WebDriver.
func Selenium() (WebDriver, error) {
	address, err := freeAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to locate a free port: %s", err)
	}

	port := strings.SplitN(address, ":", 2)[1]
	url := fmt.Sprintf("http://%s/wd/hub", address)
	command := []string{"selenium-server", "-port", port}
	service := &service.Service{URL: url, Timeout: 5 * time.Second, Command: command}

	return &webdriver.Driver{Service: service}, nil
}

func freeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}
