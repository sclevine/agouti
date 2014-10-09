package page

import (
	"fmt"
	"github.com/sclevine/agouti/page/internal/webdriver"
)

type Page struct {
	Driver driver
}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]webdriver.Element, error)
	GetWindow() (webdriver.Window, error)
	SetCookie(cookie *webdriver.Cookie) error
	GetURL() (string, error)
}

func (p *Page) Navigate(url string) error {
	if err := p.Driver.Navigate(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (p *Page) SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error {
	cookie := webdriver.Cookie{name, value, path, domain, secure, httpOnly, expiry}
	if err := p.Driver.SetCookie(&cookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
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

func (p *Page) Size(height, width int) error {
	window, err := p.Driver.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to retrieve window: %s", err)
	}

	if err := window.SetSize(640, 480); err != nil {
		return fmt.Errorf("failed to set window size: %s", err)
	}

	return nil
}

func (p *Page) Find(selector string) Selection {
	return &selection{p.Driver, []string{selector}}
}

func (p *Page) Selector() string {
	return p.body().Selector()
}

func (p *Page) Click() error {
	return p.body().Click()
}

func (p *Page) Text() (string, error) {
	return p.body().Text()
}

func (p *Page) Attribute(attribute string) (string, error) {
	return p.body().Attribute(attribute)
}

func (p *Page) body() *selection {
	return &selection{p.Driver, []string{"body"}}
}
