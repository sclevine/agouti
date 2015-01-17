package page

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sclevine/agouti/core/internal/api"
)

type Page struct {
	Client apiClient
	logs   map[string][]Log
}

type Log struct {
	Message  string
	Location string
	Level    string
	Time     time.Time
}

type apiClient interface {
	DeleteSession() error
	GetWindow() (*api.Window, error)
	GetWindows() ([]*api.Window, error)
	SetWindow(window *api.Window) error
	SetWindowByName(name string) error
	DeleteWindow() error
	GetScreenshot() ([]byte, error)
	SetCookie(cookie interface{}) error
	DeleteCookie(name string) error
	DeleteCookies() error
	GetURL() (string, error)
	SetURL(url string) error
	GetTitle() (string, error)
	GetSource() (string, error)
	Frame(frame *api.Element) error
	FrameParent() error
	Execute(body string, arguments []interface{}, result interface{}) error
	Forward() error
	Back() error
	Refresh() error
	GetAlertText() (string, error)
	SetAlertText(text string) error
	NewLogs(logType string) ([]api.Log, error)
	GetLogTypes() ([]string, error)
}

func (p *Page) Destroy() error {
	if err := p.Client.DeleteSession(); err != nil {
		return fmt.Errorf("failed to destroy session: %s", err)
	}
	return nil
}

func (p *Page) Navigate(url string) error {
	if err := p.Client.SetURL(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (p *Page) SetCookie(cookie interface{}) error {
	if err := p.Client.SetCookie(cookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
	}
	return nil
}

func (p *Page) DeleteCookie(name string) error {
	if err := p.Client.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %s", name, err)
	}
	return nil
}

func (p *Page) ClearCookies() error {
	if err := p.Client.DeleteCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %s", err)
	}
	return nil
}

func (p *Page) URL() (string, error) {
	url, err := p.Client.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %s", err)
	}
	return url, nil
}

func (p *Page) Size(width, height int) error {
	window, err := p.Client.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to retrieve window: %s", err)
	}

	if err := window.SetSize(width, height); err != nil {
		return fmt.Errorf("failed to set window size: %s", err)
	}

	return nil
}

func (p *Page) Screenshot(filename string) error {
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

func (p *Page) Title() (string, error) {
	title, err := p.Client.GetTitle()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page title: %s", err)
	}
	return title, nil
}

func (p *Page) HTML() (string, error) {
	html, err := p.Client.GetSource()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page HTML: %s", err)
	}
	return html, nil
}

func (p *Page) RunScript(body string, arguments map[string]interface{}, result interface{}) error {
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

func (p *Page) PopupText() (string, error) {
	text, err := p.Client.GetAlertText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve popup text: %s", err)
	}
	return text, nil
}

func (p *Page) EnterPopupText(text string) error {
	if err := p.Client.SetAlertText(text); err != nil {
		return fmt.Errorf("failed to enter popup text: %s", err)
	}
	return nil
}

func (p *Page) ConfirmPopup() error {
	if err := p.Client.SetAlertText("\u000d"); err != nil {
		return fmt.Errorf("failed to confirm popup: %s", err)
	}
	return nil
}

func (p *Page) CancelPopup() error {
	if err := p.Client.SetAlertText("\u001b"); err != nil {
		return fmt.Errorf("failed to cancel popup: %s", err)
	}
	return nil
}

func (p *Page) Forward() error {
	if err := p.Client.Forward(); err != nil {
		return fmt.Errorf("failed to navigate forward in history: %s", err)
	}
	return nil
}

func (p *Page) Back() error {
	if err := p.Client.Back(); err != nil {
		return fmt.Errorf("failed to navigate backwards in history: %s", err)
	}
	return nil
}

func (p *Page) Refresh() error {
	if err := p.Client.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh page: %s", err)
	}
	return nil
}

func (p *Page) SwitchToParentFrame() error {
	if err := p.Client.FrameParent(); err != nil {
		return fmt.Errorf("failed to switch to parent frame: %s", err)
	}
	return nil
}

func (p *Page) SwitchToRootFrame() error {
	if err := p.Client.Frame(nil); err != nil {
		return fmt.Errorf("failed to switch to original page frame: %s", err)
	}
	return nil
}

func (p *Page) SwitchToWindow(name string) error {
	if err := p.Client.SetWindowByName(name); err != nil {
		return fmt.Errorf("failed to switch to named window: %s", err)
	}
	return nil
}

func (p *Page) NextWindow() error {
	windows, err := p.Client.GetWindows()
	if err != nil {
		return fmt.Errorf("failed to find available windows: %s", err)
	}

	var windowIDs []string
	for _, window := range windows {
		windowIDs = append(windowIDs, window.ID)
	}

	// order not defined according to W3 spec
	sort.Strings(windowIDs)

	activeWindow, err := p.Client.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to find active window: %s", err)
	}

	for position, windowID := range windowIDs {
		if windowID == activeWindow.ID {
			activeWindow.ID = windowIDs[(position+1)%len(windowIDs)]
			break
		}
	}

	if err := p.Client.SetWindow(activeWindow); err != nil {
		return fmt.Errorf("failed to change active window: %s", err)
	}

	return nil
}

func (p *Page) CloseWindow() error {
	if err := p.Client.DeleteWindow(); err != nil {
		return fmt.Errorf("failed to close active window: %s", err)
	}
	return nil
}

func (p *Page) WindowCount() (int, error) {
	windows, err := p.Client.GetWindows()
	if err != nil {
		return 0, fmt.Errorf("failed to find available windows: %s", err)
	}
	return len(windows), nil
}

func (p *Page) LogTypes() ([]string, error) {
	types, err := p.Client.GetLogTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve log types: %s", err)
	}
	return types, nil
}

func (p *Page) ReadLogs(logType string, all ...bool) ([]Log, error) {
	if p.logs == nil {
		p.logs = map[string][]Log{}
	}

	clientLogs, err := p.Client.NewLogs(logType)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs: %s", err)
	}

	messageMatcher := regexp.MustCompile(`^(?s:(.+))\s\(([^)]+:\w+)\)$`)

	var logs []Log
	for _, clientLog := range clientLogs {
		matches := messageMatcher.FindStringSubmatch(clientLog.Message)
		message, location := clientLog.Message, ""
		if len(matches) > 2 {
			message, location = matches[1], matches[2]
		}

		log := Log{message, location, clientLog.Level, msToTime(clientLog.Timestamp)}
		logs = append(logs, log)
	}
	p.logs[logType] = append(p.logs[logType], logs...)

	if len(all) > 0 && all[0] {
		return p.logs[logType], nil
	}

	return logs, nil
}

func msToTime(ms int64) time.Time {
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1000000
	return time.Unix(seconds, nanoseconds)
}
