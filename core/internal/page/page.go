package page

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
)

type Page struct {
	Session interface {
		Delete() error
		GetWindow() (*api.Window, error)
		GetWindows() ([]*api.Window, error)
		SetWindow(window *api.Window) error
		SetWindowByName(name string) error
		DeleteWindow() error
		GetScreenshot() ([]byte, error)
		SetCookie(cookie map[string]interface{}) error
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
		AcceptAlert() error
		DismissAlert() error
		NewLogs(logType string) ([]api.Log, error)
		GetLogTypes() ([]string, error)
	}
	logs map[string][]Log
}

type Log struct {
	Message  string
	Location string
	Level    string
	Time     time.Time
}

func (p *Page) Destroy() error {
	if err := p.Session.Delete(); err != nil {
		return fmt.Errorf("failed to destroy session: %s", err)
	}
	return nil
}

func (p *Page) Navigate(url string) error {
	if err := p.Session.SetURL(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (p *Page) SetCookie(cookie map[string]interface{}) error {
	if err := p.Session.SetCookie(cookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
	}
	return nil
}

func (p *Page) DeleteCookie(name string) error {
	if err := p.Session.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %s", name, err)
	}
	return nil
}

func (p *Page) ClearCookies() error {
	if err := p.Session.DeleteCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %s", err)
	}
	return nil
}

func (p *Page) URL() (string, error) {
	url, err := p.Session.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %s", err)
	}
	return url, nil
}

func (p *Page) Size(width, height int) error {
	window, err := p.Session.GetWindow()
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

	screenshot, err := p.Session.GetScreenshot()
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
	title, err := p.Session.GetTitle()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page title: %s", err)
	}
	return title, nil
}

func (p *Page) HTML() (string, error) {
	html, err := p.Session.GetSource()
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

	if err := p.Session.Execute(cleanBody, values, result); err != nil {
		return fmt.Errorf("failed to run script: %s", err)
	}

	return nil
}

func (p *Page) PopupText() (string, error) {
	text, err := p.Session.GetAlertText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve popup text: %s", err)
	}
	return text, nil
}

func (p *Page) EnterPopupText(text string) error {
	if err := p.Session.SetAlertText(text); err != nil {
		return fmt.Errorf("failed to enter popup text: %s", err)
	}
	return nil
}

func (p *Page) ConfirmPopup() error {
	if err := p.Session.AcceptAlert(); err != nil {
		return fmt.Errorf("failed to confirm popup: %s", err)
	}
	return nil
}

func (p *Page) CancelPopup() error {
	if err := p.Session.DismissAlert(); err != nil {
		return fmt.Errorf("failed to cancel popup: %s", err)
	}
	return nil
}

func (p *Page) Forward() error {
	if err := p.Session.Forward(); err != nil {
		return fmt.Errorf("failed to navigate forward in history: %s", err)
	}
	return nil
}

func (p *Page) Back() error {
	if err := p.Session.Back(); err != nil {
		return fmt.Errorf("failed to navigate backwards in history: %s", err)
	}
	return nil
}

func (p *Page) Refresh() error {
	if err := p.Session.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh page: %s", err)
	}
	return nil
}

func (p *Page) SwitchToParentFrame() error {
	if err := p.Session.FrameParent(); err != nil {
		return fmt.Errorf("failed to switch to parent frame: %s", err)
	}
	return nil
}

func (p *Page) SwitchToRootFrame() error {
	if err := p.Session.Frame(nil); err != nil {
		return fmt.Errorf("failed to switch to original page frame: %s", err)
	}
	return nil
}

func (p *Page) SwitchToWindow(name string) error {
	if err := p.Session.SetWindowByName(name); err != nil {
		return fmt.Errorf("failed to switch to named window: %s", err)
	}
	return nil
}

func (p *Page) NextWindow() error {
	windows, err := p.Session.GetWindows()
	if err != nil {
		return fmt.Errorf("failed to find available windows: %s", err)
	}

	var windowIDs []string
	for _, window := range windows {
		windowIDs = append(windowIDs, window.ID)
	}

	// order not defined according to W3 spec
	sort.Strings(windowIDs)

	activeWindow, err := p.Session.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to find active window: %s", err)
	}

	for position, windowID := range windowIDs {
		if windowID == activeWindow.ID {
			activeWindow.ID = windowIDs[(position+1)%len(windowIDs)]
			break
		}
	}

	if err := p.Session.SetWindow(activeWindow); err != nil {
		return fmt.Errorf("failed to change active window: %s", err)
	}

	return nil
}

func (p *Page) CloseWindow() error {
	if err := p.Session.DeleteWindow(); err != nil {
		return fmt.Errorf("failed to close active window: %s", err)
	}
	return nil
}

func (p *Page) WindowCount() (int, error) {
	windows, err := p.Session.GetWindows()
	if err != nil {
		return 0, fmt.Errorf("failed to find available windows: %s", err)
	}
	return len(windows), nil
}

func (p *Page) LogTypes() ([]string, error) {
	types, err := p.Session.GetLogTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve log types: %s", err)
	}
	return types, nil
}

func (p *Page) ReadLogs(logType string, all ...bool) ([]Log, error) {
	if p.logs == nil {
		p.logs = map[string][]Log{}
	}

	clientLogs, err := p.Session.NewLogs(logType)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs: %s", err)
	}

	messageMatcher := regexp.MustCompile(`^(?s:(.+))\s\(([^)]*:\w*)\)$`)

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

// Patch to allow Page to work with new matchers after core deprecation
func (p *Page) ReadAllLogs(logType string) ([]agouti.Log, error) {
	coreLogs, err := p.ReadLogs(logType, true)
	if err != nil {
		return nil, err
	}
	var agoutiLogs []agouti.Log
	for _, log := range coreLogs {
		agoutiLogs = append(agoutiLogs, agouti.Log(log))
	}
	return agoutiLogs, nil
}

func msToTime(ms int64) time.Time {
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1000000
	return time.Unix(seconds, nanoseconds)
}
