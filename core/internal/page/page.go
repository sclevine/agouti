package page

import (
	"fmt"
	"github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Driver driver
}

type driver interface {
	DeleteSession() error
	GetWindow() (types.Window, error)
	GetScreenshot() ([]byte, error)
	SetCookie(cookie *types.Cookie) error
	DeleteCookie(name string) error
	DeleteCookies() error
	GetURL() (string, error)
	SetURL(url string) error
	GetTitle() (string, error)
	GetSource() (string, error)
	GetElements(selector types.Selector) ([]types.Element, error)
	DoubleClick() error
	MoveTo(element types.Element, point types.Point) error
	Execute(body string, arguments []interface{}, result interface{}) error
	Forward() error
	Back() error
	Refresh() error
}

func (p *Page) Destroy() error {
	if err := p.Driver.DeleteSession(); err != nil {
		return fmt.Errorf("failed to destroy session: %s", err)
	}
	return nil
}

func (p *Page) Navigate(url string) error {
	if err := p.Driver.SetURL(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (p *Page) SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error {
	cookie := types.Cookie{Name: name, Value: value, Path: path, Domain: domain, Secure: secure, HTTPOnly: httpOnly, Expiry: expiry}
	if err := p.Driver.SetCookie(&cookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
	}
	return nil
}

func (p *Page) DeleteCookie(name string) error {
	if err := p.Driver.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %s", name, err)
	}
	return nil
}

func (p *Page) ClearCookies() error {
	if err := p.Driver.DeleteCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %s", err)
	}
	return nil
}

func (p *Page) URL() (string, error) {
	url, err := p.Driver.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %s", err)
	}
	return url, nil
}

func (p *Page) Size(width, height int) error {
	window, err := p.Driver.GetWindow()
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

	screenshot, err := p.Driver.GetScreenshot()
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
	title, err := p.Driver.GetTitle()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page title: %s", err)
	}
	return title, nil
}

func (p *Page) HTML() (string, error) {
	html, err := p.Driver.GetSource()
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

	if err := p.Driver.Execute(cleanBody, values, result); err != nil {
		return fmt.Errorf("failed to run script: %s", err)
	}

	return nil
}

func (p *Page) Forward() error {
	if err := p.Driver.Forward(); err != nil {
		return fmt.Errorf("failed to navigate forward in history: %s", err)
	}
	return nil
}

func (p *Page) Back() error {
	if err := p.Driver.Back(); err != nil {
		return fmt.Errorf("failed to navigate backwards in history: %s", err)
	}
	return nil
}

func (p *Page) Refresh() error {
	if err := p.Driver.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh page: %s", err)
	}
	return nil
}

func (p *Page) Find(selector string) types.Selection {
	selection := &selection.Selection{Driver: p.Driver}
	return selection.Find(selector)
}

func (p *Page) FindXPath(selector string) types.Selection {
	selection := &selection.Selection{Driver: p.Driver}
	return selection.FindXPath(selector)
}

func (p *Page) FindLink(text string) types.Selection {
	selection := &selection.Selection{Driver: p.Driver}
	return selection.FindLink(text)
}

func (p *Page) FindByLabel(text string) types.Selection {
	selection := &selection.Selection{Driver: p.Driver}
	return selection.FindByLabel(text)
}
