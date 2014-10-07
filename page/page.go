package page

import (
	"github.com/sclevine/agouti/webdriver"
)

type Page interface {
	Navigate(url string) Page
	SetCookie(cookie webdriver.Cookie) Page
	URL() string
	Selection
}

type page struct {
	driver driver
	failer failer
}

type failer interface {
	Fail(message string)
	Skip()
	UnSkip()
	Async()
	Sync()
}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]webdriver.Element, error)
	SetCookie(cookie *webdriver.Cookie) error
	GetURL() (string, error)
}

type callable interface {
	Call(selection Selection)
}

func NewPage(driver driver, failer failer) Page {
	return &page{driver, failer}
}

func (p *page) Navigate(url string) Page {
	p.failer.Skip()
	if err := p.driver.Navigate(url); err != nil {
		p.failer.Fail("Failed to navigate: "+err.Error())
	}

	return p
}

func (p *page) SetCookie(cookie webdriver.Cookie) Page {
	p.failer.Skip()
	if err := p.driver.SetCookie(&cookie); err != nil {
		p.failer.Fail("Failed to set cookie: "+err.Error())
	}

	return p
}

func (p *page) URL() string {
	p.failer.Skip()
	url, err := p.driver.GetURL()
	if err != nil {
		p.failer.Fail("Failed to retrieve URL: "+err.Error())
	}
	return url
}

func (p *page) Should() FinalSelection {
	return p.body()
}

func (p *page) ShouldNot() FinalSelection {
	body := p.body()
	body.invert = true
	return body
}

func (p *page) ShouldEventually() FinalSelection {
	return p.body().ShouldEventually()
}

func (p *page) Within(selector string, bodies ...callable) Selection {
	firstSelection := &selection{[]string{selector}, p, false}
	for _, body := range bodies {
		body.Call(firstSelection)
	}
	return firstSelection
}

func (p *page) Selector() string {
	return p.body().Selector()
}

func (p *page) Click() {
	p.failer.Skip()
	p.body().Click()
}

func (p *page) body() *selection {
	return &selection{[]string{"body"}, p, false}
}
