package mocks

import "github.com/sclevine/agouti/core/internal/selection"

type Page struct {
	TitleCall struct {
		ReturnTitle string
		Err         error
	}
}

func (p *Page) Navigate(url string) error {
	return nil
}

func (p *Page) SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error {
	return nil
}

func (p *Page) DeleteCookie(name string) error {
	return nil
}

func (p *Page) ClearCookies() error {
	return nil
}

func (p *Page) URL() (string, error) {
	return "", nil
}

func (p *Page) Size(width, height int) error {
	return nil
}

func (p *Page) Screenshot(filename string) error {
	return nil
}

func (p *Page) Title() (string, error) {
	return p.TitleCall.ReturnTitle, p.TitleCall.Err
}

func (p *Page) RunScript(body string, arguments map[string]interface{}, result interface{}) error {
	return nil
}

func (p *Page) Find(selector string) selection.Selection {
	return nil
}

func (p *Page) FindXPath(selector string) selection.Selection {
	return nil
}
