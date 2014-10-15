package page

import (
	"fmt"
	"github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/webdriver/types"
	"os"
	"path/filepath"
	"strings"
)

type Page interface {
	Navigate(url string) error
	SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error
	DeleteCookie(name string) error
	ClearCookies() error
	URL() (string, error)
	Size(width, height int) error
	Screenshot(filename string) error
	Title() (string, error)
	RunScript(body string, arguments map[string]interface{}, result interface{}) error
	Find(selector string) selection.Selection
	FindXPath(selector string) selection.Selection
}

type page struct {
	driver driver
}

type driver interface {
	GetWindow() (types.Window, error)
	GetScreenshot() ([]byte, error)
	SetCookie(cookie *types.Cookie) error
	DeleteCookie(name string) error
	DeleteCookies() error
	GetURL() (string, error)
	SetURL(url string) error
	GetTitle() (string, error)
	GetElements(selector types.Selector) ([]types.Element, error)
	DoubleClick() error
	MoveTo(element types.Element, point types.Point) error
	Execute(body string, arguments []interface{}, result interface{}) error
}

func New(driver driver) Page {
	return &page{driver}
}

func (p *page) Navigate(url string) error {
	if err := p.driver.SetURL(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (p *page) SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error {
	cookie := types.Cookie{name, value, path, domain, secure, httpOnly, expiry}
	if err := p.driver.SetCookie(&cookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
	}
	return nil
}

func (p *page) DeleteCookie(name string) error {
	if err := p.driver.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %s", name, err)
	}
	return nil
}

func (p *page) ClearCookies() error {
	if err := p.driver.DeleteCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %s", err)
	}
	return nil
}

func (p *page) URL() (string, error) {
	url, err := p.driver.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %s", err)
	}
	return url, nil
}

func (p *page) Size(width, height int) error {
	window, err := p.driver.GetWindow()
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

	screenshot, err := p.driver.GetScreenshot()
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
	title, err := p.driver.GetTitle()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page title: %s", err)
	}
	return title, nil
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

	if err := p.driver.Execute(cleanBody, values, result); err != nil {
		return fmt.Errorf("failed to run script: %s", err)
	}

	return nil
}

func (p *page) Find(selector string) selection.Selection {
	return selection.New(p.driver, types.Selector{"css selector", selector})
}

func (p *page) FindXPath(selector string) selection.Selection {
	return selection.New(p.driver, types.Selector{"xpath", selector})
}
